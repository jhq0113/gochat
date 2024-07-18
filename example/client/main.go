package main

import (
	"context"
	"math/rand"
	"strconv"

	"github.com/jhq0113/gochat/lib/constants"
	"github.com/jhq0113/gochat/lib/pogo"
	"github.com/jhq0113/gochat/lib/pubsub"

	"github.com/redis/go-redis/v9"
)

func main() {
	red := redis.NewClient(&redis.Options{
		Network:  "tcp",
		Addr:     "127.0.0.1:6379",
		Password: "12345678",
	})

	msg := pubsub.CmdMsg{
		TraceId: strconv.FormatInt(rand.Int63(), 10),
		AppId:   0,
		Id:      constants.RoomDestroy,
		Param: pogo.Param{
			"room": 6007,
		},
	}

	red.Publish(context.Background(), "gochat_events", msg.Marshal())
}
