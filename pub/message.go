package pub

import (
	"encoding/binary"
	"encoding/json"
)

const (
	magic    = "chat"
	magicLen = len(magic)

	dataSizeLen = 2
	headerLen   = magicLen + dataSizeLen

	MessageTypeSize = 2
)

type Message interface {
	MsgID() uint16
}

// 消息序列化
func MessagePack(msg Message) (res []byte, err error) {
	bytes, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	// 创建最终序列化数组，结构为4+2+2+len(data)
	b := make([]byte, magicLen+dataSizeLen+MessageTypeSize+len(bytes))
	// [0:4] 放chat
	copy(b, []byte(magic))
	// [4:6] 放datalen 等与实际的数据长度+消息头
	dataLen := len(bytes) + MessageTypeSize
	binary.BigEndian.PutUint16(b[magicLen:], uint16(dataLen))
	// [6:8] 放消息头编号
	binary.BigEndian.PutUint16(b[headerLen:], msg.MsgID())
	// [8:] 放实际的数据
	copy(b[headerLen+MessageTypeSize:], bytes)
	return b, nil
}

// 消息反序列化
func MessageUnpack(prtToObj Message, data []byte) error {
	// 进入到这儿时，data就是实际的数据
	return json.Unmarshal(data, prtToObj)
}
