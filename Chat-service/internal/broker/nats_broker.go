package broker

import (
	"encoding/json"
	"fmt"

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

func (b *natsBroker) Publish(msg *proto_gen.Message) error {
	subject := fmt.Sprintf("chat.%d", msg.ChatId)

	data, err := protojson.Marshal(msg)
	if err != nil {
		return err
	}

	return b.conn.Publish(subject, data)
}

func (b *natsBroker) Subscribe(chatID int64, handler func(msg *proto_gen.Message)) error {
	subject := fmt.Sprintf("chat.%d", chatID)
	_, err := b.conn.Subscribe(subject, func(msg *nats.Msg) {
		var m proto_gen.Message
		if err := json.Unmarshal(msg.Data, &m); err == nil {
			handler(&m)
		}
	})

	return err
}

func (b *natsBroker) Close() error {
	b.conn.Close()

	return nil
}
