package datasource

import (
	"application/app"
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Addr     string `koanf:"addr"`
	Password string `koanf:"password"`
	DB       int    `koanf:"db"`
}

func NewRedisConfig() (*RedisConfig, error) {
	return &RedisConfig{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	}, nil
}

type RedisDS struct {
	*redis.Client
}

func NewRedisDS(ctx context.Context, cfg *RedisConfig, controller app.Controller) *RedisDS {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	controller.RegisterShutdown(
		"redis",
		func(ctx context.Context) error {
			return rdb.Shutdown(ctx).Err()
		})

	return &RedisDS{
		Client: rdb,
	}
}
