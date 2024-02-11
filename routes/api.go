package routes

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/route"
	"github.com/goravel/framework/facades"

	"goravel/app/http/controllers"
	userMiddleware "goravel/app/http/middleware/user"
)

func Api() {
	facades.Route().Prefix("api").Group(func(router route.Router) {
		// ===== Users =====
		userController := controllers.NewUserController()
		router.Get("/users/{id}", userController.Show)

		// ===== Client =====
		clientController := controllers.NewClientController()
		router.Middleware(userMiddleware.Authorization()).Group(func(router route.Router) {
			router.Post("/register_user", clientController.RegisterUser)
		})

	})
	facades.Route().Fallback(func(ctx http.Context) http.Response {
		return ctx.Response().String(404, "not foundasdasd")
	})
}
