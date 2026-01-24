-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS content_tags (
    content_id UUID NOT NULL,
    tag_id INT NOT NULL,
    PRIMARY KEY (content_id, tag_id),
    CONSTRAINT fk_content FOREIGN KEY (content_id) REFERENCES content(id) ON DELETE CASCADE,
    CONSTRAINT fk_tag FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS content_tags;
-- +goose StatementEnd
