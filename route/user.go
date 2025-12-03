package route

import (
	apps "be-metalsteel/app/controller"

	"github.com/labstack/echo/v4"
)

func RouteUser(api *echo.Group) {
	api.GET("/user", apps.GetUser)
	api.POST("/user", apps.PostUser)
	api.POST("/user/login", apps.UserLogin)
}
