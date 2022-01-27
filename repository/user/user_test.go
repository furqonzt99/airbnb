package user

import (
	"testing"

	"github.com/furqonzt99/airbnb/config"
	"github.com/furqonzt99/airbnb/model"
	"github.com/furqonzt99/airbnb/seed"
	"github.com/furqonzt99/airbnb/util"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var configTest *config.AppConfig
var db *gorm.DB
var userRepo *UserRepository

func TestRegisterUser(t *testing.T) {
	configTest = config.GetConfig()
	db = util.InitDB(configTest)

	db.Migrator().DropTable(&model.User{})
	db.Migrator().DropTable(&model.House{})
	db.Migrator().DropTable(&model.Feature{})
	db.Migrator().DropTable(&model.HouseHasFeatures{})
	db.Migrator().DropTable(&model.Transaction{})
	db.Migrator().DropTable(&model.Rating{})

	userRepo = NewUserRepo(db)

	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.House{})
	db.AutoMigrate(&model.Feature{})
	db.AutoMigrate(&model.HouseHasFeatures{})
	db.AutoMigrate(&model.Transaction{})
	db.AutoMigrate(&model.Rating{})

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
}

func TestLoginUser(t *testing.T) {
	configTest = config.GetConfig()
	db = util.InitDB(configTest)

	db.Migrator().DropTable(&model.User{})
	db.Migrator().DropTable(&model.House{})
	db.Migrator().DropTable(&model.Feature{})
	db.Migrator().DropTable(&model.HouseHasFeatures{})
	db.Migrator().DropTable(&model.Transaction{})
	db.Migrator().DropTable(&model.Rating{})

	userRepo = NewUserRepo(db)

	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.House{})
	db.AutoMigrate(&model.Feature{})
	db.AutoMigrate(&model.HouseHasFeatures{})
	db.AutoMigrate(&model.Transaction{})
	db.AutoMigrate(&model.Rating{})

	dummyUser := model.User{
		Email:    "test@gmail.com",
		Password: "test1234",
	}
	userRepo.Register(dummyUser)

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
}

func TestGetUser(t *testing.T) {
	configTest = config.GetConfig()
	db = util.InitDB(configTest)

	db.Migrator().DropTable(&model.User{})
	db.Migrator().DropTable(&model.House{})
	db.Migrator().DropTable(&model.Feature{})
	db.Migrator().DropTable(&model.HouseHasFeatures{})
	db.Migrator().DropTable(&model.Transaction{})
	db.Migrator().DropTable(&model.Rating{})

	userRepo = NewUserRepo(db)

	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.House{})
	db.AutoMigrate(&model.Feature{})
	db.AutoMigrate(&model.HouseHasFeatures{})
	db.AutoMigrate(&model.Transaction{})
	db.AutoMigrate(&model.Rating{})

	seed.UserSeed(db)

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
}

func TestUpdateUser(t *testing.T) {
	configTest = config.GetConfig()
	db = util.InitDB(configTest)

	db.Migrator().DropTable(&model.User{})
	db.Migrator().DropTable(&model.House{})
	db.Migrator().DropTable(&model.Feature{})
	db.Migrator().DropTable(&model.HouseHasFeatures{})
	db.Migrator().DropTable(&model.Transaction{})
	db.Migrator().DropTable(&model.Rating{})

	userRepo = NewUserRepo(db)

	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.House{})
	db.AutoMigrate(&model.Feature{})
	db.AutoMigrate(&model.HouseHasFeatures{})
	db.AutoMigrate(&model.Transaction{})
	db.AutoMigrate(&model.Rating{})

	seed.UserSeed(db)

	t.Run("Update User ", func(t *testing.T) {
		var mockUser model.User
		mockUser.Email = "test@gmail.com"
		mockUser.Password = "test4321"
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
}

func TestDeleteUser(t *testing.T) {
	configTest = config.GetConfig()
	db = util.InitDB(configTest)

	db.Migrator().DropTable(&model.User{})
	db.Migrator().DropTable(&model.House{})
	db.Migrator().DropTable(&model.Feature{})
	db.Migrator().DropTable(&model.HouseHasFeatures{})
	db.Migrator().DropTable(&model.Transaction{})
	db.Migrator().DropTable(&model.Rating{})

	userRepo = NewUserRepo(db)

	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.House{})
	db.AutoMigrate(&model.Feature{})
	db.AutoMigrate(&model.HouseHasFeatures{})
	db.AutoMigrate(&model.Transaction{})
	db.AutoMigrate(&model.Rating{})

	seed.UserSeed(db)

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
