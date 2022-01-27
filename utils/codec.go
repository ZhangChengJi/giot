package utils

import (
	"errors"
	"giot/utils/json"
	"github.com/vmihailenco/msgpack/v5"
)

var ErrUnsupportCodec = errors.New("unsupported codec")

type CodecMethodEnum int

const (
	Json CodecMethodEnum = iota
	MessagePack
)

type Codec interface {
	CodecMethod() CodecMethodEnum
	Partition() int32
}

func ToKafkaBytes(k Codec) ([]byte, error) {
	switch k.CodecMethod() {
	case Json:
		return json.Marshal(k)
	case MessagePack:
		return msgpack.Marshal(k)
	default:
		return nil, ErrUnsupportCodec
	}
}

func FromKafkaBytes(bytes []byte, record Codec) error {
	switch record.CodecMethod() {
	case Json:
		return json.Unmarshal(bytes, record)
	case MessagePack:
		return msgpack.Unmarshal(bytes, record)
	default:
		return ErrUnsupportCodec
	}
}
