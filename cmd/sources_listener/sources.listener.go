package main

import (
	"flag"
	"fmt"

	"github.com/maskarb/sources-integration-go/listener"
	"github.com/maskarb/sources-integration-go/strset"
	"github.com/maskarb/sources-integration-go/xconfig"
	log "github.com/sirupsen/logrus"
)

var EVENTTYPES = []string{
	"Application.create", "Application.destroy", "Authentication.create", "Authentication.update", "Source.update", "Source.destroy",
}

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	s := strset.New(EVENTTYPES...)

	flag.Parse()

	cfg, err := xconfig.LoadConfig("sources_listener")
	if err != nil {
		log.Fatal(err)
	}
	c := listener.NewConsumer(cfg.KafkaMain)

	for {
		msg, err := c.ReadMessage(-1)
		if err != nil {
			log.Error("Consumer error: %v (%v)\n", err, msg)
		}
		msgValue := listener.GetMsgData(msg)
		if !s.Contains(msgValue.EventType) {
			continue
		}
		log.Info(fmt.Sprintf("%+v", msgValue)) // format the KafkaEvent with Sprintf %+v
		// if msgValue.EventType == "Application.create" {
		// 	source := getSource(*db, msgValue)
		// 	fmt.Println(source)
		// }
	}

	c.Close()

}
