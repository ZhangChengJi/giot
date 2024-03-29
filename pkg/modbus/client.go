package modbus

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"giot/utils/encoding"
)

type client struct {
	rtuHandler *RtuHandler
}

// NewClient creates a new modbus client with given backend handler.
func NewClient(handler *RtuHandler) Client {
	return &client{rtuHandler: handler}
}
func (c *client) ReadHoldingRegisters(slaveId byte, address, quantity uint16) (results []byte, err error) {
	c.rtuHandler.SlaveId = slaveId
	if quantity < 1 || quantity > 125 {
		err = fmt.Errorf("modbus: quantity '%v' must be between '%v' and '%v',", quantity, 1, 125)
		return
	}
	request := ProtocolDataUnit{
		FunctionCode: FuncCodeReadHoldingRegisters,
		Data:         dataBlock(address, quantity),
	}
	return c.rtuHandler.Encode(&request)
}

// Request:
//  Function code         : 1 byte (0x04)
//  Starting address      : 2 bytes
//  Quantity of registers : 2 bytes
// Response:
//  Function code         : 1 byte (0x04)
//  Byte count            : 1 byte
//  Input registers       : N bytes
func (c *client) ReadInputRegisters(slaveId byte, address, quantity uint16) (results []byte, err error) {
	c.rtuHandler.SlaveId = slaveId
	if quantity < 1 || quantity > 125 {
		err = fmt.Errorf("modbus: quantity '%v' must be between '%v' and '%v',", quantity, 1, 125)
		return
	}
	request := ProtocolDataUnit{
		FunctionCode: FuncCodeReadInputRegisters,
		Data:         dataBlock(address, quantity),
	}
	return c.rtuHandler.Encode(&request)
}
func (c *client) WriteSingleRegister(slaveId byte, address, quantity uint16, value []byte) (results []byte, err error) {
	c.rtuHandler.SlaveId = slaveId
	request := ProtocolDataUnit{
		FunctionCode: FuncAnswerSuccessRegisters,
		Data:         dataBlock1(value, address, quantity),
	}
	return c.rtuHandler.Encode(&request)
}
func (c *client) ReadIndustryCode(data []byte) (result *ProtocolDataUnit, err error) {
	return c.rtuHandler.Decode(data)
}
func (c *client) ReadIndustryF1Code(data []byte) (result []*ProtocolDataUnit, err error) {
	return c.rtuHandler.F1Decode(data)
}

func (c *client) ReadHomeCode(data []byte) (pdu *ProtocolDataUnit, err error) {
	return c.rtuHandler.HomeDecode(data)
}

// Request:
//  Function code         : 1 byte (0x10)
//  Starting address      : 2 bytes
//  Quantity of outputs   : 2 bytes
//  Byte count            : 1 byte
//  Registers value       : N* bytes
// Response:
//  Function code         : 1 byte (0x10)
//  Starting address      : 2 bytes
//  Quantity of registers : 2 bytes
func (c *client) WriteMultipleRegisters(address, quantity uint16, value []byte) (results []byte, err error) {
	if quantity < 1 || quantity > 123 {
		err = fmt.Errorf("modbus: quantity '%v' must be between '%v' and '%v',", quantity, 1, 123)
		return
	}
	a := encoding.BytesToFloat32s(encoding.BIG_ENDIAN, encoding.HIGH_WORD_FIRST, dataBlockSuffix(value, address, quantity))
	fmt.Println(a)
	request := ProtocolDataUnit{
		FunctionCode: FuncCodeWriteMultipleRegisters,
		Data:         dataBlockSuffix(value, address, quantity),
	}
	ad, _ := c.rtuHandler.Encode(&request)
	return ad, nil
}
func (c *client) CheckCrc(data []byte) (err error) {
	length := len(data)
	var crc crc

	crc.reset().pushBytes(data[0 : length-2])
	checksum := uint16(data[length-1])<<8 | uint16(data[length-2])
	if checksum != crc.value() {
		err = fmt.Errorf("modbus: response crc '%v' does not match expected '%v'", checksum, crc.value())
		return err
	}
	return nil

}

type RtuHandler struct {
	SlaveId byte
}

// Encode encodes PDU in a RTU frame:
//  Slave Address   : 1 byte
//  Function        : 1 byte
//  Data            : 0 up to 252 bytes
//  CRC             : 2 byte
func (mb *RtuHandler) Encode(pdu *ProtocolDataUnit) (adu []byte, err error) {
	length := len(pdu.Data) + 4
	if length > rtuMaxSize {
		err = fmt.Errorf("modbus: length of data '%v' must not be bigger than '%v'", length, rtuMaxSize)
		return
	}
	adu = make([]byte, length)

	adu[0] = mb.SlaveId
	adu[1] = pdu.FunctionCode
	copy(adu[2:], pdu.Data)

	// Append crc
	var crc crc
	crc.reset().pushBytes(adu[0 : length-2])
	checksum := crc.value()

	adu[length-1] = byte(checksum >> 8)
	adu[length-2] = byte(checksum)

	return
}

// Verify verifies response length and slave id.
func (mb *RtuHandler) Verify(aduRequest []byte, aduResponse []byte) (err error) {
	length := len(aduResponse)
	// Minimum size (including address, function and CRC)
	if length < rtuMinSize {
		err = fmt.Errorf("modbus: response length '%v' does not meet minimum '%v'", length, rtuMinSize)
		return
	}
	// Slave address must match
	if aduResponse[0] != aduRequest[0] {
		err = fmt.Errorf("modbus: response slave id '%v' does not match request '%v'", aduResponse[0], aduRequest[0])
		return
	}
	return
}

// readHex decodes hexa string to byte, e.g. "8C" => 0x8C.
func readHex(data []byte) (value byte, err error) {
	var dst [1]byte
	if _, err = hex.Decode(dst[:], data[0:2]); err != nil {
		fmt.Println(err)
		return
	}
	value = dst[0]
	return
}

// Decode extracts PDU from RTU frame and verify CRC. Decode 解码
func (mb *RtuHandler) Decode(adu []byte) (result *ProtocolDataUnit, err error) {
	length := len(adu)
	// Calculate checksum
	var crc crc

	crc.reset().pushBytes(adu[0 : length-2])
	checksum := uint16(adu[length-1])<<8 | uint16(adu[length-2])
	if checksum != crc.value() {
		err = fmt.Errorf("modbus: response crc '%v' does not match expected '%v'", checksum, crc.value())
		return
	}
	// Function code & data
	result = &ProtocolDataUnit{}
	result.SlaveId = adu[0]      //从机id
	result.FunctionCode = adu[1] //功能码

	result.Data = adu[3 : len(adu)-2] //数据
	return
}

// F1Decode 主机 extracts PDU from RTU frame and verify CRC. Decode 解码
func (mb *RtuHandler) F1Decode(adu []byte) (result []*ProtocolDataUnit, err error) {
	length := len(adu)
	// Calculate checksum
	var crc crc
	crc.reset().pushBytes(adu[0 : length-2])
	checksum := uint16(adu[length-1])<<8 | uint16(adu[length-2])
	if checksum != crc.value() {
		err = fmt.Errorf("modbus: response crc '%v' does not match expected '%v'", checksum, crc.value())
		return
	}
	var s uint8 = 1
	data := adu[3 : len(adu)-2]
	var size uint8 = adu[2]
	for i := 0; i < int(size); i += 2 {
		protocolDataUnit := &ProtocolDataUnit{}
		protocolDataUnit.SlaveId = s
		protocolDataUnit.FunctionCode = adu[1]
		protocolDataUnit.Data = data[i : i+2]
		result = append(result, protocolDataUnit)
		s++
	}
	return
}

var homeCode = [...]byte{0x70, 0x65, 0x66, 0x63, 0x64}

func (mb *RtuHandler) HomeDecode(adu []byte) (pdu *ProtocolDataUnit, err error) {
	length := len(adu)
	// Calculate checksum
	var crc crc
	crc.reset().pushBytes(adu[0 : length-2])
	checksum := uint16(adu[length-1])<<8 | uint16(adu[length-2])
	if checksum != crc.value() {
		err = fmt.Errorf("modbus: response crc '%v' does not match expected '%v'", checksum, crc.value())
		return
	}
	// Function code & data
	pdu = &ProtocolDataUnit{}
	pdu.SlaveId = adu[0] //从机id
	for _, b := range homeCode {
		if b == adu[1] {
			pdu.FunctionCode = adu[1] //功能码
			break
		}
	}

	pdu.Data = adu[24:32] //数据

	return
}

// dataBlock creates a sequence of uint16 data.
func dataBlock(value ...uint16) []byte {
	data := make([]byte, 2*len(value))
	for i, v := range value {
		binary.BigEndian.PutUint16(data[i*2:], v)
	}

	return data
}
func dataBlock1(suffix []byte, value ...uint16) []byte {
	length := 2 * len(value)
	data := make([]byte, length+len(suffix))
	for i, v := range value {
		binary.BigEndian.PutUint16(data[i*2:], v)
	}
	//data[length] = uint8(len(suffix))
	copy(data[length:], suffix)
	return data
}

// dataBlockSuffix creates a sequence of uint16 data and append the suffix plus its length.
func dataBlockSuffix(suffix []byte, value ...uint16) []byte {
	length := 2 * len(value)
	data := make([]byte, length+1+len(suffix))
	for i, v := range value {
		binary.BigEndian.PutUint16(data[i*2:], v)
	}
	data[length] = uint8(len(suffix))
	copy(data[length+1:], suffix)
	return data
}

//func responseError(response *ProtocolDataUnit) error {
//	mbError := &ModbusError{FunctionCode: response.FunctionCode}
//	if response.Data != nil && len(response.Data) > 0 {
//		mbError.ExceptionCode = response.Data[0]
//	}
//	return mbError
//}
