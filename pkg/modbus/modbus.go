// Copyright 2014 Quoc-Viet Nguyen. All rights reserved.
// This software may be modified and distributed under the terms
// of the BSD license. See the LICENSE file for details.

/*
Package modbus provides a client for MODBUS TCP and RTU/ASCII.
*/
package modbus

import (
	"fmt"
)

const (
	// Bit access
	FuncCodeReadDiscreteInputs = 2
	FuncCodeReadCoils          = 1
	FuncCodeWriteSingleCoil    = 5
	FuncCodeWriteMultipleCoils = 15

	// 16-bit access
	FuncCodeReadInputRegisters         = 4
	FuncCodeReadHoldingRegisters       = 3
	FuncCodeWriteSingleRegister        = 6
	FuncCodeWriteMultipleRegisters     = 16
	FuncCodeReadWriteMultipleRegisters = 23
	FuncCodeMaskWriteRegister          = 22
	FuncCodeReadFIFOQueue              = 24
	FuncAnswerSuccessRegisters         = 113
)

const (
	ExceptionCodeIllegalFunction                    = 1
	ExceptionCodeIllegalDataAddress                 = 2
	ExceptionCodeIllegalDataValue                   = 3
	ExceptionCodeServerDeviceFailure                = 4
	ExceptionCodeAcknowledge                        = 5
	ExceptionCodeServerDeviceBusy                   = 6
	ExceptionCodeMemoryParityError                  = 8
	ExceptionCodeGatewayPathUnavailable             = 10
	ExceptionCodeGatewayTargetDeviceFailedToRespond = 11
)
const (
	rtuMinSize       = 4
	rtuMaxSize       = 256
	rtuExceptionSize = 5
)

var (
	Success = []byte{0x0, 0x1}
	Error   = []byte{0x0, 0x0}
)

// ModbusError implements error interface.
type ModbusError struct {
	FunctionCode  byte
	ExceptionCode byte
}

// Error converts known modbus exception code to error message.
func (e *ModbusError) Error() string {
	var name string
	switch e.ExceptionCode {
	case ExceptionCodeIllegalFunction:
		name = "illegal function"
	case ExceptionCodeIllegalDataAddress:
		name = "illegal data address"
	case ExceptionCodeIllegalDataValue:
		name = "illegal data value"
	case ExceptionCodeServerDeviceFailure:
		name = "server logic failure"
	case ExceptionCodeAcknowledge:
		name = "acknowledge"
	case ExceptionCodeServerDeviceBusy:
		name = "server logic busy"
	case ExceptionCodeMemoryParityError:
		name = "memory parity error"
	case ExceptionCodeGatewayPathUnavailable:
		name = "processor path unavailable"
	case ExceptionCodeGatewayTargetDeviceFailedToRespond:
		name = "processor target logic failed to respond"
	default:
		name = "unknown"
	}
	return fmt.Sprintf("modbus: exception '%v' (%s), function '%v'", e.ExceptionCode, name, e.FunctionCode)
}

// ProtocolDataUnit (PDU) is independent of underlying communication layers.
type ProtocolDataUnit struct {
	SlaveId      byte   //从机地址
	FunctionCode byte   //功能码
	Data         []byte //内容
}
type ResultProtocolDataStr struct {
	SlaveId      uint8  //从机地址
	FunctionCode uint8  //功能码
	Data         string //内容
	DataType     []byte //数据浮点
}

// Packager specifies the communication layer.
type Packager interface {
	Encode(pdu *ProtocolDataUnit) (adu []byte, err error)
	Decode(adu []byte) (pdu *ProtocolDataUnit, err error)
	Verify(aduRequest []byte, aduResponse []byte) (err error)
}

// Transporter specifies the transport layer.
type Transporter interface {
	Send(aduRequest []byte) (aduResponse []byte, err error)
}
