package router

import (
	"github.com/jhq0113/gochat/core"
	"github.com/jhq0113/gochat/lib/pogo"
	"github.com/jhq0113/gochat/lib/protocol"
)

type Router struct {
	handlers map[int64]protocol.Handler
}

func NewRouter() *Router {
	return &Router{
		handlers: make(map[int64]protocol.Handler),
	}
}

func (r *Router) Add(id int64, handler protocol.Handler) {
	r.handlers[id] = handler
}

func (r *Router) Handler(c *core.Conn, event *pogo.Event) {
	h, ok := r.handlers[event.Id]
	if !ok {
		return
	}

	h(c, event)
	return
}
