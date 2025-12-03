package main

import (
	"fmt"
	"log"
	_ "log"
	"net/http"
	"os"

	. "be-metalsteel/connection"
	"be-metalsteel/route"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

// @title The Master Steel API Documentation
// @version		1.0
// @description  Dokumentasi Api The Master Steel

// @BasePath		/
// @contact.name The Master Steel Dev
// @contact.url http://usepgnwan.my.id
// @contact.email usepgnwan76@gmail.com

func main() {
	ConnectDB()
	e := echo.New()
	e.Use(middleware.CORS())

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		log.Printf("err %v", err.Error())
		var code int
		var message string

		if httpErr, ok := err.(*echo.HTTPError); ok {
			code = httpErr.Code
			message = httpErr.Message.(string)
		} else {
			code = http.StatusInternalServerError
			message = err.Error()
		}

		if he, ok := err.(*echo.HTTPError); ok && he.Code == http.StatusMethodNotAllowed {
			code = http.StatusMethodNotAllowed
			message = "The HTTP method is not allowed for this route."
		}

		c.JSON(code, map[string]interface{}{
			"status":  false,
			"message": message,
		})
	}

	route.InitRouting(e)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3026"
	}

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}
