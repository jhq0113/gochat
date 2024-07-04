package core

import "github.com/Allenxuxu/gev"

func Id(c *gev.Connection) (id uint64) {
	if c == nil {
		return
	}

	value, _ := c.Get(CtxId)
	if val, ok := value.(uint64); ok {
		return val
	}

	return
}
