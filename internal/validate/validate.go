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
