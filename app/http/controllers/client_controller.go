package controllers

import (
	extensions "goravel/app"
	userRequest "goravel/app/http/requests/user_request"
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
	var request userRequest.RegisterUserPostRequest

	errors, err := ctx.Request().ValidateRequest(&request)
	if err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"message": err.Error(),
		})
	}
	if errors != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"message": errors.One(),
		})
	}

	// Password Hashing / Dcrypt
	password, err := facades.Hash().Make(request.Password)
	if err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"message": err.Error(),
		})
	}

	// Start Creating Table
	createUser := models.User{
		IdUser:         extensions.GenerateRandomString(15),
		Name:           request.Name,
		Email:          request.Email,
		Password:       password,
		PasswordString: request.Name+request.Password+request.Email,
		Level:          request.Level,
	}
	createClient := models.Client{
		IdClient: extensions.GenerateRandomString(15),
		IdUser:   createUser.IdUser,
		Status:   request.Status,
	}

	// Begin openning DB Transaction
	ts, err := facades.Orm().Query().Begin()
	if err != nil {
		ts.Rollback()
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"message": err.Error(),
		})
	}
	if err := ts.Create(&createUser); err != nil {
		ts.Rollback()
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"message": err.Error(),
		})
	}
	if err := ts.Create(&createClient); err != nil {
		ts.Rollback()
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"message": err.Error(),
		})
	}

	ts.Commit()
	return ctx.Response().Success().Json(http.Json{
		"Hello": "Client created successfully.",
	})
}
