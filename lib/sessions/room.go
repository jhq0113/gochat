package sessions

import (
	"github.com/jhq0113/gochat/core"
	"github.com/jhq0113/gochat/lib/sessions/room"
)

const (
	CtxRoomId = `ctx:room`
)

var (
	globalRooms = room.NewRooms[uint64]()
)

func GetRoom(roomId int64) *room.Room[uint64] {
	return globalRooms.GetRoom(roomId)
}

func RangeRooms(fn func(roomId int64, room *room.Room[uint64])) {
	globalRooms.Range(fn)
}

func JoinRoom(roomId int64, c *core.Conn) {
	if roomId < 1 {
		return
	}

	rId := RoomId(c)
	if rId == roomId {
		return
	}

	LeaveRoom(c)

	c.Set(CtxRoomId, roomId)

	globalRooms.Join(roomId, c.Id())
}

func LeaveRoom(c *core.Conn) {
	defer c.Delete(CtxRoomId)

	roomId := RoomId(c)
	if roomId < 1 {
		return
	}

	globalRooms.Leave(roomId, c.Id())
}

func RoomId(c *core.Conn) int64 {
	value, _ := c.Get(CtxRoomId)
	val, _ := value.(int64)
	return val
}
