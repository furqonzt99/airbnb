package main

import (
	"github.com/furqonzt99/airbnb/config"
	"github.com/furqonzt99/airbnb/delivery/controllers/feature"
	"github.com/furqonzt99/airbnb/delivery/controllers/house"
	"github.com/furqonzt99/airbnb/delivery/controllers/user"
	mw "github.com/furqonzt99/airbnb/delivery/middleware"
	"github.com/furqonzt99/airbnb/delivery/routes"
	fr "github.com/furqonzt99/airbnb/repository/feature"
	hr "github.com/furqonzt99/airbnb/repository/house"
	ur "github.com/furqonzt99/airbnb/repository/user"
	"github.com/furqonzt99/airbnb/util"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	config := config.GetConfig()

	db := util.InitDB(config)

	util.InitialMigrate(db)

	userRepo := ur.NewUserRepo(db)
	houseRepo := hr.NewHouseRepo(db)
	featureRepo := fr.NewFeatureRepo(db)

	userCtrl := user.NewUsersControllers(userRepo)
	houseCtrl := house.NewHouseControllers(houseRepo)
	featureCtrl := feature.NewFeatureControllers(featureRepo)

	e := echo.New()
	mw.LogMiddleware(e)

	e.Pre(middleware.RemoveTrailingSlash())

	e.Validator = &user.UserValidator{Validator: validator.New()}
	e.Validator = &house.HouseValidator{Validator: validator.New()}

	routes.RegisterUserPath(e, userCtrl)
	routes.RegisterHousePath(e, houseCtrl)
	routes.RegisterFeaturePath(e, featureCtrl)

	e.Logger.Fatal(e.Start(":" + config.Port))
}
