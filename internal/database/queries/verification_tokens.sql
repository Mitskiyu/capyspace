-- name: CreateVerificationToken :exec
INSERT INTO verification_tokens (
    id, email, token, used, expires_at
) VALUES (
    $1, $2, $3, $4, $5
);
