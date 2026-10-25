-- +goose Up
ALTER TABLE users
ADD COLUMN department TEXT;

-- +goose Down
ALTER TABLE users DROP COLUMN department;
