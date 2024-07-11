package main

import (
	"fmt"
	"time"

	"github.com/jhq0113/gochat/core"

	"github.com/Allenxuxu/gev"
	"github.com/Allenxuxu/gev/plugins/websocket/ws"
)

func main() {
	s, err := core.NewServer(
		func(c *core.Conn, data []byte) (messageType ws.MessageType, out []byte) {
			return
		},
		gev.Network("tcp"),
		gev.Address(":8838"),
		gev.NumLoops(4),
		gev.LoadBalance(gev.LeastConnection()),
	)

	if err != nil {
		panic(err)
	}

	s.RunEvery(time.Second, func() {
		fmt.Printf("conn count: %d\n", s.ConnCount())
	})

	defer s.Stop()

	s.Start()
}
