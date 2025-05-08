package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type hashParams struct {
	memory  uint32
	iter    uint32
	threads uint8
	saltLen uint32
	keyLen  uint32
}

type SessionClaims struct {
	SessionID    uuid.UUID `json:"session_id"`
	UserID       uuid.UUID `json:"user_id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	RevalidateAt int64     `json:"revalidate_at"`
	jwt.RegisteredClaims
}
