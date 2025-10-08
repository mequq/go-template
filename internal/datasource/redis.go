package datasource

// import (
// 	"application/app"
// 	"context"

// 	"github.com/redis/go-redis/extra/redisotel/v9"
// 	"github.com/redis/go-redis/v9"
// )

// type RedisConfig struct {
// 	Addr     string `koanf:"addr"`
// 	Password string `koanf:"password"`
// 	DB       int    `koanf:"db"`
// }

// func NewRedisConfig(k *app.KConfig) (*RedisConfig, error) {
// 	return &RedisConfig{
// 		Addr:     "localhost:6379",
// 		Password: "",
// 		DB:       0,
// 	}, nil
// }

// type RedisDS struct {
// 	*redis.Client
// }

// func NewRedisDS(ctx context.Context, cfg *RedisConfig, controller app.Controller) *RedisDS {
// 	rdb := redis.NewClient(&redis.Options{
// 		Addr:     cfg.Addr,
// 		Password: cfg.Password,
// 		DB:       cfg.DB,
// 	})

// 	controller.RegisterShutdown(
// 		"redis",
// 		func(ctx context.Context) error {
// 			return rdb.Shutdown(ctx).Err()
// 		})

// 	controller.RegisterHealthz("redis", func(ctx context.Context) error {
// 		return rdb.Ping(ctx).Err()
// 	})

// 	if err := redisotel.InstrumentTracing(rdb); err != nil {
// 		panic(err)
// 	}

// 	if err := redisotel.InstrumentMetrics(rdb); err != nil {
// 		panic(err)
// 	}

// 	return &RedisDS{
// 		Client: rdb,
// 	}
// }
