package routers

import (
	v1 "food-search-backend/routers/api/v1"

	"github.com/labstack/echo/v4"
)

func SetupRouter(app *echo.Echo) {
	apiv1 := app.Group("/api/v1")
	apiv1.POST("/search-foods", v1.SearchFoods)
	apiv1.GET("/food-detail/:id/:ref", v1.GetFoodDetail)
}
