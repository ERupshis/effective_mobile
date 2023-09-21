CREATE SCHEMA IF NOT EXISTS persons_data;

CREATE TABLE IF NOT EXISTS persons_data.gender
(
    id   SMALLSERIAL PRIMARY KEY,
    name VARCHAR(10) NOT NULL
);

CREATE TABLE IF NOT EXISTS persons_data.country
(
    id   SMALLSERIAL PRIMARY KEY,
    name VARCHAR(10) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS persons_data.persons
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(30)               NOT NULL,
    surname    VARCHAR(30)               NOT NULL,
    patronymic VARCHAR(30),
    age        SMALLINT CHECK (age >= 0) NOT NULL,
    gender_id  SMALLINT REFERENCES persons_data.gender (id),
    country_id SMALLINT REFERENCES persons_data.country (id)
);

