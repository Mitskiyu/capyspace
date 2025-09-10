package user

import (
	"context"

	"github.com/Mitskiyu/capyspace/internal/util"
)

type EmailReq struct {
	Email string `json:"email"`
}

type UsernameReq struct {
	Username string `json:"username"`
}

type Exists struct {
	Exists bool `json:"exists"`
}

type Info struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}

func (e EmailReq) Valid(ctx context.Context) map[string]string {
	problems := make(map[string]string)

	if err := util.ValidEmail(e.Email); err != nil {
		problems["email"] = err.Error()
	}

	return problems
}

func (u UsernameReq) Valid(ctx context.Context) map[string]string {
	problems := make(map[string]string)

	if err := util.ValidUsername(u.Username); err != nil {
		problems["username"] = err.Error()
	}

	return problems
}
