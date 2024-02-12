package user_request

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type ForgotUserRequest struct {
	Email   string `form:"email" json:"email"`
	OtpCode string `form:"otp_code" json:"otp_code"`
}

func (r *ForgotUserRequest) Authorize(ctx http.Context) error {
	return nil
}

func (r *ForgotUserRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"email":    "email|max_len:255|required_with:otp_code",
		"otp_code": "string|max_len:255",
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
