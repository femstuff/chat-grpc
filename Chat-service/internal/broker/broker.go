package broker

import (
	"chat-grpc/proto_gen"
	"github.com/nats-io/nats.go"
)

type Broker interface {
	Publish(msg *proto_gen.Message) error
	Subscribe(subject string, handler func(*proto_gen.Message)) (*nats.Subscription, error)
	Close() error
}
