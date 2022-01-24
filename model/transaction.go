package model

import (
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	UserID uint
	HouseID uint
	InvoiceID string
	PaymentUrl string
	BankID string
	PaymentMethod string
	PaidAt time.Time
	CheckinDate *time.Time
	CheckoutDate *time.Time
	TotalPrice float64
	Status string
	User User
	House House
}