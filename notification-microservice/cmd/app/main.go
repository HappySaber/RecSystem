package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"notifications/internal/config"
	"notifications/internal/consumer"
	"notifications/internal/events"
	"notifications/internal/services"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		panic(fmt.Sprintf("failed to load .env: %v", err))
	}
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)
	log.Info(
		"starting application",
		slog.Any("env", cfg),
	)

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	mailer := services.NewMailSender(log)
	consumer := consumer.New(cfg.Kafka, log)

	consumer.Start(ctx, func(msg kafka.Message) error {
		log.Info(
			"received kafka message",
			slog.String("key", string(msg.Key)),
			slog.Int64("offset", msg.Offset),
		)

		to := string(msg.Key)

		var event events.UserRegisteredEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Error("failed to parse event", slog.Any("error", err))
			return err
		}

		body := fmt.Sprintf("Welcome, %s %s, to my recommendation app, hope you enjoy it!!", event.Name, event.Surname)

		sendCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		return mailer.SendWithRetry(
			sendCtx,
			to,
			"Добро пожаловать!",
			body,
		)
	})

	<-ctx.Done()
	log.Info("shutdown signal received")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}
