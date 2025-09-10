package auth

import (
	"context"

	"github.com/Mitskiyu/capyspace/internal/user"
	"github.com/Mitskiyu/capyspace/internal/util"
)

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterReq struct {
	Credentials
	Username string `json:"username"`
}

type LoginReq struct {
	Credentials
}

type LoginRes struct {
	Message string    `json:"message"`
	User    user.Info `json:"user"`
}

func (c Credentials) Valid(ctx context.Context) map[string]string {
	problems := make(map[string]string)

	if err := util.ValidEmail(c.Email); err != nil {
		problems["email"] = err.Error()
	}

	if err := util.ValidPassword(c.Password); err != nil {
		problems["password"] = err.Error()
	}

	return problems
}

func (r RegisterReq) Valid(ctx context.Context) map[string]string {
	problems := r.Credentials.Valid(ctx)

	if err := util.ValidUsername(r.Username); err != nil {
		problems["username"] = err.Error()
	}

	return problems
}
