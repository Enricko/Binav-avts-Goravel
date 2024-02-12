package user_request

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type LoginUserRequest struct {
	Email string `form:"email" json:"email"`
	Password string `form:"password" json:"password"`
}

func (r *LoginUserRequest) Authorize(ctx http.Context) error {
	return nil
}

func (r *LoginUserRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"email":                 "required|email|max_len:255",
		"password":              "required|min_len:8|max_len:255",
	}
}

func (r *LoginUserRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *LoginUserRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *LoginUserRequest) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
