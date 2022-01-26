package transaction

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/furqonzt99/airbnb/constant"
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

	checkinDate, _ := time.Parse(time.RFC3339, transactionRequest.CheckinDate + "T00:00:00.000Z")
	checkoutDate, _ := time.Parse(time.RFC3339, transactionRequest.CheckoutDate + "T00:00:00.000Z")

	// if checkout date < checkin date
	if !checkoutDate.After(checkinDate) {
		return c.JSON(http.StatusBadRequest, common.ErrorResponse(http.StatusBadRequest, "Checkout date must after checkin date!"))
	}

	today := time.Now()

	// checkin or checkout cant past time/date 
	if !checkinDate.After(today) && !checkoutDate.After(today) {
		return c.JSON(http.StatusBadRequest, common.ErrorResponse(http.StatusBadRequest, "Checkin date or checkout date cant past date!"))
	}
	// check availability
	isAvailable, _ := tc.Repository.IsHouseAvailable(transactionRequest.HouseID, checkinDate, checkoutDate)
	if !isAvailable {
		return c.JSON(http.StatusBadRequest, common.ErrorResponse(http.StatusBadRequest, "House already booked at the date, please choose another date!"))
	}

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

	// add payment url to db
	updateData := model.Transaction{}
	updateData.PaymentUrl = transactionPayment.PaymentUrl
	updateData.TotalPrice = transactionPayment.TotalPrice
	tc.Repository.Update(invoiceId, updateData)

	// reformat response
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

func (tc TransactionController) Callback(c echo.Context) error {

	req := c.Request()
	headers := req.Header

	xCallbackToken := headers.Get("X-Callback-Token")

	if xCallbackToken != constant.XENDIT_CALLBACK_TOKEN {
		return c.JSON(http.StatusNotAcceptable, common.NewStatusNotAcceptable())
	}

	var callbackRequest common.CallbackRequest
	if err := c.Bind(&callbackRequest); err != nil {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
	}

	_, err := tc.Repository.GetByInvoice(callbackRequest.ExternalID) 
	if err != nil {
		return c.JSON(http.StatusNotFound, common.NewNotFoundResponse())
	}

	var data model.Transaction
	data.PaidAt, _ = time.Parse(time.RFC3339, callbackRequest.PaidAt)
	data.PaymentMethod = callbackRequest.PaymentMethod
	data.BankID = callbackRequest.BankID
	data.Status = callbackRequest.Status

	_, err = tc.Repository.Update(callbackRequest.ExternalID, data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
	}

	return c.JSON(http.StatusOK, common.NewSuccessOperationResponse())
}

func (tc TransactionController) GetAll(c echo.Context) error {

	user, err := mw.ExtractTokenUser(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
	}

	status := c.QueryParam("status")

	transactions, err := tc.Repository.GetAll(user.UserID, status)

	if err != nil {
		return c.JSON(http.StatusNotFound, common.NewNotFoundResponse())
	}

	transactionDatas := []TransactionResponse{}

	for _, td := range transactions {

		transactionDatas = append(transactionDatas, TransactionResponse{
			ID:            int(td.ID),
			UserID:        int(td.UserID),
			HouseID:       int(td.HouseID),
			InvoiceID:     td.InvoiceID,
			PaymentUrl:    td.PaymentUrl,
			BankID:        td.BankID,
			PaymentMethod: td.PaymentMethod,
			PaidAt:        fmt.Sprint(td.PaidAt),
			CheckinDate:   fmt.Sprint(td.CheckinDate),
			CheckoutDate:  fmt.Sprint(td.CheckoutDate),
			TotalPrice:    td.TotalPrice,
			Status:        td.Status,
		})
	}

	return c.JSON(http.StatusOK, common.SuccessResponse(transactionDatas))
}