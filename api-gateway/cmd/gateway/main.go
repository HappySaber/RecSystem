package main

import (
	"api-gateway/internal/app"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		panic(fmt.Sprintf("failed to load .env: %v", err))
	}
	application, err := app.New(":8080")
	if err != nil {
		log.Fatal(err)
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	go func() {
		log.Println("API Gateway started on :8080")
		if err := application.Run(); err != nil {
			log.Println(err)
		}
	}()

	<-ctx.Done()

	log.Println("shutting down api gateway...")

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	_ = application.Shutdown(shutdownCtx)
}
