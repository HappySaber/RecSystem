-- +goose Up
-- +goose StatementBegin
CREATE TABLE series_details (
    content_id UUID PRIMARY KEY REFERENCES content(id) ON DELETE CASCADE,

    tmdb_id INT NOT NULL,

    original_name TEXT,
    status TEXT,                -- Returning Series / Ended
    first_air_date DATE,
    last_air_date DATE,

    number_of_seasons INT,
    number_of_episodes INT,

    language TEXT,

    genres JSONB,               -- ["Drama", "Crime"]
    networks JSONB,             -- [{"id":..., "name":...}]
    cast_members JSONB,         -- ["Actor 1", "Actor 2"]

    images JSONB,
    videos JSONB,

    raw_data JSONB NOT NULL
);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS series_details;
-- +goose StatementEnd
