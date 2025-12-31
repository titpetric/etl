-- MySQL datetime portability test schema
DROP TABLE IF EXISTS datetime_test;

CREATE TABLE datetime_test (
    id INT PRIMARY KEY AUTO_INCREMENT,
    mysql_date DATE,
    mysql_datetime DATETIME,
    mysql_datetime_fsp DATETIME(6),
    mysql_timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    mysql_bigint BIGINT
);
