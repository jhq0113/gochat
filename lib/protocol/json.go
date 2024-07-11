package protocol

import (
	"github.com/jhq0113/gochat/core"
	"github.com/jhq0113/gochat/lib/pogo"

	"github.com/Allenxuxu/gev/plugins/websocket/ws"
	"github.com/goccy/go-json"
)

type Json struct {
	*Protocol
}

func NewJson(handler Handler) *Json {
	return &Json{
		Protocol: &Protocol{
			handler: handler,
		},
	}
}

func (j *Json) Handler(c *core.Conn, data []byte) (messageType ws.MessageType, out []byte) {
	event := pogo.AcqEvent()
	defer event.Close()

	if err := json.Unmarshal(data, event); err != nil {
		c.Close()
		return
	}

	return j.handler(c, event)
}
