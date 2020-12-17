package main

import (
	"context"

	"github.com/maskarb/sources-integration-go/database"
	"github.com/maskarb/sources-integration-go/listener"
	"github.com/maskarb/sources-integration-go/sources"
	"github.com/maskarb/sources-integration-go/strset"
	log "github.com/sirupsen/logrus"
)

var EVENTTYPES = []string{
	"Application.create", "Application.destroy", "Authentication.create", "Authentication.update", "Source.update", "Source.destroy",
}

func main() {
	log.SetFormatter(&log.TextFormatter{})
	ctx := context.Background()

	log.Info("getting DB connection")
	db := database.NewConnection()

	log.Info("listener started")
	r := listener.NewConsumer()
	s := strset.New(EVENTTYPES...)
	for {
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			break
		}
		m := listener.GetMsgData(&msg)
		if !s.Contains(m.EventType) {
			log.Info("skipping msg")
			continue
		}
		log.Infof("message: %+v\nAUTH: %+v\n", m, *m.AuthHeader)
		if m.SourceID != 0 {
			source, err := sources.SelectSource(ctx, db, m.SourceID)
			if err != nil {
				panic(err)
			}
			log.Infof("SOURCE: %+v", source)
		}

	}

	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
}
