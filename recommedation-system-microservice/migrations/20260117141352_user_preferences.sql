-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_genre_preferences (
    user_id UUID NOT NULL,
    genre TEXT NOT NULL,

    score FLOAT NOT NULL DEFAULT 0,
    updated_at TIMESTAMP NOT NULL DEFAULT now(),

    PRIMARY KEY (user_id, genre)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_genre_preferences;
-- +goose StatementEnd
