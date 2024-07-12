package core

import (
	"sync"

	"github.com/Allenxuxu/gev"
	"github.com/Allenxuxu/gev/plugins/websocket/ws"
	"github.com/Allenxuxu/gev/plugins/websocket/ws/util"
	"github.com/gobwas/pool/pbytes"
)

var (
	connPool = sync.Pool{
		New: func() any {
			return &Conn{}
		},
	}

	AcquireConn = func(c *gev.Connection) *Conn {
		conn := connPool.Get().(*Conn)
		conn.Connection = c
		return conn
	}

	ReleaseConn = func(conn *Conn) {
		conn.reset()
		connPool.Put(conn)
	}
)

type Conn struct {
	*gev.Connection
	protocol   Protocol
	acceptedAt int64
}

func (c *Conn) Id() uint64 {
	return Id(c.Connection)
}

func (c *Conn) SendText(text []byte) error {
	text, err := c.protocol.Pack(c, text)
	if err != nil {
		return err
	}

	msg, err := util.PackData(ws.MessageText, text)
	if err != nil {
		return err
	}

	if c.Connection == nil {
		return gev.ErrConnectionClosed
	}

	return c.Connection.Send(msg)
}

func (c *Conn) SendTextAsync(text []byte) error {
	text, err := c.protocol.Pack(c, text)
	if err != nil {
		return err
	}

	msg, err := util.PackData(ws.MessageText, text)
	if err != nil {
		return err
	}

	if len(msg) == 0 {
		if c.Connection != nil {
			return c.Connection.Send(msg)
		}

		return gev.ErrConnectionClosed
	}

	b := pbytes.Get(0, len(msg))
	b = append(b, msg...)

	if c.Connection == nil {
		return gev.ErrConnectionClosed
	}

	return c.Connection.Send(b, gev.SendInLoop(func(i interface{}) {
		pbytes.Put(i.([]byte))
	}))
}

func (c *Conn) reset() {
	c.acceptedAt = 0
	c.Connection = nil
	c.protocol = nil
	return
}
