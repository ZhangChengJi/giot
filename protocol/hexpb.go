package protocol

import (
	"encoding/json"
	"errors"
	"fmt"
	"giot/internal/manager/modbus"
	"github.com/panjf2000/gnet"
)

const (
	GUIDLength = 24
	G          = iota
	D
)

var as int8 = 0

type RtuData struct {
	DataType int8
	Guid     string
	Payload  string
}

type HexLengthFieldProtocol struct {
	h modbus.RtuHandler
	RtuData
}

// Encode 编码
func (hex *HexLengthFieldProtocol) Encode(c gnet.Conn, buf []byte) ([]byte, error) {

	return buf, nil
}

// Decode 解码
func (hex *HexLengthFieldProtocol) Decode(c gnet.Conn) ([]byte, error) {
	fmt.Println("1")
	buf := c.Read()
	length := len(buf)
	if length > 0 && length > 7 {
		if length >= 24 {
			hex.DataType = G
			hex.Payload = string(buf)
			fmt.Println(hex.RtuData)
			hex.Guid = string(buf)
			hex.DataType = as + 1

			a, _ := json.Marshal(hex.RtuData)
			return a, nil
		} else {
			//正常上数据
			response, err := hex.h.Decode(buf)
			if err != nil {
				return nil, err
			}
			if response.Data == nil || len(response.Data) == 0 {
				// Empty response
				err = fmt.Errorf("modbus: response data is empty")
				return nil, err
			}
			return response.Data, nil
		}
	} else {
		return nil, errors.New("logic payload failed")
	}

}
