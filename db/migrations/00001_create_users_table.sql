-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS db (id SERIAL PRIMARY KEY);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE db;

-- +goose StatementEnd
