package actions

import (
	"github.com/jhq0113/gochat/core"
	"github.com/jhq0113/gochat/lib/pogo"
	"github.com/jhq0113/gochat/lib/sessions"

	"github.com/Allenxuxu/gev/plugins/websocket/ws"
)

func Login(c *core.Conn, event *pogo.Event) (messageType ws.MessageType, out []byte) {
	userId := event.Data.Int64("userId", 0)
	if userId%2 != 0 {
		sessions.Login(userId, c)
	}

	return
}
