package message

type UserCtrlEvent interface {
}

type StreamBegin struct {
	StreamID uint32
}
