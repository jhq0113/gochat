package actions

import (
	"github.com/jhq0113/gochat/core"
	"github.com/jhq0113/gochat/lib/constants"
	"github.com/jhq0113/gochat/lib/pogo"
	"github.com/jhq0113/gochat/lib/sessions"
)

func Login(c *core.Conn, event *pogo.Event) {
	userId := event.Data.Int64("userId", 0)
	if userId > 0 {
		sessions.Login(userId, c)
		msg := pogo.AcqEventWithId(constants.LoginOk).WithData(pogo.Param{
			"userId": userId,
		})

		_ = c.SendTextAsync(msg.Marshal())
	}

	return
}
