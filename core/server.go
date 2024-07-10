package core

import (
	"errors"
	"net/http"
	"time"

	"github.com/Allenxuxu/gev"
	"github.com/Allenxuxu/gev/plugins/websocket/ws"
	"github.com/Allenxuxu/toolkit/convert"
)

type Server struct {
	session   *Session
	onConnect func(c *Conn)
	onHeader  func(c *Conn, key, value []byte) error
	onRequest func(c *Conn, uri []byte) error
	onAccept  func(c *Conn, uri string, headers http.Header) error
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

func (s *Server) OnHeader(c *gev.Connection, key, value []byte) error {
	var (
		id      = Id(c)
		conn, _ = s.session.Get(id)
	)

	if conn == nil || id == 0 {
		c.Close()
		return errors.New("conn not exists")
	}

	if s.onHeader != nil {
		if err := s.onHeader(conn, key, value); err != nil {
			return err
		}
	}

	var header http.Header
	_header, ok := c.Get(CtxHeader)
	if ok {
		header = _header.(http.Header)
	} else {
		header = make(http.Header)
	}

	header.Set(convert.BytesToString(key), convert.BytesToString(value))
	c.Set(CtxHeader, header)

	return nil
}

func (s *Server) OnRequest(c *gev.Connection, uri []byte) error {
	var (
		id      = Id(c)
		conn, _ = s.session.Get(id)
	)

	if conn == nil || id == 0 {
		c.Close()
		return errors.New("conn not exists")
	}

	if s.onRequest != nil {
		if err := s.onRequest(conn, uri); err != nil {
			return err
		}
	}

	c.Set(CtxUri, convert.BytesToString(uri))
	return nil
}

func (s *Server) OnMessage(c *gev.Connection, data []byte) (messageType ws.MessageType, out []byte) {
	var (
		id      = Id(c)
		conn, _ = s.session.Get(id)
	)

	if id == 0 || conn == nil {
		c.Close()
		return
	}

	if conn.acceptedAt == 0 {
		conn.acceptedAt = time.Now().Unix()

		var (
			header http.Header
			uri    string
		)

		_uri, ok := c.Get(CtxUri)
		if ok {
			uri = _uri.(string)
		}

		_header, ok := c.Get(CtxHeader)
		if ok {
			header = _header.(http.Header)
		}

		if err := s.onAccept(conn, uri, header); err != nil {
			conn.Close()
			return
		}
	}

	messageType = ws.MessageBinary
	return
}

func (s *Server) OnClose(c *gev.Connection) {
	var (
		id      = Id(c)
		conn, _ = s.session.Get(id)
	)

	defer func() {
		s.session.Remove(id)
		if conn != nil {
			ReleaseConn(conn)
		}
	}()

	if s.onClose != nil && conn != nil {
		s.onClose(conn)
	}
}
