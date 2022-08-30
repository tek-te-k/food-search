package main

import (
	"food-search-backend/pkg/validate"
	"food-search-backend/routers"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func main() {
	// TODO DB用のdocker-composeを作成する
	// database.Connect()
	app := echo.New()
	app.Validator = &validate.Validator{Validator: validator.New()}
	routers.SetupRouter(app)
	app.Logger.Fatal(app.Start(":8080"))
}
