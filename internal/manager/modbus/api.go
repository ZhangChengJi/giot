package modbus

type Client interface {
	// ReadInputRegisters reads from 1 to 125 contiguous input registers in
	// a remote device and returns input registers.
	ReadInputRegisters(slaveId byte, address, quantity uint16) (results []byte, err error)
	// ReadHoldingRegisters reads the contents of a contiguous block of
	// holding registers in a remote logic and returns register value.
	ReadHoldingRegisters(slaveId byte, address, quantity uint16) (results []byte, err error)
}
