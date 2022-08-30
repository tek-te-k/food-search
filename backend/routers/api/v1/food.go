package v1

import "github.com/labstack/echo/v4"

func SearchFoods(c echo.Context) error {
	return c.String(200, "Hello, World!")
}
