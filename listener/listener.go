package listener

import (
	"github.com/segmentio/kafka-go"
)

func NewConsumer() *kafka.Reader {
	// make a new reader that consumes from topic-A
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"localhost:29091"},
		GroupID:  "hccm-sources",
		Topic:    "platform.sources.event-stream",
		MinBytes: 10e0, // 10KB
		MaxBytes: 10e6, // 10MB
	})
}
