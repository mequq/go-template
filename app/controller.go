package app

import (
	"context"
	"time"
)

type healthzOptions struct {
	timeout time.Duration
}

// HealthzOption is a function option for healthz registration.
type HealthzOption func(*healthzOptions)

// WithTimeout sets a timeout for the healthz check.
func WithTimeout(d time.Duration) HealthzOption {
	return func(o *healthzOptions) {
		o.timeout = d
	}
}

type Controller interface {
	GetSutdowners() map[string]func(ctx context.Context) error
	GetStarters() map[string]func(ctx context.Context) error
	RegisterShutdown(name string, shutdown func(ctx context.Context) error)
	RegisterStartup(name string, startup func(ctx context.Context) error)
	RegisterHealthz(name string, healthz func(ctx context.Context) error, opts ...HealthzOption)
	GetHealthz() map[string]func(ctx context.Context) error
}

var _ Controller = (*controller)(nil)

type controller struct {
	shutdowners map[string]func(ctx context.Context) error
	starters    map[string]func(ctx context.Context) error
	healthz     map[string]func(ctx context.Context) error
}

func NewController() *controller {
	return &controller{
		shutdowners: make(map[string]func(ctx context.Context) error),
		starters:    make(map[string]func(ctx context.Context) error),
		healthz:     make(map[string]func(ctx context.Context) error),
	}
}

// RegisterHealthz registers a healthz check with a name.
func (c *controller) RegisterHealthz(name string, healthz func(ctx context.Context) error, opts ...HealthzOption) {
	options := &healthzOptions{
		timeout: 5 * time.Second,
	}

	for _, o := range opts {
		o(options)
	}

	c.healthz[name] = func(ctx context.Context) error {
		ctx, cancel := context.WithTimeout(ctx, options.timeout)
		defer cancel()

		return healthz(ctx)
	}
}

func (c *controller) GetHealthz() map[string]func(ctx context.Context) error {
	return c.healthz
}

// RegisterShutdown registers a shutdown function with a name.
func (c *controller) RegisterShutdown(name string, shutdown func(ctx context.Context) error) {
	c.shutdowners[name] = shutdown
}

func (c *controller) RegisterStartup(name string, startup func(ctx context.Context) error) {
	c.starters[name] = startup
}

func (c *controller) GetSutdowners() map[string]func(ctx context.Context) error {
	return c.shutdowners
}

func (c *controller) GetStarters() map[string]func(ctx context.Context) error {
	return c.starters
}
