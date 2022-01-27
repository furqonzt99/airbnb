package model

import (
	"gorm.io/gorm"
)

type House struct {
	gorm.Model
	UserID   uint    `gorm:"NOT NULL"`
	Title    string  `gorm:"NOT NULL"`
	Address  string  `gorm:"NOT NULL"`
	City     string  `gorm:"NOT NULL"`
	Price    float64 `gorm:"NOT NULL"`
	Status   string  `gorm:"NOT NULL;default:open"`
	User     User
	Features []Feature `gorm:"many2many:house_has_features;"`
	Ratings  []Rating
}

type HouseHasFeatures struct {
	HouseID   uint `gorm:"primaryKey"`
	FeatureID uint `gorm:"primaryKey"`
}

func (HouseHasFeatures) BeforeCreate(db *gorm.DB) error {
	err := db.SetupJoinTable(&House{}, "Features", &HouseHasFeatures{})
	if err != nil {
		return err
	}
	return nil
}
