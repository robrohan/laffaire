-- +migrate Up
CREATE TABLE IF NOT EXISTS event (
    uuid TEXT primary key
    , user_uuid TEXT
    , title TEXT
    , description TEXT
);

CREATE TABLE IF NOT EXISTS entry (
    uuid TEXT primary key
    , event_uuid TEXT
    , subject TEXT
    , start_date TEXT -- as ISO8601 strings ("YYYY-MM-DD HH:MM:SS.SSS")
    , start_time TEXT -- as ISO8601 strings ("YYYY-MM-DD HH:MM:SS.SSS")
    , end_date TEXT -- as ISO8601 strings ("YYYY-MM-DD HH:MM:SS.SSS")
    , end_time TEXT -- as ISO8601 strings ("YYYY-MM-DD HH:MM:SS.SSS")
    , all_day_event INTEGER
    , description TEXT
    , location TEXT
    , private INTEGER
);

-- +migrate Down
DROP TABLE event;
DROP TABLE entry;
