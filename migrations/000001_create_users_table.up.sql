CREATE TABLE users
(
    id            CHAR(20) PRIMARY KEY,
    username      VARCHAR(32) NOT NULL UNIQUE,
    display_name  VARCHAR(32),
    password_hash BYTEA       NOT NULL,
    public_key    BYTEA       NOT NULL
);