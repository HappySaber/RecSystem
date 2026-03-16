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

func (s *Storage) GetContent(ctx context.Context, id string) (models.Content, error) {
	const op = "storage.posgresql.GetContent"

	query := `
        SELECT id, type, external_source, external_id,
               title,
               COALESCE(description, ''),
               COALESCE(poster_url, ''),
               COALESCE(release_date::text, ''),
               created_at::text,
               updated_at::text
        FROM content WHERE id = $1`

	var c models.Content
	err := s.DB.QueryRowContext(ctx, query, id).Scan(
		&c.ID,
		&c.Type,
		&c.ExternalSource,
		&c.ExternalID,
		&c.Title,
		&c.Description,
		&c.PosterURL,
		&c.ReleaseDate,
		&c.CreatedAt,
		&c.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Content{}, fmt.Errorf("%s: %w", op, storage.ErrNotFound)
		}

		return models.Content{}, fmt.Errorf("%s: %w", op, err)
	}

	return c, nil
}

func (s *Storage) GetContentByIDs(ctx context.Context, ids []string) ([]models.ContentShort, error) {
	const op = "storage.posgresql.GetContentByIDs"

	if len(ids) == 0 {
		return []models.ContentShort{}, nil
	}

	query := `
		SELECT id, type, title
		FROM content
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
	const op = "storage.postgresql.AnimeDetails"

	query := `
        SELECT content_id, anilist_id,
               mal_id,
               COALESCE(original_title, ''),
               COALESCE(format, ''),
               COALESCE(status, ''),
               COALESCE(season, ''),
               season_year,
               episodes_count,
               episode_duration,
               start_date::text,
               end_date::text,
               COALESCE(language, ''),
               genres, tags, studios, characters, voice_actors,
               mean_score,
               COALESCE(popularity, 0),
               COALESCE(favourites, 0),
               trailer_url,
               raw_data
        FROM anime_details WHERE content_id = $1`

	var a models.AnimeDetails
	err := s.DB.QueryRowContext(ctx, query, id).Scan(
		&a.ContentID,
		&a.AniListID,
		&a.MALID, // *int — принимает NULL
		&a.OriginalTitle,
		&a.Format,
		&a.Status,
		&a.Season,
		&a.SeasonYear,      // *int — принимает NULL
		&a.EpisodesCount,   // *int — принимает NULL
		&a.EpisodeDuration, // *int — принимает NULL
		&a.StartDate,       // *string — принимает NULL
		&a.EndDate,         // *string — принимает NULL
		&a.Language,
		&a.Genres,
		&a.Tags,
		&a.Studios,
		&a.Characters,
		&a.VoiceActors,
		&a.MeanScore,  // *int — принимает NULL
		&a.Popularity, // int — COALESCE нужен
		&a.Favourites, // int — COALESCE нужен
		&a.TrailerURL, // *string — принимает NULL
		&a.RawData,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.AnimeDetails{}, fmt.Errorf("%s: %w", op, storage.ErrNotFound)
		}
		return models.AnimeDetails{}, fmt.Errorf("%s: %w", op, err)
	}
	return a, nil
}

func (s *Storage) MovieDetails(ctx context.Context, id string) (models.MovieDetails, error) {
	const op = "storage.posgresql.MovieDetails"

	query := `
        SELECT content_id, tmdb_id,
               COALESCE(original_title, ''), runtime, COALESCE(tagline, ''),
               COALESCE(status, ''), COALESCE(budget, 0), COALESCE(revenue, 0),
               COALESCE(language, ''),
               genres, cast_members, crew, images, videos, raw_data
        FROM movies_details WHERE content_id = $1`

	var m models.MovieDetails
	err := s.DB.QueryRowContext(ctx, query, id).Scan(
		&m.ContentID, &m.TmdbID,
		&m.OriginalTitle, &m.Runtime, &m.Tagline,
		&m.Status, &m.Budget, &m.Revenue, &m.Language,
		&m.Genres, &m.CastMembers, &m.Crew, &m.Images, &m.Videos, &m.RawData,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.MovieDetails{}, fmt.Errorf("%s: %w", op, storage.ErrNotFound)
		}
		return models.MovieDetails{}, fmt.Errorf("%s: %w", op, err)
	}
	return m, nil
}
func (s *Storage) SeriesDetails(ctx context.Context, id string) (models.SeriesDetails, error) {
	const op = "storage.postgresql.SeriesDetails"

	query := `
        SELECT content_id, tmdb_id,
               COALESCE(original_name, ''),
               COALESCE(status, ''),
               first_air_date::text,
               last_air_date::text,
               COALESCE(number_of_seasons, 0),
               COALESCE(number_of_episodes, 0),
               COALESCE(language, ''),
               genres, networks, cast_members, images, videos,
               raw_data
        FROM series_details WHERE content_id = $1`

	var sd models.SeriesDetails
	err := s.DB.QueryRowContext(ctx, query, id).Scan(
		&sd.ContentID,
		&sd.TmdbID,
		&sd.OriginalName,
		&sd.Status,
		&sd.FirstAirDate,
		&sd.LastAirDate,
		&sd.NumberOfSeasons,
		&sd.NumberOfEpisodes,
		&sd.Language,
		&sd.Genres,
		&sd.Networks,
		&sd.CastMembers,
		&sd.Images,
		&sd.Videos,
		&sd.RawData,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.SeriesDetails{}, fmt.Errorf("%s: %w", op, storage.ErrNotFound)
		}
		return models.SeriesDetails{}, fmt.Errorf("%s: %w", op, err)
	}
	return sd, nil
}

func (s *Storage) BookDetails(ctx context.Context, id string) (models.BookDetails, error) {
	const op = "storage.posgresql.BookDetails"
	query := `SELECT * FROM books_details WHERE id = $1`
	var bookDetails models.BookDetails
	err := s.DB.QueryRowContext(ctx, query, id).Scan(&bookDetails)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.BookDetails{}, fmt.Errorf("%s: %w", op, storage.ErrNotFound)
		}
		return models.BookDetails{}, fmt.Errorf("%s: %w", op, err)
	}
	return bookDetails, nil
}

func (s *Storage) GameDetails(ctx context.Context, id string) (models.GameDetails, error) {
	const op = "storage.posgresql.GameDetails"
	query := `SELECT * FROM games_details WHERE content_id = $1`
	var gameDetails models.GameDetails
	err := s.DB.QueryRowContext(ctx, query, id).Scan(&gameDetails)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.GameDetails{}, fmt.Errorf("%s: %w", op, storage.ErrNotFound)
		}
		return models.GameDetails{}, fmt.Errorf("%s: %w", op, err)
	}
	return gameDetails, nil
}

func (s *Storage) AllAnimeDetails(ctx context.Context) ([]models.AnimeDetails, error) {
	const op = "storage.postgresql.AllAnimeDetails"

	query := `
        SELECT content_id, anilist_id,
               mal_id,
               COALESCE(original_title, ''),
               COALESCE(format, ''),
               COALESCE(status, ''),
               COALESCE(season, ''),
               season_year,
               episodes_count,
               episode_duration,
               start_date::text,
               end_date::text,
               COALESCE(language, ''),
               genres, tags, studios, characters, voice_actors,
               mean_score,
               COALESCE(popularity, 0),
               COALESCE(favourites, 0),
               trailer_url,
               raw_data
        FROM anime_details`

	rows, err := s.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var result []models.AnimeDetails
	for rows.Next() {
		var a models.AnimeDetails
		if err := rows.Scan(
			&a.ContentID,
			&a.AniListID,
			&a.MALID,
			&a.OriginalTitle,
			&a.Format,
			&a.Status,
			&a.Season,
			&a.SeasonYear,
			&a.EpisodesCount,
			&a.EpisodeDuration,
			&a.StartDate,
			&a.EndDate,
			&a.Language,
			&a.Genres,
			&a.Tags,
			&a.Studios,
			&a.Characters,
			&a.VoiceActors,
			&a.MeanScore,
			&a.Popularity,
			&a.Favourites,
			&a.TrailerURL,
			&a.RawData,
		); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		result = append(result, a)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return result, nil
}
func (s *Storage) AllMovieDetails(ctx context.Context) ([]models.MovieDetails, error) {
	const op = "storage.postgresql.AllMovieDetails"

	query := `
        SELECT content_id, tmdb_id,
               COALESCE(original_title, ''),
               runtime,
               COALESCE(tagline, ''),
               COALESCE(status, ''),
               COALESCE(budget, 0),
               COALESCE(revenue, 0),
               COALESCE(language, ''),
               genres, cast_members, crew, images, videos,
               raw_data
        FROM movies_details`

	rows, err := s.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var result []models.MovieDetails
	for rows.Next() {
		var m models.MovieDetails
		if err := rows.Scan(
			&m.ContentID,
			&m.TmdbID,
			&m.OriginalTitle,
			&m.Runtime,
			&m.Tagline,
			&m.Status,
			&m.Budget,
			&m.Revenue,
			&m.Language,
			&m.Genres,
			&m.CastMembers,
			&m.Crew,
			&m.Images,
			&m.Videos,
			&m.RawData,
		); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		result = append(result, m)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return result, nil
}

func (s *Storage) AllSeriesDetails(ctx context.Context) ([]models.SeriesDetails, error) {
	const op = "storage.postgresql.AllSeriesDetails"

	query := `
        SELECT content_id, tmdb_id,
               COALESCE(original_name, ''),
               COALESCE(status, ''),
               first_air_date::text,
               last_air_date::text,
               COALESCE(number_of_seasons, 0),
               COALESCE(number_of_episodes, 0),
               COALESCE(language, ''),
               genres, networks, cast_members, images, videos,
               raw_data
        FROM series_details`

	rows, err := s.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var result []models.SeriesDetails
	for rows.Next() {
		var s models.SeriesDetails
		if err := rows.Scan(
			&s.ContentID,
			&s.TmdbID,
			&s.OriginalName,
			&s.Status,
			&s.FirstAirDate,
			&s.LastAirDate,
			&s.NumberOfSeasons,
			&s.NumberOfEpisodes,
			&s.Language,
			&s.Genres,
			&s.Networks,
			&s.CastMembers,
			&s.Images,
			&s.Videos,
			&s.RawData,
		); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		result = append(result, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return result, nil
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
			return models.Content{}, fmt.Errorf("%s: %w", op, storage.ErrNotFound)
		}

		return models.Content{}, fmt.Errorf("%s: %w", op, err)
	}

	return details, nil
}
