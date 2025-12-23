-- +goose Up
-- +goose StatementBegin
CREATE TABLE books_details (
    content_id UUID PRIMARY KEY REFERENCES content(id) ON DELETE CASCADE,

    isbn_10 TEXT,
    isbn_13 TEXT,

    authors JSONB,
    publishers JSONB,
    pages INT,
    language TEXT,

    raw_data JSONB NOT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS books_details;
-- +goose StatementEnd
