package code

import (
	"bytes"
	"errors"

	"github.com/vmihailenco/msgpack/v5"
)

// 使用msgpack 作为序列化方式
type Serialization interface {
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte, interface{}) error
}

var DefaultSerialization = NewSerialization()

var NewSerialization = func() Serialization {
	return &MsgpackSerialization{}
}

type MsgpackSerialization struct{}

func (c *MsgpackSerialization) Marshal(v interface{}) ([]byte, error) {
	if v == nil {
		return nil, errors.New("marshal nil interface{}")
	}

	return msgpack.Marshal(v)
}

func (c *MsgpackSerialization) Unmarshal(data []byte, v interface{}) error {
	if len(data) == 0 {
		return errors.New("unmarshal nil or empty bytes")
	}

	decoder := msgpack.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(v)
	return err
}
