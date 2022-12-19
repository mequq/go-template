package data

import (
	"app/internal/biz"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
)

// userRepo struct
type UserRepo struct {
	data *Data
}

// NewUserRepo creates a new user repo.
func NewUserRepo(data *Data) biz.UserRepo {
	return &UserRepo{
		data: data,
	}
}

// Get gets the user by the user ID.
func (r *UserRepo) Get(c *gin.Context, id int64) (*biz.User, error) {
	_, span := otel.Tracer("healthz").Start(c.Request.Context(), "readiness-db-check")
	defer span.End()
	time.Sleep(20 * time.Millisecond)
	// TODO: implement the repository.
	return &biz.User{
		ID:   id,
		Name: "user",
	}, nil
}

// List gets the user list.
func (r *UserRepo) List(c *gin.Context) ([]*biz.User, error) {
	_, span := otel.Tracer("healthz").Start(c.Request.Context(), "readiness-db-check")
	defer span.End()
	time.Sleep(20 * time.Millisecond)
	return []*biz.User{
		{
			ID:   1,
			Name: "user1",
		},
		{
			ID:   2,
			Name: "user2",
		},
	}, nil
}
