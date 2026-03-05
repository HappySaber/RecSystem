-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS content_genres (
    content_id UUID NOT NULL,
    genre_id INT REFERENCES genres(id),

    PRIMARY KEY (content_id, genre)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS content_genres;
-- +goose StatementEnd
