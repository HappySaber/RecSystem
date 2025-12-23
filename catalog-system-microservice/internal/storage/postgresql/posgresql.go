package postgresql

import (
	"catalog-microservice/internal/domain/models"
	"catalog-microservice/internal/storage"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/lib/pq"
)

type Storage struct {
	DB *sql.DB
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
		DB: db,
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
		Host:     os.Getenv("DB_HOST"),
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

func (s *Storage) Create(content *models.Content) (string, error) {
	const op = "storage.posgresql.Create"

	var id string
	query := `
	INSERT INTO content (
		type,
		external_source,
		external_id,
		title,
		description,
		poster_url,
		release_date
	)
	VALUES ($1,$2,$3,$4,$5,$6,$7)
	RETURNING id
	`

	err := s.DB.QueryRow(query, content.Type, content.ExternalSource, content.ExternalID, content.Title, content.Description, content.PosterURL, content.ReleaseDate).Scan(&id)
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return "", fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) CreateOrUpdateContent(content *models.Content) (string, error) {
	const op = "storage.posgresql.CreateOrUpdate"
	var releaseDate interface{}
	if content.ReleaseDate == "" {
		releaseDate = nil
	} else {
		releaseDate = content.ReleaseDate
	}
	var id string
	query := `
		INSERT INTO content (type, external_source, external_id, title, description, poster_url, release_date)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
		ON CONFLICT (external_source, external_id)
		DO UPDATE SET
			type = EXCLUDED.type,
			title = EXCLUDED.title,
			description = EXCLUDED.description,
			poster_url = EXCLUDED.poster_url,
			release_date = EXCLUDED.release_date,
			updated_at = now()
		RETURNING id
	`

	err := s.DB.QueryRow(query,
		content.Type,
		content.ExternalSource,
		content.ExternalID,
		content.Title,
		content.Description,
		content.PosterURL,
		releaseDate,
	).Scan(&id)

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) SaveMovie(movie *models.MovieDetails) error {
	_, err := s.DB.Exec(`
INSERT INTO movies_details (
	content_id, tmdb_id, original_title, runtime, tagline,
	status, budget, revenue, language,
	genres, cast_members, crew, images, videos, raw_data
)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)
ON CONFLICT (content_id) DO UPDATE SET
	tmdb_id = EXCLUDED.tmdb_id,
	raw_data = EXCLUDED.raw_data
`,
		movie.ContentID,
		movie.TmdbID,
		movie.OriginalTitle,
		movie.Runtime,
		movie.Tagline,
		movie.Status,
		movie.Budget,
		movie.Revenue,
		movie.Language,
		movie.Genres,
		movie.CastMembers,
		movie.Crew,
		movie.Images,
		movie.Videos,
		movie.RawData,
	)

	return err
}

func (s *Storage) SaveSeries(contentID string, tmdbID int, rawJSON string) error {
	_, err := s.DB.Exec(`
		INSERT INTO series_details (content_id, tmdb_id, raw_data)
		VALUES ($1, $2, $3)
		ON CONFLICT (content_id)
		DO UPDATE SET
			tmdb_id = EXCLUDED.tmdb_id,
			raw_data = EXCLUDED.raw_data
	`, contentID, tmdbID, rawJSON)
	return err
}
