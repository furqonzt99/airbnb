package util

import (
	"github.com/furqonzt99/airbnb/config"
	"github.com/furqonzt99/airbnb/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB(config *config.AppConfig) *gorm.DB {

	conn := config.Database.Username + ":" + config.Database.Password + "@tcp(" + config.Database.Host + ":" + config.Database.Port + ")/" + config.Database.Name + "?parseTime=true&loc=Asia%2FJakarta&charset=utf8mb4&collation=utf8mb4_unicode_ci"

	db, err := gorm.Open(mysql.Open(conn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	return db
}

func InitialMigrate(db *gorm.DB) {
	if config.Mode == "development" {
		// db.Migrator().DropTable(&model.User{})
		// db.Migrator().DropTable(&model.HouseHasFeatures{})
		// db.Migrator().DropTable(&model.House{})
		// db.Migrator().DropTable(&model.Feature{})
		// db.Migrator().DropTable(&model.Rating{})
		// db.Migrator().DropTable(&model.Transaction{})

		db.AutoMigrate(&model.User{})
		db.AutoMigrate(&model.House{})
		db.AutoMigrate(&model.Feature{})
		db.AutoMigrate(&model.Rating{})
		db.AutoMigrate(&model.Transaction{})

		// seed.FeatureSeed(db)
		// seed.UserSeed(db)
		// seed.HouseSeed(db)
	} else {
		db.AutoMigrate(&model.User{})
		db.AutoMigrate(&model.House{})
		db.AutoMigrate(&model.Feature{})
		db.AutoMigrate(&model.Rating{})
		db.AutoMigrate(&model.Transaction{})

	}

}
