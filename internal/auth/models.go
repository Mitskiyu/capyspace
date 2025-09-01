package auth

import (
	"github.com/Mitskiyu/capyspace/internal/user"
)

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterReq struct {
	Credentials
}

type LoginReq struct {
	Credentials
}

type LoginRes struct {
	Message string `json:"message"`
	User    user.Info
}
