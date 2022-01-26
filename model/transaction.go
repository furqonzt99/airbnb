package model

import (
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	UserID uint `gorm:"not null"`
	HouseID uint `gorm:"not null"`
	InvoiceID string
	PaymentUrl string
	BankID string
	PaymentMethod string
	PaidAt time.Time `gorm:"default:null"`
	CheckinDate time.Time
	CheckoutDate time.Time
	TotalPrice float64
	Status string `gorm:"not null;default:PENDING"`
	User User
	House House
}