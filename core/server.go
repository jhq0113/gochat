package core

import (
	"github.com/Allenxuxu/gev"
)

type Server struct {
	session   *Session
	onConnect func(c *Conn)
	onClose   func(c *Conn)
}

func (s *Server) ConnCount() int64 {
	return s.session.length.Load()
}

func (s *Server) OnConnect(c *gev.Connection) {
	conn := AcquireConn(c)
	if s.onConnect != nil {
		s.onConnect(conn)
	}

	s.session.Set(conn)
}

func (s *Server) OnClose(c *gev.Connection) {
	id := Id(c)

	defer s.session.Remove(id)

	if s.onClose != nil {
		conn, _ := s.session.Get(id)
		s.onClose(conn)
	}
}
