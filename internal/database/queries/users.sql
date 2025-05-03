-- name: CreateUser :one
INSERT INTO users (
    id, name, email, password, salt, email_verified
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING id;

-- name: GetUserByEmail :one
SELECT
    id,
    name,
    email,
    password,
    salt,
    email_verified
FROM users
WHERE email = $1;

-- name: GetUser :one
SELECT
    id,
    name,
    email,
    password,
    salt,
    email_verified
FROM users
WHERE id = $1;
