package core

import (
	"hash/crc32"
	"sync"

	"github.com/Allenxuxu/toolkit/convert"
	"go.uber.org/atomic"
)

type Key interface {
	~uint64 | ~int64 | ~string
}

func KeyEmpty(key any) bool {
	switch value := key.(type) {
	case uint64:
		return value == 0
	case int64:
		return value == 0
	}

	value, _ := key.(string)
	return value == ""
}

type Session[K Key] struct {
	length atomic.Int64
	bucket []*session[K]
}

func NewSession[K Key](bucketSize int) *Session[K] {
	s := &Session[K]{
		bucket: make([]*session[K], bucketSize),
	}

	for i := 0; i < bucketSize; i++ {
		s.bucket[i] = &session[K]{
			conns: make(map[K]*Conn),
		}
	}

	return s
}

func (s *Session[K]) Len() int64 {
	return s.length.Load()
}

func (s *Session[K]) index(key any) int {
	switch value := key.(type) {
	case uint64:
		return int(value % uint64(len(s.bucket)))
	case int64:
		return int(value % int64(len(s.bucket)))
	}

	value, _ := key.(string)
	return int(crc32.ChecksumIEEE(convert.StringToBytes(value))) % len(s.bucket)
}

func (s *Session[K]) Set(key K, c *Conn) (isNew bool) {
	isNew = s.bucket[s.index(key)].set(key, c)
	if isNew {
		s.length.Inc()
	}
	return
}

func (s *Session[K]) Get(key K) (c *Conn, exists bool) {
	if KeyEmpty(key) {
		return
	}

	return s.bucket[s.index(key)].get(key)
}

func (s *Session[K]) Range(fn func(c *Conn)) {
	for _, bucket := range s.bucket {
		bucket.Range(fn)
	}
}

func (s *Session[K]) Remove(key K) (delNum int) {
	if KeyEmpty(key) {
		return
	}

	delNum = s.bucket[s.index(key)].remove(key)
	if delNum > 0 {
		s.length.Sub(int64(delNum))
	}
	return
}

type session[K Key] struct {
	mutex sync.RWMutex
	conns map[K]*Conn
}

func (s *session[K]) set(key K, c *Conn) (isNew bool) {
	if KeyEmpty(key) {
		return
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()
	_, exists := s.conns[key]
	s.conns[key] = c

	return !exists
}

func (s *session[K]) get(key K) (c *Conn, exists bool) {
	if KeyEmpty(key) {
		return
	}

	s.mutex.RLock()
	defer s.mutex.RUnlock()

	c, exists = s.conns[key]
	return
}

func (s *session[K]) Range(fn func(c *Conn)) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for _, conn := range s.conns {
		fn(conn)
	}
}

func (s *session[K]) remove(key K) (delNum int) {
	if KeyEmpty(key) {
		return
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	_, exists := s.conns[key]
	delete(s.conns, key)

	if exists {
		return 1
	}

	return
}
