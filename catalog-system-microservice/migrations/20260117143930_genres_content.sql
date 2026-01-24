-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS content_genres (
    content_id UUID NOT NULL,
    genre_id INT NOT NULL,
    PRIMARY KEY (content_id, genre_id),
    CONSTRAINT fk_content FOREIGN KEY (content_id) REFERENCES content(id) ON DELETE CASCADE,
    CONSTRAINT fk_genre FOREIGN KEY (genre_id) REFERENCES genres(id) ON DELETE CASCADE
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS content_genres;
-- +goose StatementEnd
