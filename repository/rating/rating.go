package rating

import (
	"github.com/furqonzt99/airbnb/model"
	"gorm.io/gorm"
)

type RatingRepository struct {
	db *gorm.DB
}

func NewRatingRepository(db *gorm.DB) *RatingRepository {
	return &RatingRepository{db: db}
}

func (rr RatingRepository) Create(rating model.Rating) (model.Rating, error) {
	if err := rr.db.Create(&rating).Error; err != nil {
		return rating, err
	}

	return rating, nil
}

func (rr *RatingRepository) Update(rating model.Rating) (model.Rating, error) {
	var r model.Rating

	if err := rr.db.First(&r, "user_id = ? AND house_id", rating.HouseID, rating.UserID).Error; err != nil {
		return r, err
	}

	rr.db.Model(&r).Updates(rating)

	return r, nil
}

func (rr *RatingRepository) Delete(userId, houseId int) (model.Rating, error) {
	rating := model.Rating{}
	
	if err := rr.db.First(&rating, "user_id = ? AND house_id = ?", userId, houseId).Error; err != nil {
		return rating, err
	}

	rr.db.Delete(&rating)

	return rating, nil
}