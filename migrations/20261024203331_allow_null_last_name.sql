-- +goose Up
ALTER TABLE users
ALTER COLUMN last_name DROP NOT NULL;

-- +goose Down
ALTER TABLE users
ALTER COLUMN last_name SET NOT NULL;
