-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_content_action (
    id BIGSERIAL PRIMARY KEY,

    user_id UUID NOT NULL,
    content_id UUID NOT NULL,

    action_id SMALLINT NOT NULL,
    rating INT,
    duration_sec INT,

    created_at TIMESTAMP NOT NULL DEFAULT now(),

    CONSTRAINT fk_user_content_action_action
        FOREIGN KEY (action_id)
        REFERENCES actions(id)
        ON DELETE RESTRICT,

    CONSTRAINT rating_check
        CHECK (rating IS NULL OR rating BETWEEN 1 AND 10),

    CONSTRAINT duration_check
        CHECK (duration_sec IS NULL OR duration_sec >= 0)
);

CCREATE INDEX idx_user_content_action_user
ON user_content_action(user_id);

CREATE INDEX idx_user_content_action_content
ON user_content_action(content_id);

CREATE INDEX idx_user_content_action_created
ON user_content_action(created_at);

CREATE INDEX idx_user_content_action_action
ON user_content_action(action_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_content_action;
-- +goose StatementEnd
