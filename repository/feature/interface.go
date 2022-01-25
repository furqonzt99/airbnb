package feature

import "github.com/furqonzt99/airbnb/model"

type FeatureInterface interface {
	GetAll() ([]model.Feature, error)
}
