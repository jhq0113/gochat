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

func Login(userId int64, c *core.Conn) {
	c.Set(CtxUserId, userId)
	login.Set(userId, c)
}

func LoginOut(c *core.Conn) {
	userId, _ := c.Get(CtxUserId)
	uid, _ := userId.(int64)
	login.Remove(uid)
}
