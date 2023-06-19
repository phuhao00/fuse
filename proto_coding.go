package fuse

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
)

// Unmarshal ...
func (r *Router) Unmarshal(data []byte, msg interface{}) error {
	if len(data) < 2 {
		return errors.New("protobuf data too short")
	}
	pbMsg, ok := msg.(proto.Message)
	if !ok {
		return fmt.Errorf("msg is not protobuf message")
	}
	return proto.Unmarshal(data, pbMsg)
}

// Marshal ...
func (r *Router) Marshal(msgID uint16, msg interface{}) ([]byte, error) {
	pbMsg, ok := msg.(proto.Message)
	if !ok {
		return []byte{}, fmt.Errorf("msg is not protobuf message")
	}
	data, err := proto.Marshal(pbMsg)
	if err != nil {
		return data, err
	}
	return data, err
}
