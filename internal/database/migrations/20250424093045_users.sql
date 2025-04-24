-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY,
    name VARCHAR(36),
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255),
    email_verified TIMESTAMPTZ NOT NULL
);

-- +goose Down
DROP TABLE users;
