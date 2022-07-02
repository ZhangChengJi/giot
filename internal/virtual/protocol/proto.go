package protocol

import (
	"bytes"
	"encoding/binary"
	"errors"
	"giot/pkg/log"
	"giot/pkg/modbus"

	"github.com/panjf2000/gnet/v2"
)

var ErrIncompletePacket = errors.New("incomplete packet")

const (
	magicNumber     = 1314
	magicNumberSize = 2
	bodySize        = 4
)

var magicNumberBytes []byte

func init() {
	magicNumberBytes = make([]byte, magicNumberSize)
	binary.BigEndian.PutUint16(magicNumberBytes, uint16(magicNumber))
}

// SimpleCodec Protocol format:
//
// * 0           2                       6
// * +-----------+-----------------------+
// * |   magic   |       body len        |
// * +-----------+-----------+-----------+
// * |                                   |
// * +                                   +
// * |           body bytes              |
// * +                                   +
// * |            ... ...                |
// * +-----------------------------------+
type ModbusCodec struct {
	modbus.Client
}

//加
func (codec ModbusCodec) Encode(buf []byte) (*modbus.ResultProtocolDataUnit16, error) {
	data, err := codec.ReadHomeCode(buf) //解码
	if err != nil {
		log.Errorf("data Decode failed:%s", err)
		return nil, err
	}
	return data, nil
}

//解
func (codec *ModbusCodec) Decode(c gnet.Conn) ([]byte, error) {
	bodyOffset := magicNumberSize + bodySize
	buf, _ := c.Peek(bodyOffset)
	if len(buf) < bodyOffset {
		return nil, ErrIncompletePacket
	}

	if !bytes.Equal(magicNumberBytes, buf[:magicNumberSize]) {
		return nil, errors.New("invalid magic number")
	}

	bodyLen := binary.BigEndian.Uint32(buf[magicNumberSize:bodyOffset])
	msgLen := bodyOffset + int(bodyLen)
	if c.InboundBuffered() < msgLen {
		return nil, ErrIncompletePacket
	}
	buf, _ = c.Peek(msgLen)
	_, _ = c.Discard(msgLen)

	return buf[bodyOffset:msgLen], nil
}

func (codec ModbusCodec) Unpack(buf []byte) ([]byte, error) {
	bodyOffset := magicNumberSize + bodySize
	if len(buf) < bodyOffset {
		return nil, ErrIncompletePacket
	}

	if !bytes.Equal(magicNumberBytes, buf[:magicNumberSize]) {
		return nil, errors.New("invalid magic number")
	}

	bodyLen := binary.BigEndian.Uint32(buf[magicNumberSize:bodyOffset])
	msgLen := bodyOffset + int(bodyLen)
	if len(buf) < msgLen {
		return nil, ErrIncompletePacket
	}

	return buf[bodyOffset:msgLen], nil
}
