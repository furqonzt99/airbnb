package transaction

import "github.com/furqonzt99/airbnb/model"

type Transaction interface {
	GetAll(userId int) ([]model.Transaction, error)
	GetByStatus(userId int, status string) ([]model.Transaction, error)
	Get(userId int) (model.Transaction, error)
	GetByInvoice(userId int, invId string) (model.Transaction, error)
	
	// CheckAvailability(userId int) (model.Transaction, error)
	
	Create(model.Transaction) (model.Transaction, error)

	Update(userId int, trxId int, transaction model.Transaction) (model.Transaction, error)
}