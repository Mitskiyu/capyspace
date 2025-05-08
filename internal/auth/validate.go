package auth

import (
	"fmt"
	"net/mail"
	"strings"
)

// Validates the email according to RFC 5322
func validateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("empty")
	}

	_, err := mail.ParseAddress(email)
	if err != nil {
		return err
	}

	return nil
}

// Validates the verification code to be 8 digits
func validateVerificationCode(code string) error {
	if code == "" {
		return fmt.Errorf("empty")
	}

	if len(code) != 8 {
		return fmt.Errorf("length")
	}

	for _, c := range code {
		if c < '0' || c > '9' {
			return fmt.Errorf("not numeric")
		}
	}

	return nil
}

// Validates the password to be at least 8 characters long
func validatePassword(pw string) error {
	if pw == "" {
		return fmt.Errorf("empty")
	}

	if len(pw) > 255 || len(pw) < 8 {
		return fmt.Errorf("length")
	}

	if strings.TrimSpace(pw) == "" {
		return fmt.Errorf("pw validation error: whitespace")
	}

	return nil
}
