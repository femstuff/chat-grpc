package broker

import (
	"context"

	"chat-grpc/proto_gen"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"
)

type NatsConsumer struct {
	nc      *nats.Conn
	subject string
	log     *zap.Logger
	handler func(ctx context.Context, m *proto_gen.Message)
}

func NewNatsConsumer(nc *nats.Conn, subject string, log *zap.Logger, handler func(ctx context.Context, m *proto_gen.Message)) *NatsConsumer {
	return &NatsConsumer{
		nc:      nc,
		subject: subject,
		log:     log,
		handler: handler,
	}
}

func (c *NatsConsumer) Start() error {
	_, err := c.nc.Subscribe(c.subject, func(msg *nats.Msg) {
		var m proto_gen.Message
		if err := protojson.Unmarshal(msg.Data, &m); err != nil {
			c.log.Error("Failed to parse proto message", zap.Error(err))
			return
		}
		c.handler(context.Background(), &m)
	})

	if err != nil {
		c.log.Error("Failed to subscribe", zap.Error(err))
		return err
	}
	c.log.Info("Subscribed to NATS", zap.String("subject", c.subject))
	return nil
}
