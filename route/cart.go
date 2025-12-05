package route

import (
	apps "be-metalsteel/app/controller"

	"github.com/labstack/echo/v4"
)

func RouteCart(api *echo.Group) {
	api.GET("/cart", apps.GetCart)
	api.POST("/cart", apps.PostCart)
}
