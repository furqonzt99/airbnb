package model

import "gorm.io/gorm"

type House struct {
	gorm.Model
	UserID uint
	Title string
	Address string
	City string
	Price string
	User User
	Features []Feature `gorm:"many2many:house_has_features;"`
}