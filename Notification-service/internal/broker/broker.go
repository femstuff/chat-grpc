package broker

import (
	"encoding/json"

	"chat-grpc/Notification-service/internal/entity"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type NatsConsumer struct {
	nc      *nats.Conn
	subject string
	log     *zap.Logger
	handler func(message entity.Message)
}

func NewNatsConsumer(nc *nats.Conn, subject string, log *zap.Logger, handler func(message entity.Message)) *NatsConsumer {
	return &NatsConsumer{
		nc:      nc,
		subject: subject,
		log:     log,
		handler: handler,
	}
}

func (c *NatsConsumer) Start() error {
	_, err := c.nc.Subscribe(c.subject, func(msg *nats.Msg) {
		var m entity.Message
		if err := json.Unmarshal(msg.Data, &m); err != nil {
			c.log.Error("Failed to parse message", zap.Error(err))
		}
		c.handler(m)
	})

	if err != nil {
		c.log.Error("Failed to subscribe", zap.Error(err))
		return err
	}
	c.log.Info("Subscribe NATS", zap.String("subject", c.subject))
	return nil
}
