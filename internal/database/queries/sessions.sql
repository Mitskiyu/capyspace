-- name: CreateSession :one
INSERT INTO sessions (
    id, user_id, revoked, expires_at, created_at
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING id;

-- name: GetSessionExpiration :one
SELECT expires_at FROM sessions
WHERE id = $1;

-- name: GetSession :one
SELECT
    id,
    user_id,
    revoked,
    expires_at,
    created_at
FROM sessions
WHERE id = $1;
