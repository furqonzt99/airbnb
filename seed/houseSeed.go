package seed

import (
	"fmt"
	"math/rand"

	"github.com/furqonzt99/airbnb/model"
	"gorm.io/gorm"
)

func HouseSeed(db *gorm.DB) {
	for i := 1; i <= 15; i++ {
		house := model.House{
			UserID:   uint(rand.Intn(5-1) + 1),
			Title:    fmt.Sprint("House ", i),
			Address:  fmt.Sprint("Address ", i),
			City:     fmt.Sprint("City ", i),
			Price:    150000,
		}
		db.Create(&house)
	}

	for i := 1; i <= 45; i++ {
		houseFeatures := model.HouseHasFeatures{
			HouseID:   uint(rand.Intn(15-1) + 1),
			FeatureID: uint(rand.Intn(10-1) + 1),
		}
		db.Create(&houseFeatures)
	}
}
