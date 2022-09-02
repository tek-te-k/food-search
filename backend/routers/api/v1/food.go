package v1

import (
	"context"
	"fmt"
	"food-search-backend/routers/api"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"googlemaps.github.io/maps"
)

const imageWidth = 400

type SearchFoodsRequest struct {
	Latitude     float64 `json:"latitude" validate:"required"`
	Longitude    float64 `json:"longitude" validate:"required"`
	SearchRadius uint    `json:"search_radius" validate:"required"`
	Keyword      string  `json:"keyword" validate:"required"`
}

// SearchFoods search foods
// @Summary search foods
// @Description search foods
// @Tags food
// @Accept  json
// @Produce  json
// @Param search_foods_request body SearchFoodsRequest true "search foods request"
// @Success 200 {object} maps.NearbySearchResponse
// @Failure 400 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /api/v1/food [post]
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
	res, err := client.NearbySearch(context.Background(), &maps.NearbySearchRequest{
		Location: &maps.LatLng{
			Lat: req.Latitude,
			Lng: req.Longitude,
		},
		Radius:   req.SearchRadius,
		Keyword:  req.Keyword,
		Language: "ja",
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if len(res.Results) == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "not found")
	}
	return c.JSON(200, res)
}

type GetFoodDetailResponse struct {
	Detail   maps.PlaceDetailsResult `json:"detail"`
	PhotoUrl []string                `json:"photo_url"`
}

// GetFoodDetail get food detail
// @Summary get food detail
// @Description get food detail
// @Tags food
// @Accept  json
// @Produce  json
// @Param id path string true "place id"
// @Success 200 {object} GetFoodDetailResponse
// @Failure 400 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /api/v1/food/{id} [get]
func GetFoodDetail(c echo.Context) error {
	placeID := c.Param("id")
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
	// URL 内に API キーを含めないために，リダイレクト先の画像 URL を取得する
	redirectedPhotoUrl := []string{}
	for _, p := range detail.Photos {
		photoUrl := fmt.Sprintf(api.GoogleMapsPhotoUrl, imageWidth, p.PhotoReference, apiKey)
		client := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
		req, _ := http.NewRequest("GET", photoUrl, nil)
		res, err := client.Do(req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		redirectedPhotoUrl = append(redirectedPhotoUrl, res.Header.Get("Location"))
	}
	return c.JSON(200, &GetFoodDetailResponse{
		Detail:   detail,
		PhotoUrl: redirectedPhotoUrl,
	})
}
