-- name: CreateSpace :one
INSERT INTO spaces (
    id, user_id
) VALUES ($1, $2)
RETURNING *;

-- name: GetSpaceByUsername :one
SELECT s.*
FROM spaces s
JOIN users u ON s.user_id = u.id
WHERE u.username = $1;
