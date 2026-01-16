package postgresql

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
)

type Storage struct {
	db *sql.DB
}

func New() (*Storage, error) {
	const op = "storage.postgresql.New"
	dbConfig := buildDBConfig()
	db, err := sql.Open("postgres", dbConfig.dsn())
	if err != nil {
		log.Fatalf("Error checking database connection: %v", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Println("Successfully connected to the database!")

	return &Storage{
		db: db,
	}, nil
}

type DBConfig struct {
	Host     string
	Port     int
	User     string
	DBName   string
	Password string
}

func buildDBConfig() *DBConfig {
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Fatalf("Invalid DB_PORT: %v", err)
	}
	return &DBConfig{
		Host:     os.Getenv("DB_HOST_LOCAL"),
		Port:     port,
		User:     os.Getenv("DB_USER"),
		DBName:   os.Getenv("DB_NAME"),
		Password: os.Getenv("DB_PASSWORD"),
	}
}

func (config *DBConfig) dsn() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.DBName,
	)
}
