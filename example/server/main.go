package main

import (
	"fmt"
	"github.com/jhq0113/gochat/lib/constants"
	"github.com/jhq0113/gochat/lib/pogo"
	"net/http"
	"time"

	"github.com/jhq0113/gochat/actions"
	"github.com/jhq0113/gochat/core"
	"github.com/jhq0113/gochat/lib/protocol"

	"github.com/Allenxuxu/gev"
)

func main() {
	var (
		router = actions.LoadRouter()
		proto  = protocol.NewJson(router.Handler)
		s, err = core.NewServer(
			proto.Handler,
			gev.IdleTime(time.Second*60),
			gev.Network("tcp"),
			gev.Address(":8838"),
			gev.NumLoops(4),
			gev.LoadBalance(gev.LeastConnection()),
		)
	)

	if err != nil {
		panic(err)
	}

	s.BindAcceptHandler(func(c *core.Conn, uri string, headers http.Header) error {
		fmt.Printf("accept id: %d uri: %s headers: %+v\n", c.Id(), uri, headers)
		return nil
	})

	s.RunEvery(time.Second, func() {
		event := pogo.AcqEventWithId(constants.Login)
		event.WithData(pogo.Param{
			"code": 100,
			"data": pogo.Param{},
			"msg":  "ok",
		})

		msg := event.Marshal()
		event.Close()

		s.Range(func(c *core.Conn) {
			if err = c.SendTextAsync(msg); err != nil {
				fmt.Printf("send msg error: %v\n", err)
			}
		})
		fmt.Printf("conn total: %d\n", s.ConnCount())
	})

	defer s.Stop()

	s.Start()
}
