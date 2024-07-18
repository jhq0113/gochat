package events

import (
	"fmt"

	"github.com/jhq0113/gochat/lib/pubsub"
	"github.com/jhq0113/gochat/lib/sessions"
)

func RoomDestroy(msg pubsub.CmdMsg) {
	roomId := msg.Param.Int64("room", 0)
	if roomId > 0 {
		fmt.Printf("destroy room: %d\n", roomId)
		sessions.RoomDestroy(roomId)
	}
}
