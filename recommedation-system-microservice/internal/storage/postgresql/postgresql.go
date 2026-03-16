package postgresql

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

type Storage struct {
	DB *sql.DB
}

func New() (*Storage, error) {
	return NewWithDSN(buildDSNFromEnv())
}

type DBConfig struct {
	Host     string
	Port     int
	User     string
	DBName   string
	Password string
}

func NewWithDSN(dsn string) (*Storage, error) {
	const op = "storage.postgresql.NewWithDSN"

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// sql.Open не проверяет соединение, Ping проверяет
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: ping failed: %w", op, err)
	}

	log.Println("successfully connected to the database")
	return &Storage{DB: db}, nil
}

func buildDSNFromEnv() string {
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Fatalf("invalid DB_PORT: %v", err)
	}
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		port,
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
}
