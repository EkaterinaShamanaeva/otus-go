package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	configuration "github.com/EkaterinaShamanaeva/otus-go/hw12_13_14_15_calendar/internal/config"
	"github.com/EkaterinaShamanaeva/otus-go/hw12_13_14_15_calendar/internal/logger"
	"github.com/EkaterinaShamanaeva/otus-go/hw12_13_14_15_calendar/internal/rabbitmq"
	"github.com/EkaterinaShamanaeva/otus-go/hw12_13_14_15_calendar/internal/storage"
	"log"
	"os"
	"os/signal"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/config_rmq.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	rmqConfig := configuration.NewConfigRMQ()

	if err := rmqConfig.BuildConfigRMQ(configFile); err != nil {
		log.Fatalf("Config error: %v", err)
	}

	logg, err := logger.New(rmqConfig.Logger.Level, rmqConfig.Logger.Path)
	if err != nil {
		log.Fatalf("Logger error: %v", err)
	}
	logg.Info("sender start")
	defer logg.Info("\nsender end")

	consumer := rabbitmq.NewConsumer(rmqConfig.RMQ.Uri, rmqConfig.RMQ.Queue, logg)

	err = consumer.Connect()
	if err != nil {
		logg.Error("failed to connect RabbitMQ")
	}

	message, err := consumer.Consume()
	if err != nil {
		logg.Error("failed to consume")
	}

	sigChan := make(chan os.Signal, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-message:
				var notice storage.Notification
				if err = json.Unmarshal(msg.Body, &notice); err != nil {
					logg.Error(fmt.Sprintf("can't unmarshal notice: %v", err))
					continue
				}
				logg.Info(fmt.Sprintf(
					"ID: %d, title: \"%s\", datetime: %s, owner id: %d",
					notice.ID, notice.Title, notice.Datetime, notice.UserID,
				))
			}
		}
	}(ctx)

	signal.Notify(sigChan, os.Interrupt)
	<-sigChan
}
