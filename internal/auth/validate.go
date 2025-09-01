package auth

import (
	"context"
	"fmt"
	"net/mail"
)

func validEmail(email string) error {
	if len(email) > 255 {
		return fmt.Errorf("email cannot be longer than 255 chars")
	}

	if _, err := mail.ParseAddress(email); err != nil {
		return fmt.Errorf("email is not a valid address")
	}

	return nil
}

func validPassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password cannot be less than 8 chars")
	}

	if len(password) > 64 {
		return fmt.Errorf("password cannot be longer than 64 chars")
	}

	return nil
}

func (c Credentials) Valid(ctx context.Context) map[string]string {
	problems := make(map[string]string)

	if err := validEmail(c.Email); err != nil {
		problems["email"] = err.Error()
	}

	if err := validPassword(c.Password); err != nil {
		problems["password"] = err.Error()
	}

	return problems
}
