package model

type Rating struct {
	HouseID uint `gorm:"primaryKey"`
	UserID uint `gorm:"primaryKey"`
	Rating int
	Comment string
}