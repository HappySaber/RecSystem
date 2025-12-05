package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"sso-microservice/internal/domain/models"
	"sso-microservice/internal/storage"
	"strconv"

	"github.com/lib/pq"
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

func (s *Storage) SaveUser(ctx context.Context, email, name, surname, role string, passHash []byte) (string, error) {
	const op = "storage.posgresql.SaveUser"

	var id string

	query := "INSERT INTO users (email, name, surname pass_hash) VALUES ($1, $2, $3, $4) RETURNING id"

	err := s.db.QueryRowContext(ctx, query, email, name, surname, passHash).Scan(&id)
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return "", fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	const op = "storage.postgresql.User"

	query := "SELECT * FROM users WHERE email = $1"
	var user models.User
	err := s.db.QueryRowContext(ctx, query, email).Scan(&user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}

		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *Storage) App(ctx context.Context, id int) (models.App, error) {
	const op = "storage.posgresql.App"

	stmt, err := s.db.Prepare("SELECT id, name, secret FROM apps WHERE id = ?")
	if err != nil {
		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, id)

	var app models.App
	err = row.Scan(&app.ID, &app.Name, &app.Secret)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.App{}, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}

		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}

	return app, nil
}

func (s *Storage) UserRole(ctx context.Context, userID string) (string, error) {
	const op = "storage.posgresql.IsAdmin"

	stmt, err := s.db.Prepare("SELECT role FROM users WHERE id = ?")
	if err != nil {
		return "none", fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, userID)

	var role string

	err = row.Scan(&role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "none", fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}

		return "none", fmt.Errorf("%s: %w", op, err)
	}

	return role, nil
}
