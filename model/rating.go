package model

import "gorm.io/gorm"

type Rating struct {
	gorm.Model
	HouseID uint `gorm:"primaryKey"`
	UserID uint `gorm:"primaryKey"`
	Rating int
	Comment string
	House House
	User User
}