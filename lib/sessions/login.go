package sessions

import (
	"github.com/jhq0113/gochat/core"
)

const (
	CtxUserId = `ctx:uid`
)

var (
	login = core.NewSession[int64](32)
)

func LoginCount() int64 {
	return login.Len()
}

func UserId(c *core.Conn) int64 {
	userId, _ := c.Get(CtxUserId)
	uid, _ := userId.(int64)
	return uid
}

func Login(userId int64, c *core.Conn) {
	c.Set(CtxUserId, userId)
	login.Set(userId, c)
}

func LoginOut(c *core.Conn) {
	userId := UserId(c)
	login.Remove(userId)
}
