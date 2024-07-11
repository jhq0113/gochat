package protocol

import (
	"net/http"

	"github.com/Allenxuxu/gev/plugins/websocket/ws"
	"github.com/jhq0113/gochat/core"
)

type Protocol interface {
	Accept(c *core.Conn, uri string, headers http.Header) error
	Handler(c *core.Conn, data []byte) (messageType ws.MessageType, out []byte)
	Pack(c *core.Conn, data []byte) ([]byte, error)
}
