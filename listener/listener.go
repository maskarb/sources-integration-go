package listener

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/maskarb/sources-integration-go/xconfig"
)

func NewConsumer(cfg *xconfig.Kafka) *kafka.Consumer {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": cfg.Addr,
		"group.id":          cfg.GroupID,
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		panic(err)
	}

	c.SubscribeTopics([]string{cfg.Topic, "^aRegex.*[Tt]opic"}, nil)

	return c
}
