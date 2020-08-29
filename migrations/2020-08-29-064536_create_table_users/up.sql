CREATE TABLE users (
  id CHAR(21) PRIMARY KEY,
  username VARCHAR(32) NOT NULL UNIQUE,
  display_name VARCHAR(32) NOT NULL,
  public_key TEXT NOT NULL
);