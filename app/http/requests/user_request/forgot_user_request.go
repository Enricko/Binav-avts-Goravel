package user_request

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type ForgotUserRequest struct {
	Email                string `form:"email" json:"email"`
	OtpCode              string `form:"otp_code" json:"otp_code"`
	Password             string `form:"password" json:"password"`
	PasswordConfirmation string `form:"password_confirmation" json:"password_confirmation"`
}

func (r *ForgotUserRequest) Authorize(ctx http.Context) error {
	return nil
}

func (r *ForgotUserRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"email":                 "email|max_len:255|required_with:otp_code,password,password_confirmation",
		"otp_code":              "string|max_len:255|required_with:password,password_confirmation",
		"password":              "min_len:8|max_len:255|eq_field:password_confirmation",
		"password_confirmation": "min_len:8|max_len:255|eq_field:password",
	}
}

func (r *ForgotUserRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *ForgotUserRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *ForgotUserRequest) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
