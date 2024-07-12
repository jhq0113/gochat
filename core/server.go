package core

import (
	"net/http"
	"time"

	"github.com/Allenxuxu/gev"
	"github.com/Allenxuxu/gev/plugins/websocket"
	"github.com/Allenxuxu/gev/plugins/websocket/ws"
	"github.com/Allenxuxu/toolkit/convert"
	"github.com/RussellLuo/timingwheel"
)

type Server struct {
	server    *gev.Server
	session   *Session[uint64]
	protocol  Protocol
	onConnect func(c *Conn)
	onHeader  func(c *Conn, key, value []byte) error
	onRequest func(c *Conn, uri []byte) error
	onAccept  func(c *Conn, uri string, headers http.Header) error
	onClose   func(c *Conn)
}

func NewServer(protocol Protocol, opts ...gev.Option) (*Server, error) {
	s := &Server{
		session:  NewSession[uint64](32),
		protocol: protocol,
	}

	u := &ws.Upgrader{}
	u.OnHeader = s.OnHeader
	u.OnRequest = s.OnRequest
	u.OnBeforeUpgrade = s.OnBeforeUpgrade

	opts = append(opts, gev.CustomProtocol(newGevProtocol(u)))

	ser, err := gev.NewServer(websocket.NewHandlerWrap(u, s), opts...)
	if err != nil {
		return nil, err
	}

	s.onAccept = protocol.Accept
	s.server = ser

	return s, nil
}

func (s *Server) Start() {
	s.server.Start()
}

func (s *Server) Stop() {
	s.server.Stop()
}

func (s *Server) RunEvery(d time.Duration, fn func()) *timingwheel.Timer {
	return s.server.RunEvery(d, fn)
}

func (s *Server) RunAfter(d time.Duration, fn func()) *timingwheel.Timer {
	return s.server.RunAfter(d, fn)
}

func (s *Server) BindConnectHandler(handler func(c *Conn)) {
	s.onConnect = handler
}

func (s *Server) BindHeaderHandler(handler func(c *Conn, key, value []byte) error) {
	s.onHeader = handler
}

func (s *Server) BindRequestHandler(handler func(c *Conn, uri []byte) error) {
	s.onRequest = handler
}

func (s *Server) BindCloseHandler(handler func(c *Conn)) {
	s.onClose = handler
}

func (s *Server) ConnCount() int64 {
	return s.session.length.Load()
}

func (s *Server) OnConnect(c *gev.Connection) {
	conn := AcquireConn(c)
	conn.protocol = s.protocol

	if s.onConnect != nil {
		s.onConnect(conn)
	}

	s.session.Set(conn.Id(), conn)
}

func (s *Server) OnRequest(c *gev.Connection, uri []byte) error {
	conn, err := GetConn(s, c)
	if err != nil {
		return err
	}

	if s.onRequest != nil {
		if err = s.onRequest(conn, uri); err != nil {
			return err
		}
	}

	c.Set(CtxUri, convert.BytesToString(uri))
	return nil
}

func (s *Server) OnHeader(c *gev.Connection, key, value []byte) error {
	conn, err := GetConn(s, c)
	if err != nil {
		return err
	}

	if s.onHeader != nil {
		if err = s.onHeader(conn, key, value); err != nil {
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

func (s *Server) OnBeforeUpgrade(c *gev.Connection) (header ws.HandshakeHeader, err error) {
	conn, err := GetConn(s, c)
	if err != nil {
		return nil, err
	}

	conn.acceptedAt = time.Now().Unix()

	var (
		headers http.Header
		uri     string
	)

	_uri, ok := c.Get(CtxUri)
	if ok {
		uri = _uri.(string)
	}

	_header, ok := c.Get(CtxHeader)
	if ok {
		headers = _header.(http.Header)
	}

	err = s.onAccept(conn, uri, headers)
	return nil, err
}

func (s *Server) OnMessage(c *gev.Connection, data []byte) (messageType ws.MessageType, out []byte) {
	conn, err := GetConn(s, c)
	if err != nil {
		c.Close()
		return
	}

	return s.protocol.Handler(conn, data)
}

func (s *Server) Range(fn func(c *Conn)) {
	s.session.Range(fn)
}

func (s *Server) OnClose(c *gev.Connection) {
	conn, err := GetConn(s, c)
	if err != nil {
		c.Close()
		return
	}

	defer func() {
		s.session.Remove(conn.Id())
		ReleaseConn(conn)
	}()

	if s.onClose != nil && conn != nil {
		s.onClose(conn)
	}
}
