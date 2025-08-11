package handler

import (
	"application/internal/datasource"
	"application/internal/service"
	"context"
	"log/slog"

	"github.com/nats-io/nats.go/jetstream"
)

type campainHandler struct {
	logger *slog.Logger
	nats   *datasource.Nats
}

var _ service.Handler = (*campainHandler)(nil)

func NewCampainHandler(logger *slog.Logger, nats *datasource.Nats) *campainHandler {
	return &campainHandler{
		logger: logger.With("module", "handler.campain"),
		nats:   nats,
	}
}

func (h *campainHandler) RegisterHandler(ctx context.Context) error {
	h.logger.Info("Registering campain handlers")
	// Here you would register your HTTP handlers, e.g.:

	c, err := h.nats.JetStream.Consumer(ctx, "mystream", "myconsumer")
	if err != nil {
		h.logger.Error("Failed to create consumer", "error", err)
		return err
	}

	c.Consume(h.handleCampaign)
	return nil
}

func (h *campainHandler) handleCampaign(msg jetstream.Msg) {
	h.logger.Info("Handling campaign message", "subject", msg.Subject, "data", string(msg.Data()))
	// Process the message here
}
