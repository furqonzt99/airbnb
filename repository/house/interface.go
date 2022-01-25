package house

import "github.com/furqonzt99/airbnb/model"

type HouseInterface interface {
	Create(newHouse model.House) (model.House, error)
	HouseHasFeature(houseHasFeature model.HouseHasFeatures) error
	GetAll(offset, pageSize int, search string) ([]model.House, error)
	Get(houseId int) (model.House, error)
	Update(newHouse model.House, houseId int) (model.House, error)
	Delete(houseId int) (model.House, error)
}
