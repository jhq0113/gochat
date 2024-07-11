package protocol

import (
	"github.com/jhq0113/gochat/core"
	"github.com/jhq0113/gochat/lib/pogo"

	"github.com/Allenxuxu/gev/plugins/websocket/ws"
)

type Handler func(c *core.Conn, event *pogo.Event) (messageType ws.MessageType, out []byte)

type Router struct {
	handlers map[int64]Handler
}

func NewRouter() *Router {
	return &Router{
		handlers: make(map[int64]Handler),
	}
}

func (r *Router) Add(id int64, handler Handler) {
	r.handlers[id] = handler
}

func (r *Router) Handler(c *core.Conn, event *pogo.Event) (messageType ws.MessageType, out []byte) {
	h, ok := r.handlers[event.Id]
	if !ok {
		return
	}

	return h(c, event)
}
