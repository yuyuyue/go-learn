package frame

import (
	"encoding/binary"
	"errors"
	"io"
)

/*
Frame定义

frameHeader + payload(packet)

frameHeader

	4 bytes: length 整型，帧总长度(含头及payload)

payload

	Packet
*/
type FramePayload []byte

type FrameStream interface {
	Encode(w io.Writer, payload FramePayload) error
	Decode(r io.Reader) (FramePayload, error)
}

var ErrShortWrite = errors.New("short write")
var ErrShortRead = errors.New("short read")

type FrameCodeStruct struct{}

func NewFrameStream() FrameStream {
	return &FrameCodeStruct{}
}

func (fs *FrameCodeStruct) Encode(w io.Writer, payload FramePayload) error {
	var f = payload
	var totalLen int32 = int32(len(payload)) + 4

	err := binary.Write(w, binary.BigEndian, &totalLen)
	if err != nil {
		return err
	}

	n, err := w.Write([]byte(f)) // write the frame payload to outbound stream
	if err != nil {
		return err
	}

	if n != len(payload) {
		return ErrShortWrite
	}
	return nil
}

func (fs *FrameCodeStruct) Decode(r io.Reader) (FramePayload, error) {
	var totalLen int32
	err := binary.Read(r, binary.BigEndian, &totalLen)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, totalLen-4)
	n, err := io.ReadFull(r, buf)
	if err != nil {
		return nil, err
	}

	if n != int(totalLen-4) {
		return nil, ErrShortRead
	}

	return FramePayload(buf), nil
}
