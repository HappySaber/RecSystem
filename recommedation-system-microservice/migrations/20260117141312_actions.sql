-- +goose Up
-- +goose StatementBegin
CREATE TABLE actions (
    id SMALLSERIAL PRIMARY KEY,
    code VARCHAR(32) NOT NULL UNIQUE, -- VIEW, LIKE, DISLIKE, RATE, FAVORITE
    weight INT NOT NULL DEFAULT 1
);

INSERT INTO actions (code, weight) VALUES
('VIEW', 1),
('LIKE', 3),
('DISLIKE', -3),
('RATE', 2),
('FAVORITE', 5);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS actions;
-- +goose StatementEnd
