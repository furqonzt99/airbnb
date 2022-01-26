package routes

import (
	"github.com/furqonzt99/airbnb/delivery/controllers/feature"
	"github.com/labstack/echo/v4"
)

func RegisterFeaturePath(e *echo.Echo, featureCtrl *feature.FeatureController) {

	e.GET("/features", featureCtrl.GetAllFeatureController())
}
