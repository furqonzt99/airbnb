package transaction

import (
	"time"

	"github.com/furqonzt99/airbnb/model"
)

type Transaction interface {
	GetAll(userId int, status string) ([]model.Transaction, error)
	GetAllHostTransaction(hostId int, status string) ([]model.Transaction, error)
	Get(userId int) (model.Transaction, error)
	GetByInvoice(invId string) (model.Transaction, error)
	GetByTransactionId(userId, trxId int) (model.Transaction, error)
	
	GetHostId(houseId int) (int, error)
	
	IsHouseAvailable(houseId int, checkinDate, checkoutDate time.Time) (bool, error)
	IsHouseAvailableReschedule(trxId, houseId int, checkinDate, checkoutDate time.Time) (bool, error)
	
	Create(model.Transaction) (model.Transaction, error)

	Update(invId string, transaction model.Transaction) (model.Transaction, error)
}