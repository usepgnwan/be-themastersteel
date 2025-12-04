package route

import (
	apps "be-metalsteel/app/controller"

	"github.com/labstack/echo/v4"
)

func RouteProduct(api *echo.Group) {
	api.GET("/product", apps.GetProduct)
	api.POST("/product", apps.PostProduct)
	api.GET("/product/:slug", apps.GetIdProduct)
	api.POST("/product-images", apps.PostProductImage)
}
