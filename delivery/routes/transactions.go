package routes

import (
	"github.com/furqonzt99/airbnb/constant"
	"github.com/furqonzt99/airbnb/delivery/controllers/transaction"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RegisterTransactionPath(e *echo.Echo, TransactionController *transaction.TransactionController) {

	e.POST("/transactions/booking", TransactionController.Booking, middleware.JWT([]byte(constant.JWT_SECRET_KEY)))
	e.POST("/transactions/callback", TransactionController.Callback)
	e.GET("/transactions", TransactionController.GetAll, middleware.JWT([]byte(constant.JWT_SECRET_KEY)))
}
