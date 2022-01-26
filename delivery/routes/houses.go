package routes

import (
	"github.com/furqonzt99/airbnb/constant"
	"github.com/furqonzt99/airbnb/delivery/controllers/house"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RegisterHousePath(e *echo.Echo, houseCtrl *house.HouseController) {

	e.POST("/houses", houseCtrl.CreateHouseController(), middleware.JWT([]byte(constant.JWT_SECRET_KEY)))
	e.GET("/houses", houseCtrl.GetAllHouseController())
	e.GET("/myhouses", houseCtrl.GetMyHouseController(), middleware.JWT([]byte(constant.JWT_SECRET_KEY)))
	e.GET("/houses/:id", houseCtrl.GetHouseController())
	e.PUT("/houses/:id", houseCtrl.UpdateHouseController(), middleware.JWT([]byte(constant.JWT_SECRET_KEY)))
	e.DELETE("/houses/:id", houseCtrl.DeleteHouseCtrl(), middleware.JWT([]byte(constant.JWT_SECRET_KEY)))
}
