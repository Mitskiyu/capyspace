-- name: CreateSession :one
INSERT INTO sessions (
    id, user_id, revoked, expires_at, created_at
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING id;
