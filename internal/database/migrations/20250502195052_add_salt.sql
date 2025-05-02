-- +goose Up
ALTER TABLE users
ADD salt VARCHAR(24);

-- +goose Down
ALTER TABLE users
DROP COLUMN salt;
