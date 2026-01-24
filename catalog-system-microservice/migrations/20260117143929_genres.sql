-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS genres (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS genres;
-- +goose StatementEnd
