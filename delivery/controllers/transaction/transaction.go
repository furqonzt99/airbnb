package transaction

import (
	"net/http"
	"strings"
	"time"

	"github.com/furqonzt99/airbnb/delivery/common"
	mw "github.com/furqonzt99/airbnb/delivery/middleware"
	"github.com/furqonzt99/airbnb/helper"
	"github.com/furqonzt99/airbnb/model"
	tr "github.com/furqonzt99/airbnb/repository/transaction"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type TransactionController struct {
	Repository tr.Transaction
}

func NewTransactionController(repo tr.Transaction) *TransactionController {
	return &TransactionController{Repository: repo}
}

func (tc TransactionController) Booking(c echo.Context) error {
	var transactionRequest TransactionRequest

	// bind request data
	if err := c.Bind(&transactionRequest); err != nil {
		return c.JSON(http.StatusBadRequest, common.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	if err := c.Validate(&transactionRequest); err != nil {
      return c.JSON(http.StatusBadRequest, common.ErrorResponse(http.StatusBadRequest, err.Error()))
    }

	user, _ := mw.ExtractTokenUser(c)
	invoiceId := strings.ToUpper(strings.ReplaceAll(uuid.New().String(), "-", ""))

	checkinDate, _ := time.Parse(time.RFC3339, transactionRequest.CheckinDate + "T12:00:00.000Z")
	checkoutDate, _ := time.Parse(time.RFC3339, transactionRequest.CheckoutDate + "T12:00:00.000Z")

	data := model.Transaction{
		UserID:        uint(user.UserID),
		HouseID:       uint(transactionRequest.HouseID),
		InvoiceID:     invoiceId,
		CheckinDate:   checkinDate,
		CheckoutDate:  checkoutDate,
	}

	transactionData, err := tc.Repository.Create(data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	transactionPayment, err := helper.CreateInvoice(transactionData, user.Email) 
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	response := TransactionResponse{
		ID:            int(transactionData.ID),
		UserID:        int(transactionData.UserID),
		HouseID:       transactionRequest.HouseID,
		InvoiceID:     invoiceId,
		PaymentUrl:    transactionPayment.PaymentUrl,
		CheckinDate:   transactionRequest.CheckinDate,
		CheckoutDate:  transactionRequest.CheckoutDate,
		TotalPrice:    transactionPayment.TotalPrice,
		Status:        transactionPayment.Status,
	}

	return c.JSON(http.StatusOK, common.SuccessResponse(response))
}