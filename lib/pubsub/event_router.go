package pubsub

type Handler func(msg CmdMsg)

type EventRouter struct {
	handlers map[int64]Handler
}

func NewEventRouter() *EventRouter {
	return &EventRouter{handlers: make(map[int64]Handler)}
}

func (r *EventRouter) On(id int64, handler Handler) {
	r.handlers[id] = handler
}

func (r *EventRouter) Handler(msg CmdMsg) {
	if handler, ok := r.handlers[msg.Id]; ok {
		handler(msg)
	}
}
