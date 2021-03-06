CREATE TABLE channels
(
    id     VARCHAR(20) PRIMARY KEY,
    name   VARCHAR(32) NOT NULL,
    owner  VARCHAR(20) NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    parent VARCHAR(20) REFERENCES channels (id) ON DELETE CASCADE
);

CREATE TABLE invites
(
    name    VARCHAR(32) PRIMARY KEY,
    channel VARCHAR(20) NOT NULL REFERENCES channels (id) ON DELETE CASCADE,
    owner   VARCHAR(20) NOT NULL REFERENCES users (id) ON DELETE CASCADE
);
CREATE INDEX ON invites (name, channel);

CREATE TABLE members
(
    id          VARCHAR(20) PRIMARY KEY,
    "user"      VARCHAR(20) NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    channel     VARCHAR(20) NOT NULL REFERENCES channels (id) ON DELETE CASCADE,
    permissions BIGINT
);
CREATE INDEX ON members ("user", channel);

CREATE TABLE messages
(
    id             VARCHAR(20) PRIMARY KEY,
    content_cipher TEXT        NOT NULL,
    channel        VARCHAR(20) NOT NULL REFERENCES channels (id) ON DELETE CASCADE,
    owner          VARCHAR(20) NOT NULL REFERENCES channels (id) ON DELETE CASCADE
);
CREATE INDEX ON messages (channel);