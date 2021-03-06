package transaction

import (
	"fmt"
	"net/http"
	"strconv"
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

	hostId, err := tc.Repository.GetHostId(transactionRequest.HouseID)
	if err != nil {
		return c.JSON(http.StatusNotFound, common.NewNotFoundResponse())
	}

	invoiceId := strings.ToUpper(strings.ReplaceAll(uuid.New().String(), "-", ""))

	checkinDate, _ := time.Parse(time.RFC3339, transactionRequest.CheckinDate + "T00:00:00.000Z")
	checkoutDate, _ := time.Parse(time.RFC3339, transactionRequest.CheckoutDate + "T00:00:00.000Z")

	// if checkout date < checkin date
	if !checkoutDate.After(checkinDate) {
		return c.JSON(http.StatusBadRequest, common.ErrorResponse(http.StatusBadRequest, "Checkout date must after checkin date!"))
	}

	yesterday := time.Now().AddDate(0, 0, -1)

	// checkin or checkout cant past time/date 
	if !checkinDate.After(yesterday) || !checkoutDate.After(yesterday) {
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
		HostID: uint(hostId),
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
		HostID: hostId,
		InvoiceID:     invoiceId,
		PaymentUrl:    transactionPayment.PaymentUrl,
		CheckinDate:   transactionRequest.CheckinDate,
		CheckoutDate:  transactionRequest.CheckoutDate,
		TotalPrice:    transactionPayment.TotalPrice,
		Status:        transactionPayment.Status,
	}

	return c.JSON(http.StatusOK, common.SuccessResponse(response))
}

func (tc TransactionController) Reschedule(c echo.Context) error {
	var rescheduleRequest RescheduleRequest

	// bind request data
	if err := c.Bind(&rescheduleRequest); err != nil {
		return c.JSON(http.StatusBadRequest, common.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	if err := c.Validate(&rescheduleRequest); err != nil {
      return c.JSON(http.StatusBadRequest, common.ErrorResponse(http.StatusBadRequest, err.Error()))
    }

	trxId, err := strconv.Atoi(c.Param("id")) 
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	user, _ := mw.ExtractTokenUser(c)

	prevData, err := tc.Repository.GetByTransactionId(user.UserID, trxId)
	if err != nil {
		return c.JSON(http.StatusNotFound, common.NewNotFoundResponse())
	}

	// decline reschedule if haven't paid yet
	const PAID_STATUS = "PAID"

	if prevData.Status != PAID_STATUS {
		return c.JSON(http.StatusNotAcceptable, common.NewStatusNotAcceptable())
	}

	// get prev data
	prevCheckinDate := prevData.CheckinDate
	prevCheckoutDate := prevData.CheckoutDate
	
	// count prev data night  
	countNight := helper.CountNight(prevCheckinDate, prevCheckoutDate)
	
	// assign reschedule date
	checkinDate, _ := time.Parse(time.RFC3339, rescheduleRequest.CheckinDate + "T00:00:00.000Z")
	checkoutDate := checkinDate.AddDate(0, 0, countNight)

	// if checkout date < checkin date
	if !checkoutDate.After(checkinDate) {
		return c.JSON(http.StatusBadRequest, common.ErrorResponse(http.StatusBadRequest, "Checkout date must after checkin date!"))
	}

	yesterday := time.Now().AddDate(0, 0, -1)

	// checkin or checkout cant past time/date 
	if !checkinDate.After(yesterday) || !checkoutDate.After(yesterday) {
		return c.JSON(http.StatusBadRequest, common.ErrorResponse(http.StatusBadRequest, "Checkin date or checkout date cant past date!"))
	}

	// check reschedule availability
	isAvailable, _ := tc.Repository.IsHouseAvailableReschedule(trxId, int(prevData.HouseID), checkinDate, checkoutDate)
	if !isAvailable {
		return c.JSON(http.StatusBadRequest, common.ErrorResponse(http.StatusBadRequest, "House already booked at the date, please choose another date!"))
	}

	data := model.Transaction{
		CheckinDate:   checkinDate,
		CheckoutDate:  checkoutDate,
	}

	_, err = tc.Repository.Update(prevData.InvoiceID, data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	return c.JSON(http.StatusOK, common.NewSuccessOperationResponse())
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
	data.PaymentChannel = callbackRequest.PaymentChannel
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
			PaymentChannel: td.PaymentChannel,
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

func (tc TransactionController) GetAllHostTransaction(c echo.Context) error {

	user, err := mw.ExtractTokenUser(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
	}

	status := c.QueryParam("status")

	transactions, err := tc.Repository.GetAllHostTransaction(user.UserID, status)

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
			PaymentChannel: td.PaymentChannel,
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

func (tc TransactionController) GetByTransaction(c echo.Context) error {

	user, err := mw.ExtractTokenUser(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
	}

	trxId, err := strconv.Atoi(c.Param("id")) 
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.NewBadRequestResponse())
	}

	transaction, err := tc.Repository.GetByTransactionId(user.UserID, trxId)
	if err != nil {
		return c.JSON(http.StatusNotFound, common.NewNotFoundResponse())
	}

	transactionData := TransactionResponse{
		ID:            trxId,
		UserID:        user.UserID,
		HouseID:       int(transaction.UserID),
		InvoiceID:     transaction.InvoiceID,
		PaymentUrl:    transaction.PaymentUrl,
		PaymentChannel: transaction.PaymentChannel,
		PaymentMethod: transaction.PaymentMethod,
		PaidAt:        fmt.Sprint(transaction.PaidAt),
		CheckinDate:   fmt.Sprint(transaction.CheckinDate),
		CheckoutDate:  fmt.Sprint(transaction.CheckoutDate),
		TotalPrice:    transaction.TotalPrice,
		Status:        transaction.Status,
	}

	return c.JSON(http.StatusOK, common.SuccessResponse(transactionData))
}