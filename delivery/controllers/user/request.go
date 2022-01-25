package user

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type RegisterUserRequestFormat struct {
	Name     string `json:"name" form:"name" validate:"required"`
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required,min=8"`
}

type PutUserRequestFormat struct {
	Email    string `json:"email" form:"email" validate:"required,email"`
	Name     string `json:"name" form:"name" validate:"required"`
	Password string `json:"password" form:"password" validate:"required,min=8"`
}

type UserValidator struct {
	Validator *validator.Validate
}

func (cv *UserValidator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}
