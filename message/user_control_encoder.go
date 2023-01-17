package message

import (
	"encoding/binary"
	"io"
)

type UserControlEventEncoder struct {
	w io.Writer
}

func NewUserControlEventEncoder(w io.Writer) *UserControlEventEncoder {
	return &UserControlEventEncoder{
		w: w,
	}
}

func (enc *UserControlEventEncoder) Encode(msg interface{}) error {
	switch msg := msg.(type) {
	case *StreamBegin:
		return enc.encodeStreamBegin(msg)
	default:
		panic("unreachable")
	}
}

func (enc *UserControlEventEncoder) encodeStreamBegin(msg *StreamBegin) error {
	buf := make([]byte, 2+4)
	binary.BigEndian.PutUint16(buf[0:2], 0)
	binary.BigEndian.PutUint32(buf[2:], msg.StreamID)

	_, err := enc.w.Write(buf)

	return err
}
