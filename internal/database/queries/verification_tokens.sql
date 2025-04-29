-- name: CreateVerificationToken :exec
INSERT INTO verification_tokens (
    id, email, token, used, expires_at
) VALUES (
    $1, $2, $3, $4, $5
);

-- name: GetValidVerificationToken :one
SELECT
    id,
    email,
    token,
    used,
    expires_at
FROM verification_tokens
WHERE
    email = $1
    AND token = $2
    AND used = FALSE
    AND expires_at > NOW()
LIMIT 1;

-- name: SetUsedVerificationToken :exec
UPDATE verification_tokens
SET used = TRUE
WHERE
    id = $1;
