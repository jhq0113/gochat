package actions

import (
	"github.com/jhq0113/gochat/core"
	"github.com/jhq0113/gochat/lib/constants"
	"github.com/jhq0113/gochat/lib/pogo"
	"github.com/jhq0113/gochat/lib/sessions"
)

func Join(c *core.Conn, event *pogo.Event) {
	roomId := event.Data.Int64("room", 0)
	if roomId > 0 {
		sessions.JoinRoom(roomId, c)

		msg := pogo.AcqEventWithId(constants.JoinOk).WithData(pogo.Param{
			"roomId": roomId,
		})
		_ = c.SendTextAsync(msg.Marshal())
	}

	return
}
