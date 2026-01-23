-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_content_action (
    user_id UUID NOT NULL,
    content_id UUID NOT NULL,

    action_id SMALLINT NOT NULL,
    CONSTRAINT fk_action
        FOREIGN KEY (action_id)
        REFERENCES actions(id)
        ON DELETE RESTRICT,

    rating INT,
    duration_sec INT,

    created_at TIMESTAMP NOT NULL DEFAULT now(),

    PRIMARY KEY (user_id, content_id, action_id, created_at)
);

CREATE INDEX idx_uca_user
    ON user_content_action (user_id);

CREATE INDEX idx_uca_content
    ON user_content_action (content_id);

CREATE INDEX idx_uca_action
    ON user_content_action (action_id);

CREATE INDEX idx_uca_created
    ON user_content_action (created_at);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_content_action;
-- +goose StatementEnd
