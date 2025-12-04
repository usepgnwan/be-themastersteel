package route

import (
	"net/http"
	"os"

	"strings"

	"be-metalsteel/app/helpers"
	_ "be-metalsteel/docs"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func InitRouting(e *echo.Echo) {

	// midleware bearer
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(helpers.MyJwt)
		},
		SigningKey: []byte(os.Getenv("SECRETKEYUSER")),
		ErrorHandler: func(c echo.Context, err error) error {
			var errorMessage string

			// Check the type or content of the error to handle different scenarios
			if err != nil {
				switch err.Error() {
				case "missing or malformed jwt":
					errorMessage = "JWT is missing or malformed"
				case "invalid or expired jwt":
					errorMessage = "JWT is invalid or expired"
				default:
					errorMessage = "An error occurred with the JWT token"
				}
			} else {
				errorMessage = "An unknown error occurred"
			}

			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"message": errorMessage,
				"status":  false,
			})
		},
	}

	// Swagger route
	// e.GET("/swagger/*", echoSwagger.WrapHandler)
	// e.Static("/uploads", "uploads")
	// Route group for Swagger UI with BasicAuth middleware applied
	swaggerGroup := e.Group("/documentation")
	swaggerGroup.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		// Validate username and password
		if username == os.Getenv("USER_SWAG") && password == os.Getenv("PASS_SWAG") {
			return true, nil
		}
		return false, nil
	}))

	swaggerGroup.GET("/swagger/*", echoSwagger.WrapHandler)
	e.Static("/uploads", "uploads")
	api := e.Group("/api")
	api.Use(HeaderAuthorizationMiddleware)

	jwtGroup := api.Group("")
	jwtGroup.Use(echojwt.WithConfig(config))

	RouteUser(api)
	RouteRole(api)
	RouteProduct(api)
}

func HeaderAuthorizationMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		headerValue := c.Request().Header.Get("x-tes-themastersteel")

		if headerValue == "" {
			return c.JSON(http.StatusForbidden, map[string]interface{}{
				"message": "API can't use",
				"status":  false,
				"data":    nil,
			})
		}

		if strings.TrimSpace(headerValue) != os.Getenv("SECRETKEY_HEADER") {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"message": "Invalid secret key",
				"status":  false,
				"data":    nil,
			})
		}

		return next(c)
	}
}
