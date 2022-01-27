package house

import (
	"github.com/furqonzt99/airbnb/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type HouseRepository struct {
	db *gorm.DB
}

func NewHouseRepo(db *gorm.DB) *HouseRepository {
	return &HouseRepository{db: db}
}

func (hr *HouseRepository) Create(newHouse model.House) (model.House, error) {
	if err := hr.db.Save(&newHouse).Error; err != nil {
		return newHouse, err
	}

	return newHouse, nil
}

func (hr *HouseRepository) GetAll(offset, pageSize int, search, city string) ([]model.House, error) {
	houses := []model.House{}

	hr.db.Preload("Features").Preload("User").Preload("Ratings.User").Preload(clause.Associations).Offset(offset).Limit(pageSize).Where("title LIKE ?", "%"+search+"%").Where("city LIKE ?", "%"+city+"%").Find(&houses)

	return houses, nil
}

func (hr *HouseRepository) GetAllMine(userId int) ([]model.House, error) {
	houses := []model.House{}

	hr.db.Preload("Features").Preload("User").Preload("Ratings.User").Preload(clause.Associations).Where("user_id=?", userId).Find(&houses)

	return houses, nil
}

func (hr *HouseRepository) Get(houseId int) (model.House, error) {
	house := model.House{}
	if err := hr.db.Preload("Features").Preload("User").Preload("Ratings.User").Preload(clause.Associations).Where("id = ?", houseId).First(&house).Error; err != nil {
		return house, err
	}

	return house, nil
}

func (hr *HouseRepository) Update(newHouse model.House, houseId, userId int) (model.House, error) {
	house := model.House{}

	if err := hr.db.Where("id = ? AND user_id = ?", houseId, userId).First(&house).Error; err != nil {
		return house, err
	}

	hr.db.Model(&house).Updates(newHouse)

	return house, nil
}

func (hr *HouseRepository) Delete(houseId, userId int) (model.House, error) {
	house := model.House{}
	if err := hr.db.First(&house, "id = ? AND user_id = ?", houseId, userId).Error; err != nil {
		return house, err
	}
	hr.db.Delete(&house)
	return house, nil
}

func (hr *HouseRepository) HouseHasFeature(houseHasFeature model.HouseHasFeatures) error {
	if err := hr.db.Save(&houseHasFeature).Error; err != nil {
		return err
	}
	return nil
}

func (hr *HouseRepository) HouseHasFeatureDelete(houseId int) error {
	house := []model.HouseHasFeatures{}
	hr.db.Find(&house, "house_id = ?", houseId)
	hr.db.Delete(&house)
	return nil
}
