package model

import (
	"gorm.io/gorm"
)

type House struct {
	gorm.Model
	UserID   uint
	Title    string
	Address  string
	City     string
	Price    float64
	User     User
	Features []Feature `gorm:"many2many:house_has_features;"`
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
