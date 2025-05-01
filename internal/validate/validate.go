package validate

import (
	"fmt"
	"net/mail"
	"strings"
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

func VerificationCode(code string) error {
	if code == "" {
		return fmt.Errorf("code validation error: empty")
	}

	if len(code) != 6 {
		return fmt.Errorf("code validation error: length")
	}

	for _, c := range code {
		if c < '0' || c > '9' {
			return fmt.Errorf("code validation error: must be digits")
		}
	}

	return nil
}

func Password(pw string) error {
	if pw == "" {
		return fmt.Errorf("pw validation error: empty")
	}

	if len(pw) > 255 || len(pw) < 8 {
		return fmt.Errorf("pw validation error: length")
	}

	if strings.TrimSpace(pw) == "" {
		return fmt.Errorf("pw validation error: whitespace")
	}

	return nil
}
