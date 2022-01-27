package transaction

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/furqonzt99/airbnb/constant"
	"github.com/furqonzt99/airbnb/delivery/common"
	"github.com/furqonzt99/airbnb/delivery/controllers/user"
	"github.com/furqonzt99/airbnb/model"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var jwtToken string

func TestMain(m *testing.M)  {

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	constant.JWT_SECRET_KEY = os.Getenv("JWT_SECRET_KEY")
	constant.XENDIT_CALLBACK_TOKEN = os.Getenv("XENDIT_CALLBACK_TOKEN")

	e := echo.New()
	e.Validator = &user.UserValidator{Validator: validator.New()}

	requestBody, _ := json.Marshal(map[string]string{
		"email":    "test@gmail.com",
		"password": "test1234",
	})

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
	res := httptest.NewRecorder()

	req.Header.Set("Content-Type", "application/json")
	context := e.NewContext(req, res)
	context.SetPath("/login")

	userController := user.NewUsersControllers(mockUserRepository{})
	userController.LoginController()(context)

	response := common.ResponseSuccess{}
	json.Unmarshal([]byte(res.Body.Bytes()), &response)

	jwtToken = response.Data.(string)

	m.Run()
}

func TestBooking(t *testing.T)  {
	e := echo.New()
	e.Validator = &TransactionValidator{Validator: validator.New()}

	checkInDate := fmt.Sprint(time.Now())[:10]
	checkoutDate := fmt.Sprint(time.Now().AddDate(0, 0, 2))[:10]
	
	falseAfterCheckoutDate := fmt.Sprint(time.Now().AddDate(0, 0, 5))[:10]
	falseBeforeNowCheckinDate := fmt.Sprint(time.Now().AddDate(0, 0, -5))[:10]

	t.Run("Transaction Booking Success", func(t *testing.T) {

		reqBody, _ := json.Marshal(TransactionRequest{
			HouseID:      1,
			CheckinDate:  checkInDate,
			CheckoutDate: checkoutDate,
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
		
		context := e.NewContext(req, res)
		context.SetPath("/transactions/booking")

		transactionController := NewTransactionController(mockTransactionRepository{})
		if err := middleware.JWT([]byte(constant.JWT_SECRET_KEY))(transactionController.Booking)(context); err != nil {
			log.Fatal(err)
			return
		}

		response := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, http.StatusOK, response.Code)
	})
	
	t.Run("Transaction Booking Fail Validator", func(t *testing.T) {

		reqBody, _ := json.Marshal(TransactionRequest{})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
		
		context := e.NewContext(req, res)
		context.SetPath("/transactions/booking")

		transactionController := NewTransactionController(mockTransactionRepository{})
		if err := middleware.JWT([]byte(constant.JWT_SECRET_KEY))(transactionController.Booking)(context); err != nil {
			log.Fatal(err)
			return
		}

		response := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})
	
	t.Run("Transaction Booking Fail checkin date > checkout date", func(t *testing.T) {

		reqBody, _ := json.Marshal(TransactionRequest{
			HouseID:      1,
			CheckinDate:  falseAfterCheckoutDate,
			CheckoutDate: checkoutDate,
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
		
		context := e.NewContext(req, res)
		context.SetPath("/transactions/booking")

		transactionController := NewTransactionController(mockTransactionRepository{})
		if err := middleware.JWT([]byte(constant.JWT_SECRET_KEY))(transactionController.Booking)(context); err != nil {
			log.Fatal(err)
			return
		}

		response := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})
	
	t.Run("Transaction Booking Fail checkin date < now", func(t *testing.T) {

		reqBody, _ := json.Marshal(TransactionRequest{
			HouseID:      1,
			CheckinDate:  falseBeforeNowCheckinDate,
			CheckoutDate: checkoutDate,
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
		
		context := e.NewContext(req, res)
		context.SetPath("/transactions/booking")

		transactionController := NewTransactionController(mockTransactionRepository{})
		if err := middleware.JWT([]byte(constant.JWT_SECRET_KEY))(transactionController.Booking)(context); err != nil {
			log.Fatal(err)
			return
		}

		response := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})
	
	t.Run("Transaction Booking Fail Unavailable", func(t *testing.T) {

		reqBody, _ := json.Marshal(TransactionRequest{
			HouseID:      1,
			CheckinDate:  checkInDate,
			CheckoutDate: checkoutDate,
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
		
		context := e.NewContext(req, res)
		context.SetPath("/transactions/booking")

		transactionController := NewTransactionController(mockFalseTransactionRepository{})
		if err := middleware.JWT([]byte(constant.JWT_SECRET_KEY))(transactionController.Booking)(context); err != nil {
			log.Fatal(err)
			return
		}

		response := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, http.StatusNotFound, response.Code)
	})
}

func TestReschedule(t *testing.T)  {
	e := echo.New()
	e.Validator = &TransactionValidator{Validator: validator.New()}

	checkInDate := fmt.Sprint(time.Now().AddDate(0, 0, 3))[:10]
	
	falseAfterCheckoutDate := fmt.Sprint(time.Now().AddDate(0, 0, 5))[:10]
	falseBeforeNowCheckinDate := fmt.Sprint(time.Now().AddDate(0, 0, -5))[:10]

	t.Run("Transaction Reschedule Success", func(t *testing.T) {

		reqBody, _ := json.Marshal(RescheduleRequest{
			CheckinDate:  checkInDate,
		})

		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewBuffer(reqBody))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
		
		context := e.NewContext(req, res)
		context.SetPath("/transactions/reschedule")
		context.SetParamNames("id")
		context.SetParamValues("1")

		transactionController := NewTransactionController(mockTransactionRepository{})
		if err := middleware.JWT([]byte(constant.JWT_SECRET_KEY))(transactionController.Reschedule)(context); err != nil {
			log.Fatal(err)
			return
		}

		response := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, http.StatusOK, response.Code)
	})
	
	t.Run("Transaction Reschedule Fail Validator", func(t *testing.T) {

		reqBody, _ := json.Marshal(RescheduleRequest{})

		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewBuffer(reqBody))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
		
		context := e.NewContext(req, res)
		context.SetPath("/transactions/reschedule")
		context.SetParamNames("id")
		context.SetParamValues("1")

		transactionController := NewTransactionController(mockTransactionRepository{})
		if err := middleware.JWT([]byte(constant.JWT_SECRET_KEY))(transactionController.Reschedule)(context); err != nil {
			log.Fatal(err)
			return
		}

		response := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})
	
	t.Run("Transaction Booking Fail Parameter", func(t *testing.T) {

		reqBody, _ := json.Marshal(RescheduleRequest{
			CheckinDate:  falseAfterCheckoutDate,
		})

		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewBuffer(reqBody))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
		
		context := e.NewContext(req, res)
		context.SetPath("/transactions/reschedule")
		context.SetParamNames("id")
		context.SetParamValues("ada8")

		transactionController := NewTransactionController(mockTransactionRepository{})
		if err := middleware.JWT([]byte(constant.JWT_SECRET_KEY))(transactionController.Reschedule)(context); err != nil {
			log.Fatal(err)
			return
		}

		response := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})
	
	t.Run("Transaction Booking Fail check status booking", func(t *testing.T) {

		reqBody, _ := json.Marshal(RescheduleRequest{
			CheckinDate:  falseBeforeNowCheckinDate,
		})

		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewBuffer(reqBody))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
		
		context := e.NewContext(req, res)
		context.SetPath("/transactions/reschedule")
		context.SetParamNames("id")
		context.SetParamValues("1")

		transactionController := NewTransactionController(mockTransactionRepository{})
		if err := middleware.JWT([]byte(constant.JWT_SECRET_KEY))(transactionController.Reschedule)(context); err != nil {
			log.Fatal(err)
			return
		}

		response := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})
	
	t.Run("Transaction Reschedule Fail Prev Data Not Found", func(t *testing.T) {

		reqBody, _ := json.Marshal(RescheduleRequest{
			CheckinDate:  falseBeforeNowCheckinDate,
		})

		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewBuffer(reqBody))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))
		
		context := e.NewContext(req, res)
		context.SetPath("/transactions/reschedule")
		context.SetParamNames("id")
		context.SetParamValues("1")

		transactionController := NewTransactionController(mockFalseTransactionRepository{})
		if err := middleware.JWT([]byte(constant.JWT_SECRET_KEY))(transactionController.Reschedule)(context); err != nil {
			log.Fatal(err)
			return
		}

		response := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, http.StatusNotFound, response.Code)
	})
}

func TestGetAll(t *testing.T)  {
	e := echo.New()
	
	t.Run("Get All Success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/transactions")

		transactionController := NewTransactionController(mockTransactionRepository{})
		if err := middleware.JWT([]byte(constant.JWT_SECRET_KEY))(transactionController.GetAll)(context); err != nil {
			log.Fatal(err)
			return
		}

		response := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("Get All Failed", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/transactions")

		transactionController := NewTransactionController(mockFalseTransactionRepository{})
		middleware.JWT([]byte(constant.JWT_SECRET_KEY))(transactionController.GetAll)(context)
			

		response := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, http.StatusNotFound, response.Code)
	})
}

func TestGetAllHost(t *testing.T)  {
	e := echo.New()
	
	t.Run("Get All Host Success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/transactions/host")

		transactionController := NewTransactionController(mockTransactionRepository{})
		if err := middleware.JWT([]byte(constant.JWT_SECRET_KEY))(transactionController.GetAllHostTransaction)(context); err != nil {
			log.Fatal(err)
			return
		}

		response := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("Get All Host Failed", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/transactions/host")

		transactionController := NewTransactionController(mockFalseTransactionRepository{})
		middleware.JWT([]byte(constant.JWT_SECRET_KEY))(transactionController.GetAllHostTransaction)(context)
			

		response := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, http.StatusNotFound, response.Code)
	})
}

func TestGetByTrx(t *testing.T)  {
	e := echo.New()
	
	t.Run("Get By Trx Success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/transactions")
		context.SetParamNames("id")
		context.SetParamValues("1")

		transactionController := NewTransactionController(mockTransactionRepository{})
		if err := middleware.JWT([]byte(constant.JWT_SECRET_KEY))(transactionController.GetByTransaction)(context); err != nil {
			log.Fatal(err)
			return
		}

		response := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("Get By Trx Failed", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/transactions")
		context.SetParamNames("id")
		context.SetParamValues("1")

		transactionController := NewTransactionController(mockFalseTransactionRepository{})
		middleware.JWT([]byte(constant.JWT_SECRET_KEY))(transactionController.GetByTransaction)(context)

		response := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, http.StatusNotFound, response.Code)
	})
}

func TestCallback(t *testing.T)  {
	e := echo.New()

	t.Run("Callback Success", func(t *testing.T) {

		reqBody, _ := json.Marshal(common.CallbackRequest{
			ExternalID:     "JHAKHSHJSIWOAM",
			PaymentMethod:  "BANK TRANSFER",
			PaymentChannel: "BRI",
			PaidAt:         fmt.Sprint(time.Now()),
			Status:         "PAID",
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Callback-Token", constant.XENDIT_CALLBACK_TOKEN)
		
		res := httptest.NewRecorder()
		
		context := e.NewContext(req, res)
		context.SetPath("/transactions/callback")

		transactionController := NewTransactionController(mockTransactionRepository{})
		transactionController.Callback(context)

		response := common.DefaultResponse{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, http.StatusOK, response.Code)
	})
	
	t.Run("Callback Not found", func(t *testing.T) {

		reqBody, _ := json.Marshal(common.CallbackRequest{
			ExternalID:     "JHAKHSHJSIWOAMIASS",
			PaymentMethod:  "BANK TRANSFER",
			PaymentChannel: "BRI",
			PaidAt:         fmt.Sprint(time.Now()),
			Status:         "PAID",
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Callback-Token", constant.XENDIT_CALLBACK_TOKEN)
		
		res := httptest.NewRecorder()
		
		context := e.NewContext(req, res)
		context.SetPath("/transactions/callback")

		transactionController := NewTransactionController(mockFalseTransactionRepository{})
		transactionController.Callback(context)

		response := common.DefaultResponse{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, http.StatusNotFound, res.Code)
	})

	t.Run("Callback Not Acceptable", func(t *testing.T) {

		reqBody, _ := json.Marshal(common.CallbackRequest{
			ExternalID:     "JHAKHSHJSIWOAM",
			PaymentMethod:  "BANK TRANSFER",
			PaymentChannel: "BRI",
			PaidAt:         fmt.Sprint(time.Now()),
			Status:         "PAID",
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Callback-Token", constant.XENDIT_CALLBACK_TOKEN + "false")
		
		res := httptest.NewRecorder()
		
		context := e.NewContext(req, res)
		context.SetPath("/transactions/callback")

		transactionController := NewTransactionController(mockTransactionRepository{})
		transactionController.Callback(context)

		response := common.DefaultResponse{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, http.StatusNotAcceptable, res.Code)
	})
}

type mockUserRepository struct{}

func (m mockUserRepository) Register(newUser model.User) (model.User, error) {
	hash, _ := bcrypt.GenerateFromPassword([]byte(newUser.Password), 14)
	return model.User{Email: newUser.Email, Password: string(hash), Name: newUser.Name}, nil
}

func (m mockUserRepository) Login(email string) (model.User, error) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("test1234"), 14)
	return model.User{
		Model:    gorm.Model{
			ID: 1,
		},
		Name:     "tester",
		Email:    "test@gmail.com",
		Password: string(hash),
	}, nil
}

func (m mockUserRepository) Get(userid int) (model.User, error) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("test1234"), 14)
	return model.User{Email: "test@gmail.com", Password: string(hash), Name: "tester"}, nil
}

func (m mockUserRepository) Update(newUser model.User, userId int) (model.User, error) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("test4321"), 14)
	return model.User{Email: "test2@gmail.com", Password: string(hash), Name: "tester2"}, nil
}

func (m mockUserRepository) Delete(userId int) (model.User, error) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("test4321"), 14)
	return model.User{Email: "test2@gmail.com", Password: string(hash), Name: "tester2"}, nil
}

type mockTransactionRepository struct{}

func (tr mockTransactionRepository) GetAll(userId int, status string) ([]model.Transaction, error) {
	return []model.Transaction{{
		Model:          gorm.Model{
			ID:        1,
		},
		UserID:         1,
		HouseID:        3,
		HostID:         2,
		InvoiceID:      "JHAKHSHJSIWOAM",
		PaymentUrl:     "url",
		PaymentChannel: "",
		PaymentMethod:  "",
		PaidAt:         time.Time{},
		CheckinDate:    time.Now(),
		CheckoutDate:   time.Now().AddDate(0, 0, 2),
		TotalPrice:     300000,
		Status:         "PENDING",
	}}, nil
}

func (tr mockTransactionRepository) GetAllHostTransaction(hostId int, status string) ([]model.Transaction, error) {
	return []model.Transaction{{
		Model:          gorm.Model{
			ID:        1,
		},
		UserID:         1,
		HouseID:        3,
		HostID:         2,
		InvoiceID:      "JHAKHSHJSIWOAM",
		PaymentUrl:     "url",
		PaymentChannel: "",
		PaymentMethod:  "",
		PaidAt:         time.Time{},
		CheckinDate:    time.Now(),
		CheckoutDate:   time.Now().AddDate(0, 0, 2),
		TotalPrice:     300000,
		Status:         "PENDING",
	}}, nil
}

func (tr mockTransactionRepository) GetByTransactionId(userId, trxId int) (model.Transaction, error) {
	return model.Transaction{
		Model:          gorm.Model{
			ID:        1,
		},
		UserID:         1,
		HouseID:        3,
		HostID:         2,
		InvoiceID:      "JHAKHSHJSIWOAM",
		PaymentUrl:     "url",
		PaymentChannel: "",
		PaymentMethod:  "",
		PaidAt:         time.Time{},
		CheckinDate:    time.Now(),
		CheckoutDate:   time.Now().AddDate(0, 0, 2),
		TotalPrice:     300000,
		Status:         "PAID",
	}, nil
}

func (tr mockTransactionRepository) GetHostId(houseId int) (int, error) {
	return int(1), nil
}

func (tr mockTransactionRepository) IsHouseAvailable(houseId int, checkinDate, checkoutDate time.Time) (bool, error) {
	return true, nil
}

func (tr mockTransactionRepository) IsHouseAvailableReschedule(trxId, houseId int, checkinDate, checkoutDate time.Time) (bool, error) {
	return true, nil
}

func (tr mockTransactionRepository) Get(userId int) (model.Transaction, error) {
	return model.Transaction{
		Model:          gorm.Model{
			ID:        1,
		},
		UserID:         1,
		HouseID:        3,
		HostID:         2,
		InvoiceID:      "JHAKHSHJSIWOAM",
		PaymentUrl:     "url",
		PaymentChannel: "",
		PaymentMethod:  "",
		PaidAt:         time.Time{},
		CheckinDate:    time.Now(),
		CheckoutDate:   time.Now().AddDate(0, 0, 2),
		TotalPrice:     300000,
		Status:         "PENDING",
	}, nil
}

func (tr mockTransactionRepository) GetByInvoice(invId string) (model.Transaction, error) {
	return model.Transaction{
		Model:          gorm.Model{
			ID:        1,
		},
		UserID:         1,
		HouseID:        3,
		HostID:         2,
		InvoiceID:      "JHAKHSHJSIWOAM",
		PaymentUrl:     "url",
		PaymentChannel: "",
		PaymentMethod:  "",
		PaidAt:         time.Time{},
		CheckinDate:    time.Now(),
		CheckoutDate:   time.Now().AddDate(0, 0, 2),
		TotalPrice:     300000,
		Status:         "PENDING",
	}, nil
}

func (tr mockTransactionRepository) Create(transaction model.Transaction) (model.Transaction, error) {
	return model.Transaction{
		Model:          gorm.Model{ID: 1},
		UserID:         1,
		HouseID:        3,
		HostID:         2,
		InvoiceID:      "JHAKHSHJSIWOAM",
		PaymentUrl:     "url",
		PaymentChannel: "",
		PaymentMethod:  "",
		PaidAt:         time.Time{},
		CheckinDate:    time.Now(),
		CheckoutDate:   time.Now().AddDate(0, 0, 2),
		TotalPrice:     300000,
		Status:         "PENDING",
		User:           model.User{},
		House:          model.House{
			Model:    gorm.Model{},
			UserID:   0,
			Title:    "House 1",
			Address:  "",
			City:     "",
			Price:    150000,
			Status:   "",
			User:     model.User{},
			Features: []model.Feature{},
			Ratings:  []model.Rating{},
		},
	}, nil
}

func (tr mockTransactionRepository) Update(invId string, transaction model.Transaction) (model.Transaction, error) {
	return model.Transaction{
		Model:          gorm.Model{
			ID:        1,
		},
		UserID:         1,
		HouseID:        3,
		HostID:         2,
		InvoiceID:      "JHAKHSHJSIWOAM",
		PaymentUrl:     "url",
		PaymentChannel: "",
		PaymentMethod:  "",
		PaidAt:         time.Time{},
		CheckinDate:    time.Now(),
		CheckoutDate:   time.Now().AddDate(0, 0, 2),
		TotalPrice:     300000,
		Status:         "PENDING",
	}, nil
}

type mockFalseTransactionRepository struct{}

func (tr mockFalseTransactionRepository) GetAll(userId int, status string) ([]model.Transaction, error) {
	return []model.Transaction{{}}, errors.New("Error")
}

func (tr mockFalseTransactionRepository) GetAllHostTransaction(hostId int, status string) ([]model.Transaction, error) {
	return []model.Transaction{{}}, errors.New("Error")
}

func (tr mockFalseTransactionRepository) GetByTransactionId(userId, trxId int) (model.Transaction, error) {
	return model.Transaction{}, errors.New("Error")
}

func (tr mockFalseTransactionRepository) GetHostId(houseId int) (int, error) {
	return int(0), errors.New("Error")
}

func (tr mockFalseTransactionRepository) IsHouseAvailable(houseId int, checkinDate, checkoutDate time.Time) (bool, error) {
	return false, errors.New("Error")
}

func (tr mockFalseTransactionRepository) IsHouseAvailableReschedule(trxId, houseId int, checkinDate, checkoutDate time.Time) (bool, error) {
	return false, errors.New("Error")
}

func (tr mockFalseTransactionRepository) Get(userId int) (model.Transaction, error) {
	return model.Transaction{}, errors.New("Error")
}

func (tr mockFalseTransactionRepository) GetByInvoice(invId string) (model.Transaction, error) {
	return model.Transaction{}, errors.New("Error")
}

func (tr mockFalseTransactionRepository) Create(transaction model.Transaction) (model.Transaction, error) {
	return model.Transaction{}, errors.New("Error")
}

func (tr mockFalseTransactionRepository) Update(invId string, transaction model.Transaction) (model.Transaction, error) {
	return model.Transaction{}, errors.New("Error")
}