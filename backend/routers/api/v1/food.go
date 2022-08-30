package v1

import (
	"context"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"googlemaps.github.io/maps"
)

const nearbySearchRadius = 2000

type SearchFoodsRequest struct {
	Latitude  float64 `json:"latitude" validate:"required"`
	Longitude float64 `json:"longitude" validate:"required"`
	Keyword   string  `json:"keyword" validate:"required"`
}

func SearchFoods(c echo.Context) error {
	apiKey := os.Getenv("GOOGLE_API_KEY")
	if apiKey == "" {
		return echo.NewHTTPError(http.StatusInternalServerError, "FOOD_API_KEY is not set")
	}
	req := new(SearchFoodsRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	client, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	pr, err := client.NearbySearch(context.Background(), &maps.NearbySearchRequest{
		Location: &maps.LatLng{
			Lat: req.Latitude,
			Lng: req.Longitude,
		},
		Radius:  nearbySearchRadius,
		Keyword: req.Keyword,
	})
	return c.JSON(200, pr)
}
