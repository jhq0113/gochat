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
		s.bucket[i] = &session{
			conns: make(map[uint64]*Conn),
		}
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

	id := c.Id()
	if id == 0 {
		return
	}

	isNew = s.bucket[s.index(id)].set(id, c)
	if isNew {
		s.length.Inc()
	}
	return
}

func (s *Session) GetConn(conn *gev.Connection) (c *Conn, exists bool) {
	if conn == nil {
		return
	}

	id := Id(conn)
	if id == 0 {
		return
	}

	return s.bucket[s.index(id)].get(id)
}

func (s *Session) Get(id uint64) (c *Conn, exists bool) {
	if id == 0 {
		return
	}

	return s.bucket[s.index(id)].get(id)
}

func (s *Session) Range(fn func(c *Conn)) {
	for _, bucket := range s.bucket {
		bucket.Range(fn)
	}
}

func (s *Session) Remove(id uint64) (delNum int) {
	if id == 0 {
		return
	}

	delNum = s.bucket[s.index(id)].remove(id)
	if delNum > 0 {
		s.length.Sub(int64(delNum))
	}
	return
}

func (s *Session) RemoveConn(c *Conn) (delNum int) {
	if c == nil || c.Connection == nil {
		return
	}

	return s.Remove(c.Id())
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
	if id == 0 {
		return
	}

	s.mutex.RLock()
	defer s.mutex.RUnlock()

	c, exists = s.conns[id]
	return
}

func (s *session) Range(fn func(c *Conn)) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for _, conn := range s.conns {
		fn(conn)
	}
}

func (s *session) remove(id uint64) (delNum int) {
	if id == 0 {
		return
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	_, exists := s.conns[id]
	delete(s.conns, id)

	if exists {
		return 1
	}

	return
}
