package rating

import "github.com/furqonzt99/airbnb/model"

type Rating interface {
	Create(model.Rating) (model.Rating, error)
	Update(model.Rating) (model.Rating, error)
	Delete(userId, houseId int) (model.Rating, error)
	IsCanGiveRating(userId, houseId int) (bool, error)
}