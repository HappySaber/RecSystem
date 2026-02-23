package postgresql

import (
	"catalog-microservice/internal/domain/models"
	"catalog-microservice/internal/storage"
	"context"
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

func (s *Storage) GetContent(ctx context.Context, id string) (models.Content, error) {
	const op = "storage.posgresql.GetContent"

	query := `SELECT id, type, title, external_id, description, release_date FROM content WHERE id = $1`

	var details models.Content
	err := s.DB.QueryRowContext(ctx, query, id).Scan(
		&details.ID,
		&details.Type,
		&details.Title,
		&details.ExternalID,
		&details.Description,
		&details.ReleaseDate,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Content{}, fmt.Errorf("%s: %w", op, storage.ErrShowNotFound)
		}

		return models.Content{}, fmt.Errorf("%s: %w", op, err)
	}

	return details, nil
}

func (s *Storage) GetContentByIDs(ctx context.Context, ids []string) ([]models.ContentShort, error) {
	const op = "storage.posgresql.GetContentByIDs"

	if len(ids) == 0 {
		return []models.ContentShort{}, nil
	}

	query := `
		SELECT id, type, title
		FROM anime_details
		WHERE id = ANY($1)
	`

	rows, err := s.DB.QueryContext(ctx, query, pq.Array(ids))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	contents := make([]models.ContentShort, 0, len(ids))

	for rows.Next() {
		var content models.ContentShort
		if err := rows.Scan(
			&content.ID,
			&content.Type,
			&content.Title,
		); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		contents = append(contents, content)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return contents, nil
}

func (s *Storage) AnimeDetails(ctx context.Context, id string) (models.AnimeDetails, error) {
	const op = "storage.posgresql.AnimeDetails"

	query := `SELECT * FROM anime_details WHERE id = $1`

	var animeDetails models.AnimeDetails
	err := s.DB.QueryRowContext(ctx, query, id).Scan(&animeDetails)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.AnimeDetails{}, fmt.Errorf("%s: %w", op, storage.ErrShowNotFound)
		}

		return models.AnimeDetails{}, fmt.Errorf("%s: %w", op, err)
	}

	return animeDetails, nil
}

func (s *Storage) MovieDetails(ctx context.Context, id string) (models.MovieDetails, error) {
	const op = "storage.posgresql.MovieDetails"

	query := `SELECT * FROM movies_details WHERE id = $1`
	var movieDetails models.MovieDetails
	err := s.DB.QueryRowContext(ctx, query, id).Scan(&movieDetails)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.MovieDetails{}, fmt.Errorf("%s: %w", op, storage.ErrShowNotFound)
		}

		return models.MovieDetails{}, fmt.Errorf("%s: %w", op, err)
	}

	return movieDetails, nil
}

func (s *Storage) SeriesDetails(ctx context.Context, id string) (models.SeriesDetails, error) {
	const op = "storage.posgresql.SeriesDetails"

	query := `SELECT * FROM series_details WHERE id = $1`
	var seriesDetails models.SeriesDetails
	err := s.DB.QueryRowContext(ctx, query, id).Scan(&seriesDetails)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.SeriesDetails{}, fmt.Errorf("%s: %w", op, storage.ErrShowNotFound)
		}

		return models.SeriesDetails{}, fmt.Errorf("%s: %w", op, err)
	}

	return seriesDetails, nil
}

func (s *Storage) BookDetails(ctx context.Context, id string) (models.BookDetails, error) {
	const op = "storage.posgresql.BookDetails"
	query := `SELECT * FROM books_details WHERE id = $1`
	var bookDetails models.BookDetails
	err := s.DB.QueryRowContext(ctx, query, id).Scan(&bookDetails)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.BookDetails{}, fmt.Errorf("%s: %w", op, storage.ErrShowNotFound)
		}
		return models.BookDetails{}, fmt.Errorf("%s: %w", op, err)
	}
	return bookDetails, nil
}

func (s *Storage) GameDetails(ctx context.Context, id string) (models.GameDetails, error) {
	const op = "storage.posgresql.GameDetails"
	query := `SELECT * FROM games_details WHERE id = $1`
	var gameDetails models.GameDetails
	err := s.DB.QueryRowContext(ctx, query, id).Scan(&gameDetails)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.GameDetails{}, fmt.Errorf("%s: %w", op, storage.ErrShowNotFound)
		}
		return models.GameDetails{}, fmt.Errorf("%s: %w", op, err)
	}
	return gameDetails, nil
}

func (s *Storage) AllAnimeDetails(ctx context.Context) ([]models.AnimeDetails, error) {
	const op = "storage.posgresql.AllAnimeDetails"

	query := `SELECT * FROM anime_details`
	rows, err := s.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var animeDetails []models.AnimeDetails
	for rows.Next() {
		var anime models.AnimeDetails
		if err := rows.Scan(&anime); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		animeDetails = append(animeDetails, anime)
	}
	return animeDetails, nil
}

func (s *Storage) AllMovieDetails(ctx context.Context) ([]models.MovieDetails, error) {
	const op = "storage.posgresql.AllMovieDetails"

	query := `SELECT * FROM movies_details`
	rows, err := s.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var movieDetails []models.MovieDetails
	for rows.Next() {
		var movie models.MovieDetails
		if err := rows.Scan(&movie); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		movieDetails = append(movieDetails, movie)
	}
	return movieDetails, nil
}

func (s *Storage) AllSeriesDetails(ctx context.Context) ([]models.SeriesDetails, error) {
	const op = "storage.posgresql.AllSeriesDetails"

	query := `SELECT * FROM series_details`
	rows, err := s.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var seriesDetails []models.SeriesDetails
	for rows.Next() {
		var series models.SeriesDetails
		if err := rows.Scan(&series); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		seriesDetails = append(seriesDetails, series)
	}
	return seriesDetails, nil
}

func (s *Storage) AllBookDetails(ctx context.Context) ([]models.BookDetails, error) {
	const op = "storage.posgresql.AllBookDetails"

	query := `SELECT * FROM books_details`
	rows, err := s.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var bookDetails []models.BookDetails
	for rows.Next() {
		var book models.BookDetails
		if err := rows.Scan(&book); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		bookDetails = append(bookDetails, book)
	}
	return bookDetails, nil
}

func (s *Storage) AllGameDetails(ctx context.Context) ([]models.GameDetails, error) {
	const op = "storage.posgresql.AllGameDetails"

	query := `SELECT * FROM games_details`
	rows, err := s.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var gameDetails []models.GameDetails
	for rows.Next() {
		var game models.GameDetails
		if err := rows.Scan(&game); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		gameDetails = append(gameDetails, game)
	}
	return gameDetails, nil
}

// func (s *storage) CreateContent() {
// 	const op = "storage.posgresql.CreateContent"
// 	return fmt.Errorf("%s: not implemented", op)
// }

func (s *Storage) FindContentByExternal(ctx context.Context, externalID, externalSource string) (models.Content, error) {
	const op = "storage.posgresql.FindContentByExternal"

	query := `SELECT id, type, title, external_id, description, release_date FROM content WHERE external_id = $1 AND external_source = $2`

	var details models.Content
	err := s.DB.QueryRowContext(ctx, query, externalID, externalSource).Scan(
		&details.ID,
		&details.Type,
		&details.Title,
		&details.ExternalID,
		&details.Description,
		&details.ReleaseDate,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Content{}, fmt.Errorf("%s: %w", op, storage.ErrShowNotFound)
		}

		return models.Content{}, fmt.Errorf("%s: %w", op, err)
	}

	return details, nil
}
