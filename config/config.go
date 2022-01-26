package config

import (
	"os"
	"sync"

	"github.com/furqonzt99/airbnb/constant"

	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
)

type AppConfig struct {
	Port     string
	Database struct {
		Driver   string
		Name     string
		Host     string
		Port     string
		Username string
		Password string
	}
}

var lock = &sync.Mutex{}
var appConfig *AppConfig
var Mode string

func GetConfig() *AppConfig {
	lock.Lock()
	defer lock.Unlock()

	if appConfig == nil {
		appConfig = initConfig()
	}

	return appConfig
}

func initConfig() *AppConfig {

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var defaultConfig AppConfig
	defaultConfig.Port = os.Getenv("APP_PORT")
	defaultConfig.Database.Driver = os.Getenv("DB_DRIVER")
	defaultConfig.Database.Name = os.Getenv("DB_NAME")
	defaultConfig.Database.Host = os.Getenv("DB_HOST")
	defaultConfig.Database.Port = os.Getenv("DB_PORT")
	defaultConfig.Database.Username = os.Getenv("DB_USERNAME")
	defaultConfig.Database.Password = os.Getenv("DB_PASSWORD")

	constant.JWT_SECRET_KEY = os.Getenv("JWT_SECRET_KEY")
	constant.XENDIT_CALLBACK_TOKEN = os.Getenv("XENDIT_CALLBACK_TOKEN")
	constant.PAYMENT_DURATION = os.Getenv("PAYMENT_DURATION")

	Mode = os.Getenv("MODE")

	return &defaultConfig
}