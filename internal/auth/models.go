package auth

import (
	"github.com/Mitskiyu/capyspace/internal/user"
)

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type EmailReq struct {
	Email string `json:"email"`
}

type EmailRes struct {
	Exists bool `json:"exists"`
}

type RegisterReq struct {
	Credentials
}

type LoginReq struct {
	Credentials
}

type LoginRes struct {
	Message string    `json:"message"`
	User    user.Info `json:"user"`
}
