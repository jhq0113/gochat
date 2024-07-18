package room

import (
	"sync"

	"github.com/jhq0113/gochat/core"
)

type Rooms[K core.Key] struct {
	mutex sync.RWMutex
	rooms map[int64]*Room[K]
}

func NewRooms[K core.Key]() *Rooms[K] {
	return &Rooms[K]{
		rooms: make(map[int64]*Room[K]),
	}
}

func (rs *Rooms[K]) GetRoom(roomId int64) *Room[K] {
	rs.mutex.RLock()
	defer rs.mutex.RUnlock()

	r, _ := rs.rooms[roomId]
	return r
}

func (rs *Rooms[K]) Range(fn func(roomId int64, room *Room[K])) {
	var (
		roomIds []int64
		rooms   []*Room[K]
	)

	rs.mutex.RLock()

	if len(rs.rooms) > 0 {
		roomIds = make([]int64, 0, len(rs.rooms))
		rooms = make([]*Room[K], 0, len(rs.rooms))
		for roomId, room := range rs.rooms {
			roomIds = append(roomIds, roomId)
			rooms = append(rooms, room)
		}
	}

	rs.mutex.RUnlock()

	if len(roomIds) > 0 {
		for i, room := range rooms {
			fn(roomIds[i], room)
		}
	}
}

func (rs *Rooms[K]) Join(roomId int64, key K) {
	rs.mutex.Lock()
	if _, ok := rs.rooms[roomId]; !ok {
		rs.rooms[roomId] = NewRoom[K](roomId)
	}
	room := rs.rooms[roomId]
	rs.mutex.Unlock()

	room.Join(key)
}

func (rs *Rooms[K]) Leave(roomId int64, key K) {
	rs.mutex.RLock()

	if _, ok := rs.rooms[roomId]; !ok {
		rs.mutex.RUnlock()
		return
	}

	room := rs.rooms[roomId]
	rs.mutex.RUnlock()

	room.Leave(key)
}
