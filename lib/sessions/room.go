package sessions

import (
	"sync"

	"github.com/jhq0113/gochat/core"
)

const (
	CtxRoomId = `ctx:room`
)

var (
	globalRooms = NewRooms[uint64]()
)

func GetRoom(roomId int64) *Room[uint64] {
	return globalRooms.GetRoom(roomId)
}

func JoinRoom(roomId int64, c *core.Conn) {
	c.Set(CtxRoomId, roomId)

	globalRooms.Join(roomId, c.Id())
}

func LeaveRoom(roomId int64, c *core.Conn) {
	if roomId == RoomId(c) {
		c.Delete(CtxRoomId)
	}

	globalRooms.Leave(roomId, c.Id())
}

func RoomId(c *core.Conn) int64 {
	value, _ := c.Get(CtxRoomId)
	val, _ := value.(int64)
	return val
}

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

type Room[K core.Key] struct {
	mutex  sync.RWMutex
	roomId int64
	set    map[K]struct{}
}

func NewRoom[K core.Key](id int64) *Room[K] {
	return &Room[K]{
		roomId: id,
		set:    make(map[K]struct{}),
	}
}

func (r *Room[K]) SomeRange(size int, fn func(keys []K)) {
	if size < 1 {
		return
	}

	r.mutex.RLock()
	defer r.mutex.RUnlock()

	keys := make([]K, 0, size)
	for key, _ := range r.set {
		keys = append(keys, key)
		if len(keys) >= size {
			fn(keys)
			keys = keys[:0]
		}
	}

	if len(keys) > 0 {
		fn(keys)
	}
}

func (r *Room[K]) Id() int64 {
	return r.roomId
}

func (r *Room[K]) Len() int {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	return len(r.set)
}

func (r *Room[K]) Join(key K) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.set[key] = struct{}{}
}

func (r *Room[K]) Leave(key K) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	delete(r.set, key)

	if len(r.set) == 0 {
		r.set = make(map[K]struct{})
	}
}
