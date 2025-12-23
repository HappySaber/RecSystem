-- +goose Up
-- +goose StatementBegin
CREATE TABLE series_details (
    content_id UUID PRIMARY KEY REFERENCES content(id) ON DELETE CASCADE,

    tmdb_id INT NOT NULL,

    seasons JSONB,
    episodes JSONB,
    networks JSONB,
    genres JSONB,
    cast_members JSONB,
    images JSONB,

    raw_data JSONB NOT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS series_details;
-- +goose StatementEnd
