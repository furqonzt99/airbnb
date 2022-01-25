package routes

import (
	"github.com/furqonzt99/airbnb/constant"
	"github.com/furqonzt99/airbnb/delivery/controllers/user"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RegisterUserPath(e *echo.Echo, userCtrl *user.UsersController) {

	e.POST("/register", userCtrl.RegisterController())
	e.POST("/login", userCtrl.LoginController())
	e.GET("/profile", userCtrl.GetUserController(), middleware.JWT([]byte(constant.JWT_SECRET_KEY)))
	e.PUT("/users", userCtrl.UpdateUserController(), middleware.JWT([]byte(constant.JWT_SECRET_KEY)))
	e.DELETE("/users", userCtrl.DeleteUserCtrl(), middleware.JWT([]byte(constant.JWT_SECRET_KEY)))
}
