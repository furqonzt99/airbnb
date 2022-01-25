package seed

import (
	"fmt"

	"github.com/furqonzt99/airbnb/model"
	"gorm.io/gorm"
)

func FeatureSeed(db *gorm.DB) {
	for i := 1; i <= 10; i++ {
		feature := model.Feature{
			Name: "Feature " + fmt.Sprint(i),
		}
		db.Create(&feature)
	}
}
