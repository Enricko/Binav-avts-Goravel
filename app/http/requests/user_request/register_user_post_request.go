package user_request

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type RegisterUserPostRequest struct {
	Name                 string `form:"name" json:"name"`
	Email                string `form:"email" json:"email"`
	Password             string `form:"password" json:"password"`
	PasswordConfirmation string `form:"password_confirmation" json:"password_confirmation"`
	Level                string `form:"level" json:"level"`
	Status               string `form:"status" json:"status"`
}

func (r *RegisterUserPostRequest) Authorize(ctx http.Context) error {
	return nil
}

func (r *RegisterUserPostRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"name":                  "required|max_len:255",
		"email":                 "required|email|max_len:255",
		"password":              "required|min_len:8|max_len:255|eq_field:password_confirmation",
		"password_confirmation": "required|min_len:8|max_len:255|eq_field:password",
		"level":                 "required|in:client,admin,owner",
		"status":                "required|bool",
	}
}

func (r *RegisterUserPostRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *RegisterUserPostRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *RegisterUserPostRequest) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
