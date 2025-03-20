package internal

import (
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
