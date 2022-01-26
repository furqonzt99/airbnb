package transaction

import (
	"time"

	"github.com/furqonzt99/airbnb/model"
	"gorm.io/gorm"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (tr *TransactionRepository) GetAll(userId int, status string) ([]model.Transaction, error) {
	var transactions []model.Transaction

	if err := tr.db.Preload("User").Preload("House").Where("status lIKE ?", "%"+status+"%").Find(&transactions, "user_id = ?", userId).Error; err != nil {
		return nil, err
	}

	return transactions, nil
}

func (tr *TransactionRepository) GetByTransactionId(userId, trxId int) ([]model.Transaction, error) {
	var transactions []model.Transaction

	if err := tr.db.Preload("User").Preload("House").Where("user_id = ?", userId).Find(&transactions, trxId).Error; err != nil {
		return nil, err
	}

	return transactions, nil
}

func (tr *TransactionRepository) CheckAvailability(houseId int, checkinDate, checkoutDate time.Time) (bool, error) {
	var transactions []model.Transaction

	if err := tr.db.Where("checkout_date <= ? AND checkin_date >= ?", checkinDate, checkoutDate).First(&transactions, "house_id = ?", houseId).Error; err != nil {
		return false, err
	}

	return true, nil
}

func (tr *TransactionRepository) Get(userId int) (model.Transaction, error) {
	var transaction model.Transaction

	if err := tr.db.Preload("User").Preload("House").Where("user_id = ?", userId).First(&transaction).Error; err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (tr *TransactionRepository) GetByInvoice(invId string) (model.Transaction, error) {
	var transaction model.Transaction

	if err := tr.db.Preload("User").Preload("House").Where("invoice_id = ?", invId).First(&transaction).Error; err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (tr *TransactionRepository) Create(transaction model.Transaction) (model.Transaction, error) {
	if err := tr.db.Create(&transaction).Error; err != nil {
		return transaction, err
	}

	var t model.Transaction

	if err := tr.db.Preload("User").Preload("House").First(&t, &transaction.ID).Error; err != nil {
		return transaction, err
	}

	return t, nil
}

func (tr *TransactionRepository) Update(invId string, transaction model.Transaction) (model.Transaction, error) {
	var t model.Transaction

	if err := tr.db.First(&t, "invoice_id = ?", invId).Error; err != nil {
		return t, err
	}

	tr.db.Model(&t).Updates(transaction)

	return t, nil
}