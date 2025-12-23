-- +goose Up
-- +goose StatementBegin
CREATE TABLE movies_details (
    content_id UUID PRIMARY KEY REFERENCES content(id) ON DELETE CASCADE,

    tmdb_id INT NOT NULL,

    original_title TEXT,
    runtime INT,
    tagline TEXT,
    status TEXT,
    budget BIGINT,
    revenue BIGINT,
    language TEXT,

    genres JSONB,
    cast_members JSONB,
    crew JSONB,
    images JSONB,
    videos JSONB,

    raw_data JSONB NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS movies_details;
-- +goose StatementEnd
