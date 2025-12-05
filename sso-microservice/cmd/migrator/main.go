package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

	"github.com/pressly/goose/v3"
)

func main() {
	var migrationsPath string

	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")

	log.Println("Running migrations with Goose...")
	var db *sql.DB

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	if err := goose.Up(db, migrationsPath); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	fmt.Println("migrations applied successfully")

}
