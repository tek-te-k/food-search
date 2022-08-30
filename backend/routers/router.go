package routers

import (
	"github.com/labstack/echo/v4"
)

func SetupRouter(app *echo.Echo) {
	app.GET("/", Index)
}

func Index(c echo.Context) error {
	return c.String(200, "Hello, World!")
}
