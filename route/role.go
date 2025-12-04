package route

import (
	apps "be-metalsteel/app/controller"

	"github.com/labstack/echo/v4"
)

func RouteRole(api *echo.Group) {
	api.GET("/role", apps.GetRole)
}
