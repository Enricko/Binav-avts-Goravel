package controllers

import (
	extensions "goravel/app"
	"goravel/app/http/requests/user_request"
	"goravel/app/models"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

// UserController handles user-related HTTP requests.
type UserController struct {
	// No dependencies for now
}

// NewUserController creates a new instance of UserController.
func NewUserController() *UserController {
	return &UserController{}
}

// Show is a sample method to demonstrate returning a JSON response.
func (r *UserController) Show(ctx http.Context) http.Response {
	header := ctx.Request().Header("Authorization", "")
	payload, err := facades.Auth().Parse(ctx, header)
	if err != nil {
		return extensions.HandleBadRequestError(ctx, err)
	}

	if payload.Guard != "user" {
		return ctx.Response().Json(http.StatusForbidden, http.Json{
			"message": "The token is not user token.",
			"status":  http.StatusBadRequest,
		})
	}

	var existingUser models.User
	if err := facades.Orm().Query().Model(&existingUser).Where("id_user = ?", payload.Key).First(&existingUser); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"message": err.Error(),
			"status":  http.StatusBadRequest,
		})
	}

	return ctx.Response().Json(http.StatusOK, http.Json{
		"message": "Successfully get data.",
		"user":    existingUser,
		"status":  http.StatusOK,
	})
}

// Login handles user login requests.
func (r *UserController) Login(ctx http.Context) http.Response {
	// Parse the incoming request
	var request user_request.LoginUserRequest

	errors, err := ctx.Request().ValidateRequest(&request)
	if err != nil {
		return extensions.HandleBadRequestError(ctx, err)
	}
	// Validate the request
	if errors != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"message": errors.One(),
		})
	}

	// Find the user by email
	var existingUser models.User
	if err := facades.Orm().Query().Model(&existingUser).Where("email = ?", request.Email).First(&existingUser); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"message": err.Error(),
			"status":  http.StatusBadRequest,
		})
	}

	// Check if the password matches
	if !facades.Hash().Check(request.Password, existingUser.Password) {
		return ctx.Response().Json(http.StatusForbidden, http.Json{
			"message": "Email or Password incorrect.",
			"status":  http.StatusForbidden,
		})
	}

	// Generate JWT token and return response
	token, err := facades.Auth().Login(ctx, &existingUser)
	if err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"message": err.Error(),
			"status":  http.StatusBadRequest,
		})
	}

	return ctx.Response().Json(http.StatusAccepted, http.Json{
		"message": "Welcome back " + existingUser.Name,
		"token":   token,
		"status":  http.StatusAccepted,
	})
}
