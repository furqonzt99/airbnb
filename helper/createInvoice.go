package helper

import (
	"fmt"
	"os"

	"github.com/furqonzt99/airbnb/model"
	"github.com/xendit/xendit-go"
	"github.com/xendit/xendit-go/invoice"
)

func CreateInvoice(transaction model.Transaction, email string) (model.Transaction, error) {
	xendit.Opt.SecretKey = os.Getenv("XENDIT_SECRET_KEY")

	totalNight := CountNight(transaction.CheckinDate, transaction.CheckoutDate)
	totalPrice := transaction.House.Price * float64(totalNight)

	items := []xendit.InvoiceItem{
		{
			Name:     transaction.House.Title,
			Price:    transaction.House.Price,
			Quantity: totalNight,
		},
	}

	data := invoice.CreateParams{
		ExternalID:      transaction.InvoiceID,
		Amount:          totalPrice,
		Description:     "Invoice " + transaction.InvoiceID + " for " + email,
		PayerEmail:      email,
		Items:           items,
	}

	resp, err := invoice.Create(&data)
	if err != nil {
		fmt.Println(data.Amount)
		return transaction, err
	}

	transactionSuccess := model.Transaction{
		UserID:        transaction.UserID,
		HouseID:       transaction.HouseID,
		InvoiceID:     transaction.InvoiceID,
		PaymentUrl:    resp.InvoiceURL,
		CheckinDate:   transaction.CheckinDate,
		CheckoutDate:  transaction.CheckoutDate,
		TotalPrice:    resp.Amount,
		Status:        resp.Status,
	}

	return transactionSuccess, nil
}