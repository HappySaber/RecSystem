-- +goose Up
-- +goose StatementBegin
INSERT INTO apps (id, name, secret)
VALUES (1, 'test', 'testsecret')
ON CONFLICT DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM apps WHERE id = 1;
-- +goose StatementEnd
