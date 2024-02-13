package client_request

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type UpdateClientRequest struct {
	Name   string `form:"name" json:"name"`
	Email  string `form:"email" json:"email"`
	Status string `form:"status" json:"status"`
}

func (r *UpdateClientRequest) Authorize(ctx http.Context) error {
	return nil
}

func (r *UpdateClientRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"name":   "required|max_len:255",
		"email":  "required|email|max_len:255",
		"status": "required|bool",
	}
}

func (r *UpdateClientRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *UpdateClientRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *UpdateClientRequest) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
