package main

import (
	"food-search-backend/routers"

	"github.com/labstack/echo/v4"
)

func main() {
	// TODO DB用のdocker-composeを作成する
	// database.Connect()
	app := echo.New()
	routers.SetupRouter(app)
	app.Logger.Fatal(app.Start(":8080"))
}
