package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"github.com/maskarb/sources-integration/utils"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

var TOPIC = getEnv("SOURCES_TOPIC", "platform.sources.event-stream")
var EVENTTYPES = []string{
	"Application.create", "Application.destroy", "Authentication.create", "Authentication.update", "Source.update", "Source.destroy",
}

type KafkaEvent struct {
	EventType  string `json:"event_type"`
	Offset     kafka.Offset
	Partition  int32
	ResourceID int `json:"resource_id"`
	SourceID   int `json:"source_id"`
	AuthHeader string
}

type Source struct {
	tableName        struct{}  `pg:"api_sources"`
	ID               int       `pg:"source_id"` // Id is automatically detected as primary key
	SourceUUID       uuid.UUID `pg:"source_uuid"`
	Name             string
	AuthHeader       string
	Offset           int
	AccountID        int
	SourceType       string
	Authentication   struct{}
	BillingSource    struct{}
	KokuUUID         string `pg:"koku_uuid"`
	PendingDelete    bool
	PendingUpdate    bool
	OutOfOrderDelete bool
	Status           struct{}
}

func getSource(db pg.DB, msgValue KafkaEvent) Source {
	source := new(Source)
	err := db.Model(source).Where("source_id = ?", msgValue.SourceID).Select()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	return *source
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func getHeader(headers []kafka.Header, header string) (string, error) {
	for _, head := range headers {
		if head.Key == header {
			return string(head.Value), nil
		}
	}
	return "", errors.New(header + " not found in headers")
}

func getMsgData(msg *kafka.Message) KafkaEvent {
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

func main() {
	s := utils.NewSet()
	for _, v := range EVENTTYPES {
		s.Add(v)
	}

	db := pg.Connect(&pg.Options{
		User:     "postgres",
		Password: "postgres",
		Database: "postgres",
		Addr:     "localhost:15432",
	})
	defer db.Close()
	if err := db.Ping(db.Context()); err != nil {
		panic(err)
	}

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:29091",
		"group.id":          "hccm-sources",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		panic(err)
	}

	c.SubscribeTopics([]string{TOPIC, "^aRegex.*[Tt]opic"}, nil)

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			msgValue := getMsgData(msg)
			if !s.Contains(msgValue.EventType) {
				continue
			}
			fmt.Printf("consumed from topic %s [partition: %d] at offset %v: %v\n",
				*msg.TopicPartition.Topic, msg.TopicPartition.Partition, msg.TopicPartition.Offset, string(msg.Value))
			fmt.Println(msgValue)
			fmt.Println()
			if msgValue.EventType == "Application.create" {
				source := getSource(*db, msgValue)
				fmt.Println(source)
			}
		} else {
			// The client will automatically try to recover from all errors.
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}

	c.Close()
}
