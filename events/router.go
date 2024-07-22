package events

import (
	"context"

	"github.com/jhq0113/gochat/global"
	"github.com/jhq0113/gochat/lib/constants"
	"github.com/jhq0113/gochat/lib/pubsub"

	"github.com/redis/go-redis/v9"
)

func LoadRouter() {
	er := pubsub.NewEventRouter()

	er.On(constants.RoomDestroy, RoomDestroy)

	go func() {
		var red = redis.NewClient(global.RedisOption)

		ch := pubsub.SubscribeRedis(red, context.Background(), "gochat_events")
		for msg := range ch {
			if msg.Id == 0 {
				continue
			}

			go er.Handler(msg)
		}
	}()
}
