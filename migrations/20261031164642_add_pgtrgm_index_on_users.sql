-- +goose NO TRANSACTION

-- +goose Up
CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX CONCURRENTLY idx_users_first_name_trgm ON users USING gin (first_name gin_trgm_ops);
CREATE INDEX CONCURRENTLY idx_users_last_name_trgm ON users USING gin (last_name gin_trgm_ops);
CREATE INDEX CONCURRENTLY idx_users_email_trgm ON users USING gin (email gin_trgm_ops);

-- +goose Down
DROP INDEX CONCURRENTLY IF EXISTS idx_users_name_trgm
DROP INDEX CONCURRENTLY IF EXISTS idx_users_email_trgm
