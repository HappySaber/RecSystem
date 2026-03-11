-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS content_genres (
    content_id UUID NOT NULL,
    genre_id INT REFERENCES genres(id),

    PRIMARY KEY (content_id, genre_id)
);

CREATE INDEX idx_content_genres_content
ON content_genres(content_id);

CREATE INDEX idx_content_genres_genre
ON content_genres(genre_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS content_genres;
-- +goose StatementEnd
