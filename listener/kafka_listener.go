package listener

import (
	"encoding/json"
	"errors"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/maskarb/sources-integration-go/xconfig"
)

type KafkaEvent struct {
	EventType  string `json:"event_type"`
	Offset     kafka.Offset
	Partition  int32
	ResourceID int `json:"resource_id"`
	SourceID   int `json:"source_id"`
	AuthHeader string
}

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

func GetMsgData(msg *kafka.Message) KafkaEvent {
	msgValue := KafkaEvent{}
	msgValue.Offset = msg.TopicPartition.Offset
	msgValue.Partition = msg.TopicPartition.Partition
	err := json.Unmarshal(msg.Value, &msgValue)
	if err != nil {
		panic(err)
	}
	msgValue.EventType, err = getHeader(msg.Headers, "event_type")
	if err != nil {
		panic(err)
	}
	msgValue.AuthHeader, err = getHeader(msg.Headers, "x-rh-identity")
	if err != nil {
		panic(err)
	}
	return msgValue
}

func getHeader(headers []kafka.Header, header string) (string, error) {
	for _, head := range headers {
		if head.Key == header {
			return string(head.Value), nil
		}
	}
	return "", errors.New(header + " not found in headers")
}
