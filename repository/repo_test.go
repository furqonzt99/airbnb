package repository

import (
	"testing"

	"github.com/furqonzt99/airbnb/config"
	"github.com/furqonzt99/airbnb/model"
	"github.com/furqonzt99/airbnb/repository/feature"
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
var featureRepo *feature.FeatureRepository
var houseRepo *house.HouseRepository

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
	featureRepo = feature.NewFeatureRepo(db)
	houseRepo = house.NewHouseRepo(db)

	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.House{})
	db.AutoMigrate(&model.Feature{})
	db.AutoMigrate(&model.HouseHasFeatures{})
	db.AutoMigrate(&model.Transaction{})
	db.AutoMigrate(&model.Rating{})

	seed.FeatureSeed(db)

	m.Run()
}

func TestUser(t *testing.T) {
	t.Run("Register User", func(t *testing.T) {
		var mockUser model.User
		mockUser.Email = "test@gmail.com"
		mockUser.Password = "test123"
		mockUser.Name = "tester"

		res, err := userRepo.Register(mockUser)
		assert.Nil(t, err)
		assert.Equal(t, mockUser.Name, res.Name)
		assert.Equal(t, 1, int(res.ID))
	})

	t.Run("Error Register User Duplicate Email", func(t *testing.T) {
		var mockUser model.User
		mockUser.Email = "test@gmail.com"
		mockUser.Password = "test123"
		mockUser.Name = "tester"

		_, err := userRepo.Register(mockUser)
		assert.NotNil(t, err)
	})

	t.Run("Login User", func(t *testing.T) {
		var mockUser model.User
		mockUser.Email = "test@gmail.com"

		res, err := userRepo.Login(mockUser.Email)
		assert.Nil(t, err)
		assert.Equal(t, res.Email, mockUser.Email)
	})

	t.Run("Error Login User No Email", func(t *testing.T) {
		var mockUser model.User
		mockUser.Email = "test123@gmail.com"

		_, err := userRepo.Login(mockUser.Email)
		assert.NotNil(t, err)
	})

	t.Run("Get User", func(t *testing.T) {
		userId := 1
		res, err := userRepo.Get(userId)
		assert.Nil(t, err)
		assert.Equal(t, res, res)
	})

	t.Run("Error Get User No ID", func(t *testing.T) {
		userId := 100
		_, err := userRepo.Get(userId)
		assert.NotNil(t, err)
	})

	t.Run("Update User ", func(t *testing.T) {
		var mockUser model.User
		mockUser.Email = "test@gmail.com"
		mockUser.Password = "test321"
		mockUser.Name = "tester2"

		userId := 1

		res, err := userRepo.Update(mockUser, userId)
		assert.Nil(t, err)
		assert.Equal(t, mockUser.Email, res.Email)
		assert.Equal(t, mockUser.Password, res.Password)
		assert.Equal(t, mockUser.Name, res.Name)
	})

	t.Run("Error Update User No ID", func(t *testing.T) {
		var mockUser model.User
		mockUser.Email = "test@gmail.com"

		userId := 100

		_, err := userRepo.Update(mockUser, userId)
		assert.NotNil(t, err)
	})

	t.Run("Delete User", func(t *testing.T) {
		userId := 1
		res, err := userRepo.Delete(userId)
		assert.Nil(t, err)
		assert.Equal(t, res, res)
	})

	t.Run("Error Delete User No ID", func(t *testing.T) {
		userId := 100
		_, err := userRepo.Delete(userId)
		assert.NotNil(t, err)
	})
}

func TestFeature(t *testing.T) {
	t.Run("Get All Features", func(t *testing.T) {
		res, err := featureRepo.GetAll()
		assert.Nil(t, err)
		assert.Equal(t, true, len(res) > 0)
	})
}

func TestHouse(t *testing.T) {
	t.Run("Create House", func(t *testing.T) {
		var mockHouse model.House
		mockHouse.UserID = 1
		mockHouse.Title = "rumah"
		mockHouse.Address = "jalan ujung"
		mockHouse.City = "indonesia"
		mockHouse.Price = 100000

		res, err := houseRepo.Create(mockHouse)
		assert.Nil(t, err)
		assert.Equal(t, 1, int(res.ID))
		assert.Equal(t, 1, int(res.UserID))
	})

	t.Run("Error Create House No Fields Inserted", func(t *testing.T) {
		var mockHouse model.House

		_, err := houseRepo.Create(mockHouse)
		assert.NotNil(t, err)
	})

	t.Run("Get All House", func(t *testing.T) {
		offset := 0
		pageSize := 10
		search := "rumah"

		res, err := houseRepo.GetAll(offset, pageSize, search)
		assert.Nil(t, err)
		assert.Equal(t, res[0].Title, "rumah")
		assert.Equal(t, res[0].Address, "jalan ujung")
	})

	t.Run("Get All My House", func(t *testing.T) {
		userId := 1
		res, err := houseRepo.GetAllMine(userId)
		assert.Nil(t, err)
		assert.Equal(t, res[0].Title, "rumah")
		assert.Equal(t, res[0].Address, "jalan ujung")
	})

	t.Run("Update House", func(t *testing.T) {
		var mockHouse model.House
		mockHouse.Title = "rumah2"
		mockHouse.Address = "jalan awal"
		mockHouse.City = "indonesia"
		mockHouse.Price = 100000

		res, err := houseRepo.Update(mockHouse, 1, 1)
		assert.Nil(t, err)
		assert.Equal(t, res.Title, "rumah2")
		assert.Equal(t, res.Address, "jalan awal")
	})

	t.Run("Delete House", func(t *testing.T) {
		houseId := 1
		userId := 1
		_, err := houseRepo.Delete(houseId, userId)
		assert.Nil(t, err)
	})

	t.Run("Error Delete House No ID Or UserID", func(t *testing.T) {
		houseId := 100
		userId := 100
		_, err := houseRepo.Delete(houseId, userId)
		assert.NotNil(t, err)
	})

	t.Run("Save House Has Feature ", func(t *testing.T) {
		var mockHouseFeature model.HouseHasFeatures
		mockHouseFeature.HouseID = 1
		mockHouseFeature.FeatureID = 1

		err := houseRepo.HouseHasFeature(mockHouseFeature)
		assert.Equal(t, err, nil)
	})

	t.Run("Error Save House Has Feature  No ID", func(t *testing.T) {
		var mockHouseFeature model.HouseHasFeatures

		err := houseRepo.HouseHasFeature(mockHouseFeature)
		assert.NotNil(t, err)
	})

	t.Run("Delete House Has Feature", func(t *testing.T) {
		houseId := 1

		err := houseRepo.HouseHasFeatureDelete(houseId)
		assert.Equal(t, err, nil)
	})
}
