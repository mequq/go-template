package biz

import (
	"app/config"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// user is the user entity.
type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// userRepo is the user interface.
type UserRepo interface {
	// Get gets the user by the user ID.
	Get(c *gin.Context, id int64) (*User, error)
	List(c *gin.Context) ([]*User, error)
}

// userUsecase is the user usecase.
type UserUsecase struct {
	repo   UserRepo
	config *config.Config
}

// NewUserUsecase creates a new user usecase.
func NewUserUsecase(repo UserRepo, c *config.Config) *UserUsecase {
	return &UserUsecase{
		repo:   repo,
		config: c,
	}
}

// Get gets the user by the user ID.
func (u *UserUsecase) Get(c *gin.Context, id int64) (*User, error) {
	ctx, span := otel.Tracer("healthz").Start(c.Request.Context(), "readiness-usecase")
	span.SetAttributes(attribute.Int("id", int(id)))
	defer span.End()

	url := u.config.Tracing.TUrl

	// ----------------
	// use otelhttp for calling external service
	resp, err := otelhttp.Get(ctx, url)
	if err != nil {
		return nil, err
	}
	// *** IMPORTANT ***
	// need to close the body for the connection to be reused and not leak and close the span
	// *** IMPORTANT ***
	resp.Body.Close()
	// ----------------

	fmt.Println("response: ", url)

	time.Sleep(10 * time.Millisecond)
	c.Request = c.Request.WithContext(ctx)
	span.AddEvent("get user")
	return u.repo.Get(c, id)
}

func (u *UserUsecase) List(c *gin.Context) ([]*User, error) {
	ctx, span := otel.Tracer("healthz").Start(c.Request.Context(), "readiness-usecase")
	defer span.End()
	url := u.config.Tracing.TUrl

	// ----------------
	// second way to use otelhttp for calling external service
	client := &http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	// *** IMPORTANT ***
	// need to close the body for the connection to be reused and not leak and close the span
	// *** IMPORTANT ***
	resp.Body.Close()
	// ----------------

	fmt.Println("response: ", url)
	time.Sleep(10 * time.Millisecond)
	c.Request = c.Request.WithContext(ctx)
	return u.repo.List(c)
}
