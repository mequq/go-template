package data

import (
	"app/internal/biz"
	"context"

	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// healthz repo struct
type HealthzRepo struct {
	data *Data
}

// NewHealthzRepo creates a new healthz repo.
func NewHealthzRepo(data *Data) biz.HealthzRepo {
	return &HealthzRepo{
		data: data,
	}
}

// Readyness checks the readyness of the service.
func (r *HealthzRepo) Readiness(ctx context.Context) error {
	// ping the redis server
	_, err := r.data.redis.Ping().Result()
	if err != nil {
		return err
	}
	// ping the mongo server
	err = r.data.mongo.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		return err
	}
	return nil

}

// Liveness checks the liveness of the service.
func (r *HealthzRepo) Liveness(ctx context.Context) error {
	return nil
}
