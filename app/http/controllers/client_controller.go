package controllers

import (
	userRequest "goravel/app/http/requests/user_request"

	"github.com/goravel/framework/contracts/http"
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

	return ctx.Response().Success().Json(http.Json{
		"Hello": "Goravel",
	})
}
