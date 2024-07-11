package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jhq0113/gochat/core"

	"github.com/Allenxuxu/gev"
	"github.com/Allenxuxu/gev/plugins/websocket/ws"
	"github.com/Allenxuxu/toolkit/convert"
)

func main() {
	s, err := core.NewServer(
		func(c *core.Conn, data []byte) (messageType ws.MessageType, out []byte) {
			fmt.Printf("receive msg: %s\n", data)
			return
		},
		gev.IdleTime(time.Second*60),
		gev.Network("tcp"),
		gev.Address(":8838"),
		gev.NumLoops(4),
		gev.LoadBalance(gev.LeastConnection()),
	)

	if err != nil {
		panic(err)
	}

	s.BindAcceptHandler(func(c *core.Conn, uri string, headers http.Header) error {
		fmt.Printf("accept id: %d uri: %s headers: %+v\n", c.Id(), uri, headers)
		return nil
	})

	s.RunEvery(time.Second, func() {
		s.Range(func(c *core.Conn) {
			if err = c.SendTextAsync(convert.StringToBytes(`Hello World!`)); err != nil {
				fmt.Printf("send msg error: %v\n", err)
			}
		})
		fmt.Printf("conn total: %d\n", s.ConnCount())
	})

	defer s.Stop()

	s.Start()
}
