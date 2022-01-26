package transaction

import (
	"github.com/furqonzt99/airbnb/model"
	"gorm.io/gorm"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (tr *TransactionRepository) GetAll(userId int) ([]model.Transaction, error) {
	var transactions []model.Transaction

	if err := tr.db.Preload("User").Preload("House").Where("user_id = ?", userId).Find(&transactions).Error; err != nil {
		return nil, err
	}

	return transactions, nil
}

func (tr *TransactionRepository) GetByStatus(userId int, status string) ([]model.Transaction, error) {
	var transactions []model.Transaction

	if err := tr.db.Preload("User").Preload("House").Where("user_id = ? AND status = ?", userId, status).Find(&transactions).Error; err != nil {
		return nil, err
	}

	return transactions, nil
}

func (tr *TransactionRepository) Get(userId int) (model.Transaction, error) {
	var transaction model.Transaction

	if err := tr.db.Preload("User").Preload("House").Where("user_id = ?", userId).First(&transaction).Error; err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (tr *TransactionRepository) GetByInvoice(userId int, invId string) (model.Transaction, error) {
	var transaction model.Transaction

	if err := tr.db.Preload("User").Preload("House").Where("user_id = ? AND invoice_id = ?", userId, invId).First(&transaction).Error; err != nil {
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

func (tr *TransactionRepository) Update(userId, transactionId int, transaction model.Transaction) (model.Transaction, error) {
	var t model.Transaction

	if err := tr.db.Where("user_id = ?", userId).First(&t, transactionId).Error; err != nil {
		return t, err
	}

	tr.db.Model(&t).Updates(transaction)

	return t, nil
}