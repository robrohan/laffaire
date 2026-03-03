-- +migrate Up
ALTER TABLE users ADD COLUMN timezone TEXT DEFAULT 'UTC';

-- +migrate Down
-- SQLite doesn't support DROP COLUMN; leave as no-op
