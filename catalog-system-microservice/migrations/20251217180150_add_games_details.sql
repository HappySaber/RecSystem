-- +goose Up
-- +goose StatementBegin
CREATE TABLE games_details (
    content_id UUID PRIMARY KEY REFERENCES content(id) ON DELETE CASCADE,

    igdb_id INT,

    platforms JSONB,
    screenshots JSONB,
    rating NUMERIC,

    genres JSONB,
    developers JSONB,
    publishers JSONB,

    raw_data JSONB NOT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS games_details;
-- +goose StatementEnd
