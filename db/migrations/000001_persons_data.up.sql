CREATE SCHEMA IF NOT EXISTS persons_data;

CREATE TABLE IF NOT EXISTS persons_data.genders
(
    id   SMALLSERIAL PRIMARY KEY,
    name VARCHAR(10) NOT NULL UNIQUE
);

INSERT INTO persons_data.genders(name)
VALUES ('male'),
       ('female');

CREATE TABLE IF NOT EXISTS persons_data.countries
(
    id   SMALLSERIAL PRIMARY KEY,
    name VARCHAR(25) NOT NULL UNIQUE
);


CREATE TABLE IF NOT EXISTS persons_data.persons
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(30)               NOT NULL,
    surname    VARCHAR(30)               NOT NULL,
    patronymic VARCHAR(30),
    age        SMALLINT CHECK (age >= 0) NOT NULL,
    gender_id  SMALLINT REFERENCES persons_data.genders (id),
    country_id SMALLINT REFERENCES persons_data.countries (id)
);

