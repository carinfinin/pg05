DROP DATABASE IF EXISTS go;
CREATE DATABASE go;

DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS user_data;

\c go;

CREATE TABLE users (
    id SERIAL,
    username VARCHAR PRIMARY KEY
);

CREATE TABLE user_data (
    user_id int not null,
    name VARCHAR(100),
    surname VARCHAR(100),
    description VARCHAR(100)
);