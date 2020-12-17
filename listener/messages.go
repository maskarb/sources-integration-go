package listener

import (
	"encoding/base64"
	"encoding/json"
	"errors"

	"github.com/segmentio/kafka-go"
)

type KafkaEvent struct {
	Offset    int64
	Partition int
	Topic     string

	EventType   string      `json:"event_type"`
	ResourceID  int         `json:"resource_id"`
	SourceID    int         `json:"source_id"`
	AuthType    string      `json:"authtype"`
	AuthHeader  *AuthHeader `json:"auth_header"`
	EncodedAuth string
}

func getHeader(headers []kafka.Header, header string) (string, error) {
	for _, head := range headers {
		if head.Key == header {
			return string(head.Value), nil
		}
	}
	return "", errors.New(header + " not found in headers")
}

func GetMsgData(m *kafka.Message) KafkaEvent {
	msg := KafkaEvent{
		Offset:    m.Offset,
		Partition: m.Partition,
		Topic:     m.Topic,
	}

	err := json.Unmarshal(m.Value, &msg)
	if err != nil {
		panic(err)
	}
	msg.EventType, err = getHeader(m.Headers, "event_type")
	if err != nil {
		panic(err)
	}
	msg.EncodedAuth, err = getHeader(m.Headers, "x-rh-identity")
	if err != nil {
		panic(err)
	}
	data, err := base64.StdEncoding.DecodeString(msg.EncodedAuth)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &msg.AuthHeader)
	if err != nil {
		panic(err)
	}

	return msg
}
