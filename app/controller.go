package app

import (
	"context"
	"errors"
	"fmt"
)

type Controller interface {
	GetSutdowners() map[string]func(ctx context.Context) error
	GetStarters() map[string]func(ctx context.Context) error
	RegisterShutdown(name string, shutdown func(ctx context.Context) error)
	RegisterStartup(name string, startup func(ctx context.Context) error)
}

var _ Controller = (*controller)(nil)

type controller struct {
	shutdowners map[string]func(ctx context.Context) error
	starters    map[string]func(ctx context.Context) error
}

func NewController() *controller {
	return &controller{
		shutdowners: make(map[string]func(ctx context.Context) error),
		starters:    make(map[string]func(ctx context.Context) error),
	}
}

// Shutdown implements Controller.
func (c *controller) Shutdown(ctx context.Context) error {
	for name, shutdown := range c.shutdowners {
		if err := shutdown(ctx); err != nil {
			return errors.Join(fmt.Errorf("failed to shutdown component %s: %w", name, err))
		}
	}

	return nil
}

func (c *controller) Start(ctx context.Context) error {
	for name, startup := range c.starters {
		if err := startup(ctx); err != nil {
			return errors.Join(fmt.Errorf("failed to start component %s: %w", name, err))
		}
	}

	return nil
}

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
