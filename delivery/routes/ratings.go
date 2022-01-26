package routes

import (
	"github.com/furqonzt99/airbnb/constant"
	"github.com/furqonzt99/airbnb/delivery/controllers/rating"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RegisterRatingPath(e *echo.Echo, RatingController *rating.RatingController) {
	
	e.POST("/ratings", RatingController.Create, middleware.JWT([]byte(constant.JWT_SECRET_KEY)))
	e.PUT("/ratings/:productId", RatingController.Update, middleware.JWT([]byte(constant.JWT_SECRET_KEY)))
	e.DELETE("/ratings/:productId", RatingController.Delete, middleware.JWT([]byte(constant.JWT_SECRET_KEY)))
}