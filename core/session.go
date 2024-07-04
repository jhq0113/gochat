package core

import (
	"sync"

	"github.com/Allenxuxu/gev"
	"go.uber.org/atomic"
)

type Session struct {
	length atomic.Int64
	bucket []*session
}

func NewSession(bucketSize int) *Session {
	s := &Session{
		bucket: make([]*session, bucketSize),
	}

	for i := 0; i < bucketSize; i++ {
		s.bucket[i] = &session{}
	}

	return s
}

func (s *Session) index(id uint64) int {
	return int(id % uint64(len(s.bucket)))
}

func (s *Session) Set(c *Conn) (isNew bool) {
	if c == nil {
		return
	}

	isNew = s.bucket[s.index(c.id)].set(c.id, c)
	if isNew {
		s.length.Inc()
	}
	return
}

func (s *Session) GetConn(gc *gev.Connection) (c *Conn, exists bool) {
	if gc == nil {
		return
	}

	id := Id(gc)
	if id == 0 {
		return
	}

	return s.bucket[s.index(id)].get(id)
}

func (s *Session) Get(id uint64) (c *Conn, exists bool) {
	return s.bucket[s.index(id)].get(id)
}

func (s *Session) Remove(id uint64) (delNum int) {
	delNum = s.bucket[s.index(id)].remove(id)
	if delNum > 0 {
		s.length.Sub(int64(delNum))
	}
	return
}

func (s *Session) RemoveConn(c *Conn) (delNum int) {
	if c == nil {
		return
	}

	return s.Remove(c.id)
}

type session struct {
	mutex sync.RWMutex
	conns map[uint64]*Conn
}

func (s *session) set(id uint64, c *Conn) (isNew bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	_, exists := s.conns[id]
	s.conns[id] = c

	return !exists
}

func (s *session) get(id uint64) (c *Conn, exists bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	c, exists = s.conns[id]
	return
}

func (s *session) remove(id uint64) (delNum int) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	_, exists := s.conns[id]
	delete(s.conns, id)

	if exists {
		return 1
	}

	return
}
