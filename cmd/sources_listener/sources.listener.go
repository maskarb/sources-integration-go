package main

import (
	"flag"
	"fmt"

	"github.com/maskarb/sources-integration-go/listener"
	"github.com/maskarb/sources-integration-go/utils"
	"github.com/maskarb/sources-integration-go/xconfig"
	"github.com/sirupsen/logrus"
)

var EVENTTYPES = []string{
	"Application.create", "Application.destroy", "Authentication.create", "Authentication.update", "Source.update", "Source.destroy",
}

func main() {
	s := utils.NewSet()
	for _, v := range EVENTTYPES {
		s.Add(v)
	}
	flag.Parse()

	cfg, err := xconfig.LoadConfig("sources_listener")
	if err != nil {
		logrus.Fatal(err)
	}
	c := listener.NewConsumer(cfg.KafkaMain)
	fmt.Println(c)

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			fmt.Printf("consumed from topic %s [partition: %d] at offset %v: %v\n",
				*msg.TopicPartition.Topic, msg.TopicPartition.Partition, msg.TopicPartition.Offset, string(msg.Value))
			msgValue := listener.GetMsgData(msg)
			if !s.Contains(msgValue.EventType) {
				continue
			}

			fmt.Println(msgValue)
			fmt.Println()
			// if msgValue.EventType == "Application.create" {
			// 	source := getSource(*db, msgValue)
			// 	fmt.Println(source)
			// }
		} else {
			// The client will automatically try to recover from all errors.
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}

	c.Close()

}
