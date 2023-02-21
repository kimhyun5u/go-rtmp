package rtmp

import "rtmp/message"

type Handler struct {
	OnMessage func(msg message.Message, r Stream)
}
