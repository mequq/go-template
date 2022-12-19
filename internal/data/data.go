package data

import (
	"app/config"
	"context"
	"time"

	"github.com/go-redis/redis"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// provide set
var DataProviderSet = wire.NewSet(NewData, NewHealthzRepo, NewUserRepo)

// Data provider struct
type Data struct {
	conf  *config.Config
	redis *redis.Client
	mongo *mongo.Client
}

// NewData creates a new data.
func NewData(c *config.Config) (*Data, func(), error) {

	data := &Data{
		conf: c,
	}
	data.redis = initRedis(c)

	cleanup := func() {
		data.redis.Close()
		data.mongo.Disconnect(context.Background())
	}
	client, err := initMongo(c)
	if err != nil {
		return nil, nil, err
	}
	data.mongo = client
	return data, cleanup, nil

}

// init redis
func initRedis(c *config.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     c.Redis.Host + ":" + c.Redis.Port,
		Password: c.Redis.Password,
		DB:       c.Redis.DB,
	})
}

// init mongo
func initMongo(c *config.Config) (*mongo.Client, error) {
	uri := "mongodb://localhost:27017"
	mongoOptions := options.Client().ApplyURI(uri)
	mongoOptions.SetTimeout(1 * time.Second)
	client, err := mongo.Connect(context.TODO(), mongoOptions)
	if err != nil {
		return nil, err
	}
	return client, nil
}
