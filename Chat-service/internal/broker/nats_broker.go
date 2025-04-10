package broker

import (
	"fmt"
	"log"

	"chat-grpc/proto_gen"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/encoding/protojson"
)

type natsBroker struct {
	conn *nats.Conn
}

func NewNatsBroker(url string) (Broker, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	return &natsBroker{conn: nc}, nil
}

func (b *natsBroker) Subscribe(subject string, handler func(*proto_gen.Message)) (*nats.Subscription, error) {
	return b.conn.Subscribe(subject, func(m *nats.Msg) {
		var msg proto_gen.Message
		if err := protojson.Unmarshal(m.Data, &msg); err != nil {
			log.Fatal("Failed to parse json")
			return
		}
		handler(&msg)
	})
}

func (b *natsBroker) Publish(msg *proto_gen.Message) error {
	subject := fmt.Sprintf("chat.%d", msg.ChatId)

	data, err := protojson.Marshal(msg)
	if err != nil {
		return err
	}

	return b.conn.Publish(subject, data)
}

func (b *natsBroker) Close() error {
	b.conn.Close()

	return nil
}
