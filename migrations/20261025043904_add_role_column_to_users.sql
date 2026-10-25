-- +goose Up
CREATE TYPE roles AS ENUM ('student', 'teacher', 'admin');
ALTER TABLE users
ADD COLUMN role roles NOT NULL DEFAULT 'student';

-- +goose Down
ALTER TABLE users DROP COLUMN role;
DROP TYPE roles;
