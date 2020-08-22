CREATE TABLE users (
    id CHAR(20) PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password CHAR(60) NOT NULL,
    public_key TEXT NOT NULL
);