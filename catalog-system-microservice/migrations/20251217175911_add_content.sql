-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE content (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    type TEXT NOT NULL,                 -- movie | series | anime | game | book
    external_source TEXT NOT NULL,      -- tmdb | igdb | anilist | openlibrary
    external_id TEXT NOT NULL,          -- id из внешнего API

    title TEXT NOT NULL,
    description TEXT,
    poster_url TEXT,
    release_date DATE,

    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),

    UNIQUE (external_source, external_id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS content;
-- +goose StatementEnd
