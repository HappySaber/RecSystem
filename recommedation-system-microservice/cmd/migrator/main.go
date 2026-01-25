package main

import (
	"flag"
	"fmt"
	"log"
	"rec-system-microservice/internal/storage/postgresql"

	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
)

func main() {
	var migrationsPath string

	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")
	flag.Parse()

	if err := godotenv.Load(".env"); err != nil {
		panic(fmt.Sprintf("failed to load .env: %v", err))
	}
	db, err := postgresql.New()
	if err != nil {
		log.Fatalf("Could not open database %s", err)
	}

	log.Println("Running migrations with Goose...")

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	log.Println(migrationsPath)
	if err := goose.Up(db.DB, migrationsPath); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	fmt.Println("migrations applied successfully")

}
