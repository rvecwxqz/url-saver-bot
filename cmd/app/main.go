package main

import (
	"context"
	"log"
	"os"
	"os/exec"
	tgClient "url-saver-bot/internal/clients/telegram"
	"url-saver-bot/internal/config"
	eventConsumer "url-saver-bot/internal/consumer/event-consumer"
	"url-saver-bot/internal/events/telegram"
	"url-saver-bot/internal/storage/db"
)

const batchSize = 100

func main() {
	go runPython()

	cfg := config.NewConfig()
	ctx := context.Background()

	eventProcessor := telegram.New(
		ctx,
		tgClient.NewClient(cfg.Token),
		db.NewDBStorage(ctx, cfg.DatabaseDSN),
		cfg.TagBufferSize,
	)
	log.Println("service started")

	consumer := eventConsumer.New(eventProcessor, eventProcessor, batchSize)
	if err := consumer.Start(); err != nil {
		log.Fatal(err)
	}
}

func runPython() {
	cmd := exec.Command("python", "./internal/ml/bert-classifier/main.py")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
