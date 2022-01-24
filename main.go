package main

import (
	"github.com/furqonzt99/airbnb/config"
	"github.com/furqonzt99/airbnb/util"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main()  {
	config := config.GetConfig()

	db := util.InitDB(config)

	util.InitialMigrate(db)

	e := echo.New()

	e.Pre(middleware.RemoveTrailingSlash())

	e.Logger.Fatal(e.Start(":" + config.Port))
}