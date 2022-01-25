package feature

import (
	"github.com/furqonzt99/airbnb/model"
	"gorm.io/gorm"
)

type FeatureRepository struct {
	db *gorm.DB
}

func NewFeatureRepo(db *gorm.DB) *FeatureRepository {
	return &FeatureRepository{db: db}
}

func (fr *FeatureRepository) GetAll() ([]model.Feature, error) {
	features := []model.Feature{}
	fr.db.Find(&features)

	return features, nil
}
