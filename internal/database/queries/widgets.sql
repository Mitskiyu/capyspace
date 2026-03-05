-- name: CreateWidget :one
INSERT INTO widgets (
    id, space_id, type, x_pos, y_pos, minimized, data
) VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: UpdateWidget :one
UPDATE widgets
SET x_pos = $2, y_pos = $3, minimized = $4, data = $5
WHERE id = $1
RETURNING *;
