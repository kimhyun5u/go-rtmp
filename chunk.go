package rtmp

import "rtmp/message"

type ChunkStreamIO struct {
	streamID int
	f        func(streamID int, msg message.Message) error
}

func (w *ChunkStreamIO) Write(msg message.Message) error {
	return w.f(w.streamID, msg)
}

type ChunkStreamLayer struct {
}
