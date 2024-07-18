package events

import (
	"context"

	"github.com/jhq0113/gochat/lib/constants"
	"github.com/jhq0113/gochat/lib/pubsub"

	"github.com/redis/go-redis/v9"
)

func LoadRouter() {
	er := pubsub.NewEventRouter()

	er.On(constants.RoomDestroy, RoomDestroy)

	go func() {
		var red = redis.NewClient(&redis.Options{
			Network:  "tcp",
			Addr:     "127.0.0.1:6379",
			Password: "12345678",
		})

		ch := pubsub.SubscribeRedis(red, context.Background(), "gochat_events")
		for msg := range ch {
			if msg.Id == 0 {
				continue
			}

			go er.Handler(msg)
		}
	}()
}
