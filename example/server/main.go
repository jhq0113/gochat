package main

import (
	"fmt"
	"github.com/jhq0113/gochat/lib/sessions"
	"time"

	"github.com/jhq0113/gochat/actions"
	"github.com/jhq0113/gochat/core"
	"github.com/jhq0113/gochat/lib/constants"
	"github.com/jhq0113/gochat/lib/pogo"
	"github.com/jhq0113/gochat/lib/protocol"
	"github.com/jhq0113/gochat/lib/utils"

	"github.com/Allenxuxu/gev"
)

func main() {
	privateKey := `-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCAR9fLwylDuY+3
b1oXfyBVNphtQYkx+Tu95j9xYqXldJqWuhM19BaJ2rmiyTDzfErU8jXmU+VUd05S
og40eGuOPzN9qryLtiM681pzdSLEl9EJuF9qPnY2J5ocsx3Ngd9lmMwo6QOylqvV
Sv9vldYVzzLmiIWzGPlhLE9zSaBx5gKAaHgBSbT8ov5vpc2YgPupZ1A9a4/7feTz
0WOrs2m8CS7pkn3j6SdebFu5Kn8Fw8jxARhYl6uEwDmw2wI+FGUS0G33onn5xDT8
AFH3gUDMvx2p7N6lrOmq1PVHp4os39Bg9Kk18rqTnLGVaa3HanbYnqawHmAj431G
5M3XBNuFAgMBAAECggEAYNVaiV69xHjLit2MdVYqUwjOURL6fHW16ihnVbDwp9+4
XfUCUWT6qn8oSEB1Kid12X8ovsAvye3xeqEH9gW8svj0ZnU+tHv95h8S9TrDzwEG
875wkCOsWj7Ur7tp0Nqrhuarnmoet23TMkPVxCaoH3hz5N/DhSlh/j9TjiRS8AZk
oCtaSTujZxKRW9J8WmOhr+BCn1TruOYrM+CWwY4uyAJ2C8GdW19pNLEPCGbTi3wd
Qdw5kx06nNr/w/zF7aLxckTh5ELYrrq65KA8V/4Vc6eaF8i8KAtAV8xEOD8SgAm6
d2nWJNuLaibIDKAZeQ0eEZNx/NhPdt99O8pRalQfQQKBgQDq3CnbFJMDWtR1Srj/
ikuODPSf+wHli+Pe14NPTTfyz/rruxPxUy3hPGPm93xd3A8G8CPVDBZ/keUbFe5b
GMx5gtwVzS7GI+JWjmDsi/b1EUir6+bpv9jbI6k+fKs7D9TmI/1DQvf+VfrS+9T4
RQei5w7mdPPXNMSFoMGFhdgfVQKBgQCL08pQNg5F8R15aolbDXK2Wuj7Dq7jWpu3
KhN6ORBhuFvU5XBhSfvLfvxVRA21Gim0DBX503U22KMysD3Zr26PQQmwfNnQ63V4
XWOHl8rehUMXZkakLD6FILkdLDO7Xl4TZbnwpkaXLkcecaGNX4lGkjAw6jny6XcX
wcCqjmzrcQKBgQDpWtonuNid47jn1dfc6C9MNCk3b/Khfo5qU15ABCyMEQRzBs24
4XnbquJkkhC7PbScoywnOjx8gpzOfcr4Lrq0HUbKTrWj3/G7KPq67hLxyzuWvu0P
4jP5AQTfdoW4SHG0PZweIcNArXNsARbJm+ULgmM2rou9j8uYLnM6VRO5hQKBgFx4
gspLBWNx0okyUIYbvaolMwCdNEF88Y/PTrQ8ur21W13PImPktpVcdGm7KMmE9OFy
QLJICrNrz3m9HhoxL4+jdlH2L2I/5R7lu+W3F93TCYyXAc6ex/XoryNA0TFvdg2j
77TbccGXREc45JsG/FTkZuRiclJX3X+jjdP9fsCxAoGAHwtH1TU5BxhN0615/SfC
mxjoXbkzqzWtSFb6kr5MiTYAAsfXpHrzVupJvCHvoJ6zr2bZNl2C8Kj5slc091UH
uHhjP8AH1rAw5ipncs6eYNC/wxz7gK9Yp9Vs7q0uhSDAXExPSK6pOyu+D0sY+rF4
EgnfyWP6+Suqh8gJfa2hZug=
-----END PRIVATE KEY-----`
	publicKey := `-----BEGIN PUBLIC KEY-----
    MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAgEfXy8MpQ7mPt29aF38g
VTaYbUGJMfk7veY/cWKl5XSalroTNfQWidq5oskw83xK1PI15lPlVHdOUqIONHhr
jj8zfaq8i7YjOvNac3UixJfRCbhfaj52NieaHLMdzYHfZZjMKOkDspar1Ur/b5XW
Fc8y5oiFsxj5YSxPc0mgceYCgGh4AUm0/KL+b6XNmID7qWdQPWuP+33k89Fjq7Np
vAku6ZJ94+knXmxbuSp/BcPI8QEYWJerhMA5sNsCPhRlEtBt96J5+cQ0/ABR94FA
zL8dqezepazpqtT1R6eKLN/QYPSpNfK6k5yxlWmtx2p22J6msB5gI+N9RuTN1wTb
hQIDAQAB
-----END PUBLIC KEY-----`

	rsa, err := utils.NewRsa(publicKey, privateKey)
	if err != nil {
		panic(err)
	}

	var (
		router = actions.LoadRouter()
		//proto  = protocol.NewJson(router.Handler)
		proto = protocol.NewV1(router.Handler, rsa)
	)

	s, err := core.NewServer(
		proto,
		gev.IdleTime(time.Second*60),
		gev.Network("tcp"),
		gev.Address(":8838"),
		gev.NumLoops(4),
		gev.LoadBalance(gev.LeastConnection()),
	)

	if err != nil {
		panic(err)
	}

	s.BindCloseHandler(func(c *core.Conn) {
		sessions.LoginOut(c)
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
		fmt.Printf("conn total: %d login count: %d\n", s.ConnCount(), sessions.LoginCount())
	})

	defer s.Stop()

	s.Start()
}
