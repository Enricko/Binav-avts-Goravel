package controllers

import (
	extensions "goravel/app"
	"goravel/app/http/requests/client_request"
	"goravel/app/http/requests/user_request"
	"goravel/app/models"
	"strings"

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

	// Validate the request
	if shouldReturn, returnValue := extensions.RequestValidation(ctx, &request); shouldReturn {
		return returnValue
	}

	// Password Hashing / Dcrypt
	password, shouldReturn, returnValue := extensions.PasswordHash(request.Password, ctx)
	if shouldReturn {
		return returnValue
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
		// IdUser:   createUser.IdUser,
		Status: request.Status,
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
		"message": "Client created successfully.",
		"status":  http.StatusOK,
	})
}

func (r *ClientController) Update(ctx http.Context) http.Response {
	var request client_request.UpdateClientRequest

	// Validate the request
	if shouldReturn, returnValue := extensions.RequestValidation(ctx, &request); shouldReturn {
		return returnValue
	}

	var existingClient models.Client

	// Begin openning DB Transaction
	ts, err := facades.Orm().Query().Begin()
	if err != nil {
		ts.Rollback()
		return extensions.HandleBadRequestError(ctx, err)
	}

	if err := ts.With("User").Where("id_client=?", ctx.Request().Input("id_client", "")).FirstOrFail(&existingClient); err != nil {
		ts.Rollback()
		return extensions.HandleInternalServerError(ctx, err)
	}

	passwordString := existingClient.User.PasswordString
	passwordString = strings.Replace(passwordString, existingClient.User.Name, "", 1)
	passwordString = strings.Replace(passwordString, existingClient.User.Email, "", 1)

	if _, err := ts.Where("id_client=?", ctx.Request().Input("id_client", "")).
		Update(&models.Client{
			Status: request.Status,
		}); err != nil {
		ts.Rollback()
		return extensions.HandleInternalServerError(ctx, err)
	}
	if _, err := ts.Where("id_user=?", existingClient.IdUser).
		Update(&models.User{
			Name:           request.Name,
			Email:          request.Email,
			PasswordString: request.Name + passwordString + request.Email,
		}); err != nil {
		ts.Rollback()
		return extensions.HandleInternalServerError(ctx, err)
	}
	// Commit transaction
	if err := ts.Commit(); err != nil {
		return extensions.HandleBadRequestError(ctx, err)
	}

	// Return success response
	return ctx.Response().Success().Json(http.Json{
		"message": "Client updated successfully.",
		"status":  http.StatusOK,
	})
}
func (r *ClientController) Delete(ctx http.Context) http.Response {
	var existingClient models.Client

	// Begin openning DB Transaction
	ts, err := facades.Orm().Query().Begin()
	if err != nil {
		ts.Rollback()
		return extensions.HandleBadRequestError(ctx, err)
	}

	if err := ts.FindOrFail(&existingClient, "id_client=?", ctx.Request().Input("id_client", "")); err != nil {
		ts.Rollback()
		return extensions.HandleInternalServerError(ctx, err)
	}
	if _, err := ts.Delete(&existingClient); err != nil {
		ts.Rollback()
		return extensions.HandleInternalServerError(ctx, err)
	}
	if _, err := ts.Delete(&models.User{}, "id_user=?", existingClient.IdUser); err != nil {
		ts.Rollback()
		return extensions.HandleInternalServerError(ctx, err)
	}

	// Commit transaction
	if err := ts.Commit(); err != nil {
		return extensions.HandleBadRequestError(ctx, err)
	}

	// Return success response
	return ctx.Response().Success().Json(http.Json{
		"message": "Client deleted successfully.",
		"status":  http.StatusOK,
	})
}
