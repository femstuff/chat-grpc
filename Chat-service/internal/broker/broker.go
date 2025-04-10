package broker

import (
	"chat-grpc/proto_gen"
)

type Broker interface {
	Publish(msg *proto_gen.Message) error
	Subscribe(chatID int64, handler func(*proto_gen.Message)) error
	Close() error
}
