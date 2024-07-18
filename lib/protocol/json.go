package protocol

import (
	"net/http"

	"github.com/jhq0113/gochat/core"
	"github.com/jhq0113/gochat/lib/pogo"

	"github.com/Allenxuxu/gev/plugins/websocket/ws"
	"github.com/goccy/go-json"
)

type Json struct {
	handler Handler
}

func NewJson(handler Handler) core.Protocol {
	return &Json{
		handler: handler,
	}
}

func (j *Json) Accept(c *core.Conn, uri string, headers http.Header) error {
	return nil
}

func (j *Json) Handler(c *core.Conn, data []byte) (messageType ws.MessageType, out []byte) {
	event := pogo.AcqEvent()
	defer event.Close()

	if err := json.Unmarshal(data, event); err != nil {
		c.Close()
		return
	}

	j.handler(c, event)
	return
}

func (j *Json) Pack(c *core.Conn, data []byte) ([]byte, error) {
	return data, nil
}
