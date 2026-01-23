-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS content_popularity_agg (
    content_id UUID PRIMARY KEY,

    views INT NOT NULL DEFAULT 0,
    likes INT NOT NULL DEFAULT 0,
    favorites INT NOT NULL DEFAULT 0,

    score FLOAT NOT NULL DEFAULT 0,
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS content_popularity_agg;
-- +goose StatementEnd
