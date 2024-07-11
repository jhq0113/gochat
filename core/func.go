package core

import (
	"unsafe"

	"github.com/Allenxuxu/gev"
)

func Id(c *gev.Connection) (id uint64) {
	if c == nil {
		return
	}

	return uint64(uintptr(unsafe.Pointer(c)))
}

func GetConn(ser *Server, c *gev.Connection) (conn *Conn, err error) {
	id := Id(c)
	if id == 0 {
		return nil, ErrConnNotFound
	}

	conn, _ = ser.session.Get(id)
	if conn == nil {
		return conn, ErrConnNotFound
	}

	return conn, nil
}
