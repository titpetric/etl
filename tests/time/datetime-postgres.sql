-- PostgreSQL datetime portability test schema
DROP TABLE IF EXISTS datetime_test;

CREATE TABLE datetime_test (
    id SERIAL PRIMARY KEY,
    postgres_date DATE,
    postgres_timestamp TIMESTAMP,
    postgres_timestamptz TIMESTAMPTZ,
    postgres_bigint BIGINT
);
