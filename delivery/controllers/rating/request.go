package rating

import (
	"net/http"

	"github.com/furqonzt99/airbnb/delivery/common"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type PostRatingRequest struct {
	HouseID int `json:"house_id" validate:"required"`
	Rating int `json:"rating" validate:"required,max=5,min=1"`
	Comment string `json:"comment"`
}

type UpdateRatingRequest struct {
	Rating int `json:"rating" validate:"required,max=5,min=1"`
	Comment string `json:"comment"`
}

type TransactionValidator struct {
	Validator *validator.Validate
}

func (tv *TransactionValidator) Validate(i interface{}) error {
	if err := tv.Validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(http.StatusBadRequest, common.NewBadRequestResponse())
	}
	return nil
}