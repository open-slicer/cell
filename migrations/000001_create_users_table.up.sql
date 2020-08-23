create table users
(
    id         char(20) primary key,
    email      text     not null unique,
    password   char(60) not null,
    public_key text     not null
);
