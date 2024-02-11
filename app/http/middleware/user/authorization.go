package userMiddleware

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

func Authorization() http.Middleware {

	return func(ctx http.Context) {
		header := ctx.Request().Header("Authorization", "")
		_, err := facades.Auth().Parse(ctx, header)
		if err != nil {
			ctx.Request().AbortWithStatus(401)
			return
		}
		ctx.Request().Next()
	}
}
