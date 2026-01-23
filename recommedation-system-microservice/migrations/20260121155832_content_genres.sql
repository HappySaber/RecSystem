-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS content_genres (
    content_id UUID NOT NULL,
    genre TEXT NOT NULL,

    PRIMARY KEY (content_id, genre)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS content_genres;
-- +goose StatementEnd
