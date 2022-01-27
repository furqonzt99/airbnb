package transaction

import (
	"testing"
	"time"

	"github.com/furqonzt99/airbnb/config"
	"github.com/furqonzt99/airbnb/model"
	"github.com/furqonzt99/airbnb/repository/house"
	"github.com/furqonzt99/airbnb/repository/user"
	"github.com/furqonzt99/airbnb/seed"
	"github.com/furqonzt99/airbnb/util"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var configTest *config.AppConfig
var db *gorm.DB
var userRepo *user.UserRepository
var houseRepo *house.HouseRepository
var transactionRepo *TransactionRepository

var checkinDate time.Time
var checkoutDate time.Time

var rescheduleCheckinDate time.Time
var rescheduleCheckoutDate time.Time

func TestMain(m *testing.M) {
	configTest = config.GetConfig()
	db = util.InitDB(configTest)

	db.Migrator().DropTable(&model.User{})
	db.Migrator().DropTable(&model.House{})
	db.Migrator().DropTable(&model.Feature{})
	db.Migrator().DropTable(&model.HouseHasFeatures{})
	db.Migrator().DropTable(&model.Transaction{})
	db.Migrator().DropTable(&model.Rating{})

	userRepo = user.NewUserRepo(db)
	houseRepo = house.NewHouseRepo(db)
	transactionRepo = NewTransactionRepository(db)

	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.House{})
	db.AutoMigrate(&model.Feature{})
	db.AutoMigrate(&model.HouseHasFeatures{})
	db.AutoMigrate(&model.Transaction{})
	db.AutoMigrate(&model.Rating{})

	seed.UserSeed(db)
	seed.FeatureSeed(db)
	seed.HouseSeed(db)
	seed.RatingSeed(db)

	checkinDate = time.Now()
	checkoutDate = time.Now().AddDate(0, 0, 2)

	rescheduleCheckinDate = time.Now().AddDate(0, 0, 4)
	rescheduleCheckoutDate = time.Now().AddDate(0, 0, 6)

	m.Run()
}

func TestBooking(t *testing.T)  {
	
	t.Run("Success Booking 1", func(t *testing.T) {
		mockTransaction := model.Transaction{
			Model:          gorm.Model{},
			UserID:         1,
			HouseID:        4,
			HostID:         2,
			InvoiceID:      "US89IYSD9DAHA",
			PaymentUrl:     "url",
			PaymentChannel: "",
			PaymentMethod:  "",
			PaidAt:         time.Now(),
			CheckinDate:    checkinDate,
			CheckoutDate:   checkoutDate,
			TotalPrice:     300000,
			Status: "PAID",
		}

		res, err := transactionRepo.Create(mockTransaction)
		assert.Nil(t, err)
		assert.Equal(t, 1, int(res.UserID))
		assert.Equal(t, 4, int(res.HouseID))
		assert.Equal(t, 2, int(res.HostID))
		assert.Equal(t, "PAID", res.Status)
	})
	
	t.Run("Success Booking 2", func(t *testing.T) {
		mockTransaction := model.Transaction{
			Model:          gorm.Model{},
			UserID:         5,
			HouseID:        4,
			HostID:         2,
			InvoiceID:      "US89IYSD9DAHB",
			PaymentUrl:     "url",
			PaymentChannel: "",
			PaymentMethod:  "",
			PaidAt:         time.Now(),
			CheckinDate:    checkinDate.AddDate(0, 0, 2),
			CheckoutDate:   checkoutDate.AddDate(0, 0, 4),
			TotalPrice:     300000,
			Status: "PAID",
		}

		res, err := transactionRepo.Create(mockTransaction)
		assert.Nil(t, err)
		assert.Equal(t, 5, int(res.UserID))
		assert.Equal(t, 4, int(res.HouseID))
		assert.Equal(t, 2, int(res.HostID))
		assert.Equal(t, "PAID", res.Status)
	})
	
	t.Run("Failed Booking", func(t *testing.T) {
		mockTransaction := model.Transaction{
			Model:          gorm.Model{},
			InvoiceID:      "US89IYSD9DAHA",
			PaymentUrl:     "url",
			PaymentChannel: "",
			PaymentMethod:  "",
			PaidAt:         time.Now(),
			CheckinDate:    time.Now().AddDate(0, 0, 2),
			CheckoutDate:   time.Now(),
			TotalPrice:     300000,
		}

		_, err := transactionRepo.Create(mockTransaction)
		assert.NotNil(t, err)
	})
}

func TestGet(t *testing.T)  {
	
	t.Run("Success Get All", func(t *testing.T) {
		_, err := transactionRepo.GetAll(1, "")
		assert.Nil(t, err)
	})
	
	t.Run("Success Get All Host Trx", func(t *testing.T) {
		_, err := transactionRepo.GetAllHostTransaction(2, "")
		assert.Nil(t, err)
	})
	
	t.Run("Success Get By Trx ID", func(t *testing.T) {
		_, err := transactionRepo.GetByTransactionId(1, 1)
		assert.Nil(t, err)
	})

	t.Run("Failed Get By Trx ID", func(t *testing.T) {
		_, err := transactionRepo.GetByTransactionId(4, 4)
		assert.NotNil(t, err)
	})
	
	t.Run("Success Get By Inv ID", func(t *testing.T) {
		_, err := transactionRepo.GetByInvoice("US89IYSD9DAHA")
		assert.Nil(t, err)
	})

	t.Run("Failed Get By Inv ID", func(t *testing.T) {
		_, err := transactionRepo.GetByInvoice("US89IYSD9DAHC")
		assert.NotNil(t, err)
	})
	
	t.Run("Success Get", func(t *testing.T) {
		_, err := transactionRepo.Get(1)
		assert.Nil(t, err)
	})

	t.Run("Failed Get", func(t *testing.T) {
		_, err := transactionRepo.Get(7)
		assert.NotNil(t, err)
	})

}

func TestIsAvailable(t *testing.T)  {

	t.Run("Failed Is Available", func(t *testing.T) {
		res, err := transactionRepo.IsHouseAvailable(4, checkinDate.AddDate(0, 0, 10), checkoutDate.AddDate(0, 0, 12))
		assert.NotNil(t, err)
		assert.Equal(t, true, res)
	})
	
	t.Run("Success Is Available", func(t *testing.T) {
		res, err := transactionRepo.IsHouseAvailable(4, checkinDate.AddDate(0, 0, 1), checkoutDate)
		assert.Nil(t, err)
		assert.Equal(t, false, res)
	})

	t.Run("Failed Is Available Reschedule", func(t *testing.T) {
		res, err := transactionRepo.IsHouseAvailableReschedule(1, 4, checkinDate.AddDate(0, 0, 13), checkoutDate.AddDate(0, 0, 15))
		assert.NotNil(t, err)
		assert.Equal(t, true, res)
	})
	
	t.Run("Success Is Available Reschedule", func(t *testing.T) {
		res, err := transactionRepo.IsHouseAvailableReschedule(1, 4, time.Now().AddDate(0, 0, 1), time.Now().AddDate(0, 0, 2))
		assert.Nil(t, err)
		assert.Equal(t, false, res)
	})
	
}

func TestGetHostId(t *testing.T)  {
	
	t.Run("Success Get Host ID", func(t *testing.T) {
		_, err := transactionRepo.GetHostId(2)
		assert.Nil(t, err)
	})

	t.Run("Failed Get Host ID", func(t *testing.T) {
		_, err := transactionRepo.GetHostId(26)
		assert.NotNil(t, err)
	})
	
}

func TestUpdate(t *testing.T)  {
	
	t.Run("Success Update 1", func(t *testing.T) {
		mockTransaction := model.Transaction{
			Model:          gorm.Model{},
			UserID:         1,
			HouseID:        4,
			HostID:         2,
			InvoiceID:      "US89IYSD9DAHA",
			PaymentUrl:     "url",
			PaymentChannel: "",
			PaymentMethod:  "",
			PaidAt:         time.Now(),
			CheckinDate:    checkinDate,
			CheckoutDate:   checkoutDate,
			TotalPrice:     450000,
			Status: "PAID",
		}

		res, err := transactionRepo.Create(mockTransaction)
		assert.Nil(t, err)
		assert.Equal(t, 1, int(res.UserID))
		assert.Equal(t, 4, int(res.HouseID))
		assert.Equal(t, 2, int(res.HostID))
		assert.Equal(t, float64(450000), res.TotalPrice)
		assert.Equal(t, "PAID", res.Status)
	})
	
	t.Run("Success Booking 2", func(t *testing.T) {
		mockTransaction := model.Transaction{
			Model:          gorm.Model{},
			UserID:         5,
			HouseID:        4,
			HostID:         2,
			InvoiceID:      "US89IYSD9DAHB",
			PaymentUrl:     "urlnih",
			PaymentChannel: "",
			PaymentMethod:  "",
			PaidAt:         time.Now(),
			CheckinDate:    checkinDate.AddDate(0, 0, 2),
			CheckoutDate:   checkoutDate.AddDate(0, 0, 4),
			TotalPrice:     300000,
			Status: "PAID",
		}

		res, err := transactionRepo.Update("US89IYSD9DAHB", mockTransaction)
		assert.Nil(t, err)
		assert.Equal(t, 5, int(res.UserID))
		assert.Equal(t, 4, int(res.HouseID))
		assert.Equal(t, 2, int(res.HostID))
		assert.Equal(t, "urlnih", res.PaymentUrl)
		assert.Equal(t, "PAID", res.Status)
	})
	
	t.Run("Failed Booking", func(t *testing.T) {
		mockTransaction := model.Transaction{
			Model:          gorm.Model{},
			InvoiceID:      "US89IYSD9DAHA",
			PaymentUrl:     "url",
			PaymentChannel: "",
			PaymentMethod:  "",
			PaidAt:         time.Now(),
			CheckinDate:    time.Now().AddDate(0, 0, 2),
			CheckoutDate:   time.Now(),
			TotalPrice:     300000,
		}

		_, err := transactionRepo.Update("US89IYSD9DAHV", mockTransaction)
		assert.NotNil(t, err)
	})
}