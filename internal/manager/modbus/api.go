package modbus

type Client interface {
	// ReadHoldingRegisters reads the contents of a contiguous block of
	// holding registers in a remote device and returns register value.
	ReadHoldingRegisters(address, quantity uint16) (results []byte, err error) //获取寄存器信息
}
