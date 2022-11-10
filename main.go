package main

import (
	"context"
	"flag"
	"log"

	tgClient "telegram-bot/clients/telegram"
	"telegram-bot/consumer/event-consumer"
	telegram "telegram-bot/events/telegram"
	// "telegram-bot/storage/files"
	_ "github.com/mattn/go-sqlite3"
	"telegram-bot/storage/sqlite"
)

const (
	tgBotHost         = "api.telegram.org"
	sqliteStoragePath = "data/sqlite/storage.db"
	// storagePath       = "storage"
	batchSize = 100
)

func main() {
	// s := files.New(storagePath)
	s, err := sqlite.New(sqliteStoragePath)

	if err != nil {
		log.Fatalf("can't connect to storage: ", err)
	}

	err = s.Init(context.TODO())
	if err != nil {
		log.Fatal("can't init starage:", err)
	}

	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, mustToken()),
		s,
	)
	log.Print("service started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}

func mustToken() string {
	token := flag.String("tg-bot-token", "", "token for access to telegram bot")
	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified")
	}

	return *token
}
