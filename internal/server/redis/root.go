package redis

import (
	"bookbox-backend/internal/model"
	"os"
)

var Redis model.Redis

func InitRedis() error {

	Redis = model.Redis{
		Host:     os.Getenv("REDIS_HOST"),
		Port:     os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASS"),
	}

	return nil
}
