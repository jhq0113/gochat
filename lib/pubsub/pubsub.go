package pubsub

import (
	"context"
	"strings"

	"github.com/jhq0113/gochat/lib/pogo"

	"github.com/Allenxuxu/toolkit/convert"
	"github.com/goccy/go-json"
	"github.com/redis/go-redis/v9"
)

type CmdMsg struct {
	Channel string     `json:"c"`
	TraceId string     `json:"tid"`
	AppId   uint8      `json:"aid"`
	Id      int64      `json:"id"`
	Param   pogo.Param `json:"p"`
}

func (cm *CmdMsg) Marshal() []byte {
	data, _ := json.Marshal(cm)
	return data
}

func SubscribeRedis(red *redis.Client, ctx context.Context, channels ...string) <-chan CmdMsg {
	var (
		ps  = red.Subscribe(ctx, channels...)
		ch  = make(chan CmdMsg, 1)
		msg *redis.Message
		err error
	)

	go func() {
		defer ps.Close()

		for {
			msg, err = ps.ReceiveMessage(ctx)
			if err != nil {
				if strings.HasPrefix(err.Error(), "redis: unknown message") {
					continue
				}
				close(ch)
				return
			}

			var cmdMsg CmdMsg
			if err = json.Unmarshal(convert.StringToBytes(msg.Payload), &cmdMsg); err != nil {
				continue
			}

			cmdMsg.Channel = msg.Channel
			ch <- cmdMsg
		}
	}()

	return ch
}
