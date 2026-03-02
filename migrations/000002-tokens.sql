-- +migrate Up
CREATE TABLE IF NOT EXISTS token (
    uuid      TEXT primary key,
    user_uuid TEXT,
    name      TEXT,
    token     TEXT UNIQUE,
    created_at TEXT
);

-- +migrate Down
DROP TABLE token;
