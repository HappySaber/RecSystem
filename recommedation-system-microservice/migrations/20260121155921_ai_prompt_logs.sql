-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS ai_prompt_logs (
    id BIGSERIAL PRIMARY KEY,

    user_id UUID,
    prompt TEXT NOT NULL,
    response TEXT,

    created_at TIMESTAMP NOT NULL DEFAULT now()
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS ai_prompt_logs;
-- +goose StatementEnd
