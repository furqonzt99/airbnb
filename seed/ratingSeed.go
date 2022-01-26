package seed

import (
	"fmt"
	"math/rand"

	"github.com/furqonzt99/airbnb/model"
	"gorm.io/gorm"
)

func RatingSeed(db *gorm.DB) {
	for i := 1; i <= 45; i++ {
		houseId := uint(rand.Intn(15-1) + 1)
		userId := uint(rand.Intn(5-1) + 1)
		rating := model.Rating{
			HouseID: houseId,
			UserID:  userId,
			Rating:  rand.Intn(5-1) + 1,
			Comment: fmt.Sprintf("Rating Untuk House %v dari User %v", houseId, userId),
		}
		db.Create(&rating)
	}
}
