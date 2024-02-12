package controllers

import (
	extensions "goravel/app"
	"goravel/app/http/requests/user_request"
	"goravel/app/models"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

type ClientController struct {
	//Dependent services
}

func NewClientController() *ClientController {
	return &ClientController{
		//Inject services
	}
}

func (r *ClientController) Index(ctx http.Context) http.Response {
	return nil
}

func (r *ClientController) RegisterUser(ctx http.Context) http.Response {
	// Parse and validate the request
	var request user_request.RegisterUserPostRequest

	errors, err := ctx.Request().ValidateRequest(&request)
	if err != nil {
		return extensions.HandleBadRequestError(ctx, err)
	}
	if errors != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"message": errors.One(),
		})
	}

	// Password Hashing / Dcrypt
	password, err := facades.Hash().Make(request.Password)
	if err != nil {
		return extensions.HandleBadRequestError(ctx, err)
	}

	// Generate random user and client IDs
	userID := extensions.GenerateRandomString(15)
	clientID := extensions.GenerateRandomString(15)

	// Create the user model
	createUser := models.User{
		IdUser:         userID,
		Name:           request.Name,
		Email:          request.Email,
		Password:       password,
		PasswordString: request.Name + request.Password + request.Email,
		Level:          request.Level,
	}
	// Create the client model
	createClient := models.Client{
		IdClient: clientID,
		IdUser:   createUser.IdUser,
		Status:   request.Status,
	}

	// Begin openning DB Transaction
	ts, err := facades.Orm().Query().Begin()
	if err != nil {
		ts.Rollback()
		return extensions.HandleBadRequestError(ctx, err)
	}

	// Create user and client in transaction
	if err := ts.Create(&createUser); err != nil {
		ts.Rollback()
		return extensions.HandleBadRequestError(ctx, err)
	}
	if err := ts.Create(&createClient); err != nil {
		ts.Rollback()
		return extensions.HandleBadRequestError(ctx, err)
	}
	// Commit transaction
	if err := ts.Commit(); err != nil {
		return extensions.HandleBadRequestError(ctx, err)
	}
	// Return success response
	return ctx.Response().Success().Json(http.Json{
		"Hello": "Client created successfully.",
	})
}

