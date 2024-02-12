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
		router.Post("/login", userController.Login)
		router.Middleware(userMiddleware.Authorization()).Group(func(router route.Router) {
			router.Get("/user", userController.Show)
			router.Get("/logout", userController.Logout)
		})

		// ===== Forget Password =====
		router.Post("/forgot-password", userController.ForgotPassword)
		router.Post("/check-code", userController.CheckCode)
		router.Post("/reset-password", userController.ResetPassword)

		// ===== Client =====
		clientController := controllers.NewClientController()
		// router.Middleware(userMiddleware.Authorization()).Group(func(router route.Router) {
		router.Group(func(router route.Router) {
			router.Post("/register-user", clientController.RegisterUser)
		})
	})
	facades.Route().Fallback(func(ctx http.Context) http.Response {
		return ctx.Response().String(404, "not foundasdasd")
	})
}
