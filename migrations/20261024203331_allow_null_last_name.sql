-- +goose Up
ALTER TABLE user
ALTER COLUMN last_name DROP NOT NULL;

-- +goose Down
ALTER TABLE user
ALTER COLUMN last_name SET NOT NULL;
