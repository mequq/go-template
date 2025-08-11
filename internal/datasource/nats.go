package datasource

import (
	"application/pkg/initializer/config"
	"context"
	"log/slog"
	"strings"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type Nats struct {
	logger    *slog.Logger
	Client    *nats.Conn
	JetStream jetstream.JetStream
	Stream    jetstream.Stream
	// Protocol  *cejsm.Protocol
}

type NatsConfig struct {
	DSN           string `json:"dsn"`
	InitJetstream bool   `json:"initJetstream"`
	StreamName    string `json:"streamName"`
	Subjects      string `json:"subject"`
}

func NewNats(ctx context.Context, logger *slog.Logger, config config.Config) (*Nats, error) {

	cfg := new(NatsConfig)
	if err := config.Unmarshal("datasource.nats", cfg); err != nil {
		logger.Error("Failed to unmarshal NATS config", "error", err)
		return nil, err
	}

	nc, err := nats.Connect(cfg.DSN)
	if err != nil {
		logger.Error("Failed to connect to NATS", "error", err)
		return nil, err
	}

	js, err := jetstream.New(nc)
	if err != nil {
		logger.Error("Failed to create JetStream context", "error", err)
		return nil, err
	}

	nats := &Nats{
		logger:    logger.With("module", "datasource.nats"),
		Client:    nc,
		JetStream: js,
	}

	if cfg.InitJetstream {
		if err := nats.initJetStream(ctx, cfg); err != nil {
			logger.Error("Failed to initialize JetStream", "error", err)
			return nil, err
		}
	} else {
		logger.Info("JetStream initialization skipped")
		stream, err := js.Stream(ctx, cfg.StreamName)
		if err != nil {
			logger.Error("Failed to get existing NATS stream", "stream_name", cfg.StreamName, "error", err)
			return nil, err
		}
		nats.Stream = stream
	}

	return nats, nil
}

func (n *Nats) initJetStream(ctx context.Context, cfg *NatsConfig) error {

	streamConfig := jetstream.StreamConfig{
		Name:        cfg.StreamName,
		Subjects:    strings.Split(cfg.Subjects, ","),
		Description: "Stream for FCM Campaigns",
		Storage:     jetstream.MemoryStorage,
		Replicas:    1,
		Retention:   jetstream.InterestPolicy,
	}

	s, err := n.JetStream.CreateOrUpdateStream(ctx, streamConfig)
	if err != nil {
		n.logger.Error("Failed to add NATS stream", "error", err)
		return err
	}

	n.Stream = s

	if _, err = s.CreateOrUpdateConsumer(ctx,
		jetstream.ConsumerConfig{
			Durable:       "campaigns",
			FilterSubject: "ir.lenz.fcm-campaign.campaigns",
			AckPolicy:     jetstream.AckExplicitPolicy,
		},
	); err != nil {
		return err
	}
	n.logger.Info("NATS stream initialized successfully", "stream_name", cfg.StreamName)
	return nil
}
