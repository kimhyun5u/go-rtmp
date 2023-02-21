package rtmp

import "rtmp/message"

type Stream interface {
	Write(msg message.Message) error
}
