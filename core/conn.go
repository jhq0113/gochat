package core

import (
	"sync"

	"github.com/Allenxuxu/gev"
	"go.uber.org/atomic"
)

var (
	connId atomic.Uint64

	connPool = sync.Pool{
		New: func() any {
			return &Conn{}
		},
	}

	AcquireConn = func(c *gev.Connection) *Conn {
		conn := connPool.Get().(*Conn)
		conn.id = connId.Inc()
		conn.Set(CtxId, conn.id)
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
	id uint64
}

func (c *Conn) reset() {
	c.id = 0
	return
}
