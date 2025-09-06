-- name: CreateSpace :one
INSERT INTO spaces (
    id, user_id
) VALUES ($1, $2)
RETURNING *;
