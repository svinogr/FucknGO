CREATE DATABASE db;
CREATE USER user_test WITH PASSWORD '123456';
/*GRANT ALL PRIVILEGES ON DATABASE db to user_test;
CREATE TABLE if not exists users
(
    Id SERIAL PRIMARY KEY,
    user_name CHARACTER VARYING(30),
    password CHARACTER VARYING(30),
    Email CHARACTER VARYING(30),
    UNIQUE(Id)
    );

CREATE TABLE if not exists tokens
(
    Id SERIAL PRIMARY KEY,
    token  CHARACTER VARYING(250),
    user_id INTEGER references users (Id),
    UNIQUE(Id)
    );

*/