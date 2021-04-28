CREATE DATABASE db;
CREATE USER user_test WITH PASSWORD '123456';
GRANT ALL PRIVILEGES ON DATABASE db to user_test;
CREATE TABLE if not exists users
(
    Id        SERIAL PRIMARY KEY,
    user_name CHARACTER VARYING(30),
    password  CHARACTER VARYING(30),
    Email     CHARACTER VARYING(30),
    UNIQUE (Id)
);

ALTER TABLE users ADD COLUMN type CHARACTER VARYING(30);

CREATE TABLE if not exists tokens
(
    Id      SERIAL PRIMARY KEY,
    token   CHARACTER VARYING(250),
    user_id INTEGER references users (Id),
    UNIQUE (Id)
);

CREATE TABLE if not exists refresh_sessions
(
    "id"            SERIAL PRIMARY KEY,
    "user_id"       INTEGER,
    "refresh_token" character varying(200)   NOT NULL,
    "user_agent"    character varying(200)   NOT NULL, /* user-agent */
    "fingerprint"   character varying(200)   NOT NULL,
    "ip"            character varying(15)    NOT NULL,
    "expires_in"    timestamp                NOT NULL,
    "created_at"    timestamp with time zone NOT NULL DEFAULT now()
);

CREATE TABLE if not exists coord
(
    "id"      SERIAL PRIMARY KEY,
    "coord_x" double precision,
    "coord_y" double precision
);

CREATE TABLE if not exists shops
(
    "id"      SERIAL PRIMARY KEY,
    "user_id" INTEGER,
    "coord_id" INTEGER,
    "name"  character varying(200)   NOT NULL,
    "address"  character varying(200)   NOT NULL,
    UNIQUE (Id)
);

CREATE TABLE if not exists shops_stock
(
    "id"      SERIAL PRIMARY KEY,
    "shop_id" INTEGER,
    "title"  character varying(200)   NOT NULL,
    "description"  character varying(200)   NOT NULL,
    "date_start" timestamp NOT NULL,
    "date_finish" timestamp NOT NULL,
    UNIQUE (Id)
);

CREATE TABLE if not exists product
(
    "id"      SERIAL PRIMARY KEY,
    "shop_stock_id" INTEGER,
    "title" character varying(200)   NOT NULL,
    "description"  character varying(200)   NOT NULL,
    "price"  character varying(200)   NOT NULL,
    "price_old"  character varying(200)   NOT NULL,
    UNIQUE (Id)
);

