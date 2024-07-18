package room

import (
	"sync"

	"github.com/jhq0113/gochat/core"
)

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

func (r *Room[K]) Leave(key K) (restNum int) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	delete(r.set, key)

	if len(r.set) == 0 {
		r.set = make(map[K]struct{})
	}

	return len(r.set)
}
