package redis

import (
	"os"
)

type RedisConfig struct {
	Host     string
	Port     string
	Password string
}

var Redis RedisConfig

func init() {

	Redis = RedisConfig{
		Host:     os.Getenv("REDIS_HOST"),
		Port:     os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASS"),
	}

	RedisClient()
}
