package redis

import (
	"fmt"

	"github.com/go-redis/redis"
)

var Client *redis.Client

func RedisClient() error {

	Client = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf(
			"%s:%s",
			Redis.Host,
			Redis.Port,
		),
		Password: Redis.Password, // no password set
		DB:       0,              // use default DB
	})

	return nil
}
