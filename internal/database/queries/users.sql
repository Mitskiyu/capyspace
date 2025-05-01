-- name: CreateUser :one
INSERT INTO users (
    id, name, email, password, email_verified
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING id;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1
LIMIT 1;
