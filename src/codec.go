package src

import (
	"bytes"
	"encoding/binary"
	"net"
)

type MsgIdType uint32
type Codec interface {
	Encode(MsgIdType, []byte) ([]byte, error)
	Decode(net.Conn) (MsgIdType, []byte, error)
}

type TestCodec struct{}

func (t TestCodec) Encode(msgId MsgIdType, data []byte) ([]byte, error) {

	var tarData []byte
	bf := bytes.NewBuffer(tarData)
	binary.Write(bf, binary.LittleEndian, msgId)
	binary.Write(bf, binary.LittleEndian, uint32(len(data)))
	binary.Write(bf, binary.LittleEndian, data)
	return bf.Bytes(), nil
}

func (t TestCodec) Decode(conn net.Conn) (MsgIdType, []byte, error) {
	var (
		msgId  MsgIdType
		msgLen uint32
	)

	// binary.LittleEndian.PutUint32(data, uint32(msgId))
	// msgId = MsgIdType(binary.LittleEndian.Uint32(data))

	binary.Read(conn, binary.LittleEndian, &msgId)

	binary.Read(conn, binary.LittleEndian, &msgLen)

	tarData := make([]byte, msgLen)
	binary.Read(conn, binary.LittleEndian, &tarData)
	return msgId, tarData, nil
}
