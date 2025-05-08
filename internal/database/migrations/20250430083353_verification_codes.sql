-- +goose Up
CREATE TABLE verification_codes (
    id UUID PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    code VARCHAR(8) NOT NULL,
    used TIMESTAMPTZ,
    expires_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_email ON verification_codes (email);

-- +goose Down
DROP TABLE verification_codes;
