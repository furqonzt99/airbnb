package transaction

import (
	"net/http"

	"github.com/furqonzt99/airbnb/delivery/common"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type TransactionRequest struct {
	HouseID int	`json:"house_id" validate:"required"`
	CheckinDate string `json:"checkin_date" validate:"required"`
	CheckoutDate string	`json:"checkout_date" validate:"required"`
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