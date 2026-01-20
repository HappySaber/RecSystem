-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_content (
    user_id TEXT NOT NULL,
    content_id TEXT NOT NULL,
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
