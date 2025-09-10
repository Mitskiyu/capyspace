-- name: CreateUser :one
INSERT INTO users (
    id, email, password, username, display_name
) VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username_lower = lower($1);
