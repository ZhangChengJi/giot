package modbus

type Client interface {
	// ReadInputRegisters reads from 1 to 125 contiguous input registers in
	// a remote transfer and returns input registers.
	ReadInputRegisters(slaveId byte, address, quantity uint16) (results []byte, err error)
	// ReadHoldingRegisters reads the contents of a contiguous block of
	// holding registers in a remote logic and returns register value.
	ReadHoldingRegisters(slaveId byte, address, quantity uint16) (results []byte, err error)
	WriteSingleRegister(slaveId byte, address, quantity uint16, value []byte) (results []byte, err error)
	WriteMultipleRegisters(address, quantity uint16, value []byte) (results []byte, err error)

	ReadIndustryCode(data []byte) (result *ResultProtocolDataUnit16, err error)
	ReadHomeCode(data []byte) (result *ResultProtocolDataUnit16, err error)

	CheckCrc(data []byte) (err error)
}
