-- +goose Up
-- +goose StatementBegin
CREATE TABLE actions (
    id SMALLSERIAL PRIMARY KEY,
    code VARCHAR(32) NOT NULL UNIQUE -- VIEW, LIKE, DISLIKE, RATE, FAVORITE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS actions;
-- +goose StatementEnd
