package rtmp

import (
	"encoding/binary"
	"io"
)

type chunkBasicHeader struct {
	fmt           byte
	chunkStreamID int
}

func decodeChunkBasicHeader(r io.Reader, bh *chunkBasicHeader) error {
	buf := make([]byte, 3)
	_, err := io.ReadAtLeast(r, buf[:1], 1)
	if err != nil {
		return err
	}

	fmt := (buf[0] & 0xC0 >> 6)
	csID := int(buf[0] & 0x3f)

	switch csID {
	case 0:
		panic("not implemented")
	case 1:
		panic("not implemented")
	}

	bh.fmt = fmt
	bh.chunkStreamID = csID

	return nil
}

func encodeChunkBasicHeader(w io.Writer, mh *chunkBasicHeader) error {
	buf := make([]byte, 3)
	buf[0] = byte(mh.fmt&0x03) << 6

	switch {
	case mh.chunkStreamID >= 2 && mh.chunkStreamID <= 63:
		buf[0] |= byte(mh.chunkStreamID & 0x3f)
		_, err := w.Write(buf[:1])
		return err
	case mh.chunkStreamID >= 64 && mh.chunkStreamID <= 319:
		panic("not implemented")
	case mh.chunkStreamID >= 320 && mh.chunkStreamID <= 65599:
		panic("not implemented")
	default:
		panic("unexpected chunk stream id")
	}
}

type chunkMessageHeader struct {
	timestamp       uint32
	timestampDelta  uint32
	messageLength   uint32
	messageTypeID   byte
	messageStreamID uint32
}

func decodeChunkMessageHeader(r io.Reader, fmt byte, mh *chunkMessageHeader) error {
	cache32bits := make([]byte, 4)

	switch fmt {
	case 0:
		buf := make([]byte, 11)
		_, err := io.ReadAtLeast(r, buf, len(buf))
		if err != nil {
			return err
		}

		copy(cache32bits[1:], buf[0:3])
		mh.timestamp = binary.BigEndian.Uint32(cache32bits)
		copy(cache32bits[1:], buf[3:6])
		mh.messageLength = binary.BigEndian.Uint32(cache32bits)
		mh.messageTypeID = buf[6]
		mh.messageStreamID = binary.LittleEndian.Uint32(buf[7:11])

		if mh.timestamp == 0xfffffff {
			panic("not implemented extended timestamp")
		}
	case 2:
		buf := make([]byte, 3)
		_, err := io.ReadAtLeast(r, buf, len(buf))
		if err != nil {
			return err
		}

		copy(cache32bits[1:], buf[0:3])
		mh.timestampDelta = binary.BigEndian.Uint32(cache32bits)

		if mh.timestampDelta == 0xffffff {
			panic("not implemented extended timestamp delta")
		}

	case 3:
	default:
		panic("unexpected fmt")
	}

	return nil
}

func encodeChunkMessageHeader(w io.Writer, fmt byte, mh *chunkMessageHeader) error {
	cache32bits := make([]byte, 4)

	switch fmt {
	case 0:
		buf := make([]byte, 11)
		binary.BigEndian.PutUint32(cache32bits, mh.timestamp)
		copy(buf[0:3], cache32bits[1:])
		binary.BigEndian.PutUint32(cache32bits, mh.messageLength)
		copy(buf[3:6], cache32bits[1:])
		buf[6] = mh.messageTypeID
		binary.LittleEndian.PutUint32(buf[7:11], mh.messageStreamID)

		_, err := w.Write(buf)
		return err

	case 1:
		buf := make([]byte, 7)
		binary.BigEndian.PutUint32(cache32bits, mh.timestampDelta)
		copy(buf[0:3], cache32bits[1:])
		binary.BigEndian.PutUint32(cache32bits, mh.messageLength)
		copy(buf[3:6], cache32bits[1:])
		buf[6] = mh.messageTypeID

		_, err := w.Write(buf)
		return err

	case 2:
		buf := make([]byte, 3)
		binary.BigEndian.PutUint32(cache32bits, mh.timestampDelta)
		copy(buf[0:3], cache32bits[1:])

		_, err := w.Write(buf)
		return err

	case 3:
		return nil

	default:
		panic("unexpected fmt")
	}
}
