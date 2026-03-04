package auth

import (
	"crypto/rand"
	"encoding/base64"
)

func createSessionID() string {
	token := make([]byte, 15)
	rand.Read(token)
	sessionID := base64.StdEncoding.EncodeToString(token)

	return sessionID
}
