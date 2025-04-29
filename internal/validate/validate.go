package validate

import (
	"fmt"
	"net/mail"
)

func Email(email string) error {
	if email == "" {
		return fmt.Errorf("email validation error: empty")
	}

	_, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("email validation error: %v", err)
	}

	return nil
}

func VerificationToken(token string) error {
	if token == "" {
		return fmt.Errorf("token validation error: empty")
	}

	if len(token) != 6 {
		return fmt.Errorf("token validation error: length")
	}

	for _, c := range token {
		if c < '0' || c > '9' {
			return fmt.Errorf("token validation error: must be digits")
		}
	}

	return nil
}
