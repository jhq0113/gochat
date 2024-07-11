package actions

import (
	"fmt"

	"github.com/jhq0113/gochat/core"
	"github.com/jhq0113/gochat/lib/pogo"

	"github.com/Allenxuxu/gev/plugins/websocket/ws"
)

func Login(c *core.Conn, event *pogo.Event) (messageType ws.MessageType, out []byte) {
	fmt.Printf("login: %v\n", event)
	return
}
