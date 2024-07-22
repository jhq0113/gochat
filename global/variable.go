package global

import "github.com/redis/go-redis/v9"

var (
	RedisOption = &redis.Options{
		Network:  "tcp",
		Addr:     "127.0.0.1:6379",
		Password: "12345678",
	}
)
