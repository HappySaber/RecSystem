package postgresql

import (
	"catalog-microservice/internal/domain/models"
	"catalog-microservice/internal/storage"
	"errors"
	"fmt"

	"github.com/lib/pq"
)

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

func (s *Storage) SaveSeries(series *models.SeriesDetails) error {
	_, err := s.DB.Exec(`
	INSERT INTO series_details (
		content_id,
		tmdb_id,
		original_name,
		status,
		first_air_date,
		last_air_date,
		number_of_seasons,
		number_of_episodes,
		language,
		genres,
		networks,
		cast_members,
		images,
		videos,
		raw_data
	)
	VALUES (
		$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15
	)
	ON CONFLICT (content_id) DO UPDATE SET
		tmdb_id = EXCLUDED.tmdb_id,
		original_name = EXCLUDED.original_name,
		status = EXCLUDED.status,
		first_air_date = EXCLUDED.first_air_date,
		last_air_date = EXCLUDED.last_air_date,
		number_of_seasons = EXCLUDED.number_of_seasons,
		number_of_episodes = EXCLUDED.number_of_episodes,
		language = EXCLUDED.language,
		genres = EXCLUDED.genres,
		networks = EXCLUDED.networks,
		cast_members = EXCLUDED.cast_members,
		images = EXCLUDED.images,
		videos = EXCLUDED.videos,
		raw_data = EXCLUDED.raw_data
	`,
		series.ContentID,
		series.TmdbID,
		series.OriginalName,
		series.Status,
		series.FirstAirDate,
		series.LastAirDate,
		series.NumberOfSeasons,
		series.NumberOfEpisodes,
		series.Language,
		series.Genres,
		series.Networks,
		series.CastMembers,
		series.Images,
		series.Videos,
		series.RawData,
	)

	return err
}

func (s *Storage) SaveAnime(anime *models.AnimeDetails) error {
	_, err := s.DB.Exec(`
	INSERT INTO anime_details (
		content_id,
		anilist_id,
		mal_id,

		original_title,
		format,
		status,
		season,
		season_year,

		episodes_count,
		episode_duration,

		start_date,
		end_date,

		language,

		genres,
		tags,
		studios,

		characters,
		voice_actors,

		mean_score,
		popularity,
		favourites,

		trailer_url,
		raw_data
	)
	VALUES (
		$1, $2, $3,
		$4, $5, $6, $7, $8,
		$9, $10,
		$11, $12,
		$13,
		$14, $15, $16,
		$17, $18,
		$19, $20, $21,
		$22, $23
	)
	ON CONFLICT (content_id) DO UPDATE SET
		anilist_id       = EXCLUDED.anilist_id,
		mal_id           = EXCLUDED.mal_id,

		original_title   = EXCLUDED.original_title,
		format           = EXCLUDED.format,
		status           = EXCLUDED.status,
		season           = EXCLUDED.season,
		season_year      = EXCLUDED.season_year,

		episodes_count   = EXCLUDED.episodes_count,
		episode_duration = EXCLUDED.episode_duration,

		start_date       = EXCLUDED.start_date,
		end_date         = EXCLUDED.end_date,

		language         = EXCLUDED.language,

		genres           = EXCLUDED.genres,
		tags             = EXCLUDED.tags,
		studios          = EXCLUDED.studios,

		characters       = EXCLUDED.characters,
		voice_actors     = EXCLUDED.voice_actors,

		mean_score       = EXCLUDED.mean_score,
		popularity       = EXCLUDED.popularity,
		favourites       = EXCLUDED.favourites,

		trailer_url      = EXCLUDED.trailer_url,
		raw_data         = EXCLUDED.raw_data
	`,
		anime.ContentID,
		anime.AniListID,
		anime.MALID,

		anime.OriginalTitle,
		anime.Format,
		anime.Status,
		anime.Season,
		anime.SeasonYear,

		anime.EpisodesCount,
		anime.EpisodeDuration,

		anime.StartDate,
		anime.EndDate,

		anime.Language,

		anime.Genres,
		anime.Tags,
		anime.Studios,

		anime.Characters,
		anime.VoiceActors,

		anime.MeanScore,
		anime.Popularity,
		anime.Favourites,

		anime.TrailerURL,
		anime.RawData,
	)

	return err
}
