-- SQLite datetime portability test schema
DROP TABLE IF EXISTS datetime_test;

CREATE TABLE datetime_test (
    id INTEGER PRIMARY KEY,
    sqlite_datetime DATETIME,
    sqlite_timestamp TIMESTAMP,
    sqlite_text TEXT,
    sqlite_int INTEGER
);
