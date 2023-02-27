package handshake

import (
	"bytes"
	"crypto/rand"
	"errors"
	"io"
	"time"
)

type S0C0 byte

type S1C1 struct {
	Time    uint32
	Version [4]byte
	Random  [1528]byte
}

type S2C2 struct {
	Time   uint32
	Time2  uint32
	Random [1528]byte
}

var Version = [4]byte{0, 0, 0, 0}

var timeNow = time.Now

func HandshakeWithClient(r io.Reader, w io.Writer) error {
	d := NewDecoder(r)
	e := NewEncoder(w)

	var c0 S0C0
	if err := d.DecoderS0C0(&c0); err != nil {
		return err
	}

	s0 := S0C0(0)
	if err := e.EncodeS0C0(&s0); err != nil {
		return err
	}

	s1 := S1C1{
		Time: uint32(timeNow().UnixNano() / int64(time.Millisecond)),
	}
	copy(s1.Version[:], Version[:])
	if _, err := rand.Read(s1.Random[:]); err != nil {
		return err
	}
	if err := e.EncodeS1C1(&s1); err != nil {
		return err
	}

	var c1 S1C1
	if err := d.DecodeS1C1(&c1); err != nil {
		return err
	}

	s2 := S2C2{
		Time:  c1.Time,
		Time2: uint32(timeNow().UnixNano() / int64(time.Microsecond)),
	}
	copy(s2.Random[:], c1.Random[:])
	if err := e.EncodeS2C2(&s2); err != nil {
		return err
	}

	var c2 S2C2
	if err := d.DecodeS2C2(&c2); err != nil {
		return err
	}

	if !bytes.Equal(c2.Random[:], s1.Random[:]) {
		return errors.New("Random echo is not matched")
	}

	return nil
}
