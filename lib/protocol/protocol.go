package protocol

import (
	"github.com/jhq0113/gochat/core"
	"github.com/jhq0113/gochat/lib/pogo"

	"github.com/Allenxuxu/gev/plugins/websocket/ws"
)

type Handler func(c *core.Conn, event *pogo.Event) (messageType ws.MessageType, out []byte)

type Protocol struct {
	handler Handler
}
