package protocol

import (
	"github.com/jhq0113/gochat/core"
	"github.com/jhq0113/gochat/lib/pogo"
)

type Handler func(c *core.Conn, event *pogo.Event)
