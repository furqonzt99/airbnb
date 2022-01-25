package house

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type CreateHouseRequestFormat struct {
	Title    string  `json:"title" form:"title" validate:"required"`
	Address  string  `json:"address" form:"address" validate:"required"`
	City     string  `json:"city" form:"city" validate:"required"`
	Price    float64 `json:"price" form:"price" validate:"required"`
	Features []int   `json:"features" form:"features" validate:"required"`
}

type PutHouseRequestFormat struct {
	Title    string  `json:"title" form:"title" validate:"required"`
	Address  string  `json:"address" form:"address" validate:"required"`
	City     string  `json:"city" form:"city" validate:"required"`
	Price    float64 `json:"price" form:"price" validate:"required"`
	Features []int   `json:"features" form:"features" validate:"required"`
}

type HouseValidator struct {
	Validator *validator.Validate
}

func (cv *HouseValidator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}
