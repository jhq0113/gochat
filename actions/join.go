package actions

import (
	"github.com/jhq0113/gochat/core"
	"github.com/jhq0113/gochat/lib/pogo"
	"github.com/jhq0113/gochat/lib/sessions"

	"github.com/Allenxuxu/gev/plugins/websocket/ws"
)

func Join(c *core.Conn, event *pogo.Event) (messageType ws.MessageType, out []byte) {
	roomId := event.Data.Int64("room", 0)
	sessions.JoinRoom(roomId, c)
	return
}
