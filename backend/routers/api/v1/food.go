package v1

import (
	"context"
	"fmt"
	"image"
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
	fmt.Println("moke", apiKey)

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
	res, err := client.NearbySearch(context.Background(), &maps.NearbySearchRequest{
		Location: &maps.LatLng{
			Lat: req.Latitude,
			Lng: req.Longitude,
		},
		Radius:   nearbySearchRadius,
		Keyword:  req.Keyword,
		Language: "ja",
	})
	return c.JSON(200, res)
}

type GetFoodDetailResponse struct {
	Detail maps.PlaceDetailsResult `json:"detail"`
	Photo  image.Image             `json:"photo"`
}

func GetFoodDetail(c echo.Context) error {
	placeID := c.Param("id")
	photoReference := c.Param("ref")
	apiKey := os.Getenv("GOOGLE_API_KEY")
	if apiKey == "" {
		return echo.NewHTTPError(http.StatusInternalServerError, "FOOD_API_KEY is not set")
	}

	client, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	detail, err := client.PlaceDetails(context.Background(), &maps.PlaceDetailsRequest{
		PlaceID:  placeID,
		Language: "ja",
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	ppres, err := client.PlacePhoto(context.Background(), &maps.PlacePhotoRequest{
		PhotoReference: photoReference,
		MaxHeight:      300,
		MaxWidth:       400,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	photo, err := ppres.Image()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(200, &GetFoodDetailResponse{
		Detail: detail,
		Photo:  photo,
	})
}
