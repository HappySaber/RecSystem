-- +goose Up
-- +goose StatementBegin
CREATE TABLE anime_details (
    content_id UUID PRIMARY KEY REFERENCES content(id) ON DELETE CASCADE,

    anilist_id INT NOT NULL,
    mal_id INT,

    -- Основные метаданные
    original_title TEXT,
    format TEXT,                -- TV | MOVIE | OVA | ONA | SPECIAL
    status TEXT,                -- FINISHED | RELEASING | NOT_YET_RELEASED
    season TEXT,                -- WINTER | SPRING | SUMMER | FALL
    season_year INT,

    episodes_count INT,
    episode_duration INT,       -- в минутах

    start_date DATE,
    end_date DATE,

    language TEXT,              -- ja, en

    -- Контент для рекомендаций
    genres JSONB,               -- ["Action", "Drama"]
    tags JSONB,                 -- ["Time Travel", "Psychological"]
    studios JSONB,              -- [{"id":..., "name":...}]

    characters JSONB,           -- ["Eren Yeager", "Mikasa Ackerman"]
    voice_actors JSONB,         -- ["Yuki Kaji", "Marina Inoue"]

    -- Популярность / скоринг
    mean_score INT,             -- AniList meanScore
    popularity INT,
    favourites INT,

    trailer_url TEXT,

    raw_data JSONB NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS anime_details;
-- +goose StatementEnd
