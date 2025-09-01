package auth

import (
	"crypto/rand"
	"encoding/base64"
)

func createSessionId() string {
	token := make([]byte, 15)
	rand.Read(token)
	sessionId := base64.StdEncoding.EncodeToString(token)

	return sessionId
}
