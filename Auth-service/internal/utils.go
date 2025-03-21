package internal

import (
	"errors"
	"time"

	"chat-grpc/proto_gen"
)

func ConvertTimeStamp(t *proto_gen.Timestamp) string {
	if t == nil {
		return ""
	}

	tConv := time.Unix(t.Seconds, int64(t.Nanos))
	return tConv.Format("2006-01-02 15:04:05")
}

func Validate(name, email, pass string) error {
	if name == "" {
		return errors.New("name field cannot be empty")
	}

	if email == "" {
		return errors.New("email field cannot be empty")
	}

	if len(pass) < 4 {
		return errors.New("pass must be more than 4 char")
	}

	return nil
}
