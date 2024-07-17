package sessions

import "github.com/jhq0113/gochat/core"

func Leave(c *core.Conn) {
	LeaveRoom(c)
	LoginOut(c)
}
