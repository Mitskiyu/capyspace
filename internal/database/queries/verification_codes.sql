-- name: CreateVerificationCode :exec
INSERT INTO verification_codes (
    id, email, code, used, expires_at
) VALUES (
    $1, $2, $3, $4, $5
);

-- name: GetValidVerificationCode :one
SELECT id
FROM verification_codes
WHERE
    email = $1
    AND code = $2
    AND used = FALSE
    AND expires_at > NOW()
LIMIT 1;

-- name: GetUsedVerificationCode :one
SELECT 1
FROM verification_codes
WHERE email = $1 AND used = TRUE
LIMIT 1;

-- name: SetUsedVerificationCode :exec
UPDATE verification_codes
SET used = TRUE
WHERE
    id = $1;
