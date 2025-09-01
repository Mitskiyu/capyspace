package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

func hashPassword(password string) string {
	salt := make([]byte, 16)
	rand.Read(salt)

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		1,
		19456,
		2,
		32,
	)

	return fmt.Sprintf(
		"%s$%s",
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	)
}

func comparePassword(stored, password string) error {
	sep := strings.Split(stored, "$")
	saltStr, hashStr := sep[0], sep[1]

	salt, err := base64.RawStdEncoding.DecodeString(saltStr)
	if err != nil {
		return err
	}

	hash, err := base64.RawStdEncoding.DecodeString(hashStr)
	if err != nil {
		return err
	}

	passw := argon2.IDKey(
		[]byte(password),
		salt,
		1,
		19456,
		2,
		32,
	)

	if subtle.ConstantTimeCompare(passw, hash) == 0 {
		return fmt.Errorf("password does not match stored hash")
	}

	return nil
}
