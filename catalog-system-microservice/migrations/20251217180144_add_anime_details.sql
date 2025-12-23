-- +goose Up
-- +goose StatementBegin
CREATE TABLE anime_details (
    content_id UUID PRIMARY KEY REFERENCES content(id) ON DELETE CASCADE,

    anilist_id INT,
    mal_id INT,

    episodes_count INT,
    status TEXT,
    season TEXT,

    studios JSONB,
    genres JSONB,
    trailer_url TEXT,

    raw_data JSONB NOT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS anime_details;
-- +goose StatementEnd
