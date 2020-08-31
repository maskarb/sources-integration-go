package listener

import (
	"encoding/json"
	"errors"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaEvent struct {
	EventType  string `json:"event_type"`
	Offset     kafka.Offset
	Partition  int32
	ResourceID int `json:"resource_id"`
	SourceID   int `json:"source_id"`
	AuthHeader string
}

func getHeader(headers []kafka.Header, header string) (string, error) {
	for _, head := range headers {
		if head.Key == header {
			return string(head.Value), nil
		}
	}
	return "", errors.New(header + " not found in headers")
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
