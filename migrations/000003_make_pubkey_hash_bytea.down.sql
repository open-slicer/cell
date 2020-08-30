DROP TABLE IF EXISTS users CASCADE;

CREATE TABLE users
(
    id            CHAR(20) PRIMARY KEY,
    username      VARCHAR(32) NOT NULL UNIQUE,
    display_name  VARCHAR(32),
    password_hash CHAR(60)    NOT NULL,
    public_key    TEXT        NOT NULL
);