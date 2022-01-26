package house

import "github.com/furqonzt99/airbnb/model"

type HouseInterface interface {
	Create(newHouse model.House) (model.House, error)
	GetAll(offset, pageSize int, search string) ([]model.House, error)
	GetAllMine(userId int) ([]model.House, error)
	Get(houseId int) (model.House, error)
	Update(newHouse model.House, houseId, userId int) (model.House, error)
	Delete(houseId, userId int) (model.House, error)
	HouseHasFeature(houseHasFeature model.HouseHasFeatures) error
	HouseHasFeatureDelete(houseId int) error
}
