package house

import (
	"testing"

	"github.com/furqonzt99/airbnb/config"
	"github.com/furqonzt99/airbnb/model"
	"github.com/furqonzt99/airbnb/repository/feature"
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
var houseRepo *HouseRepository

func TestCreateHouse(t *testing.T) {
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
	houseRepo = NewHouseRepo(db)

	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.House{})
	db.AutoMigrate(&model.Feature{})
	db.AutoMigrate(&model.HouseHasFeatures{})
	db.AutoMigrate(&model.Transaction{})
	db.AutoMigrate(&model.Rating{})

	seed.UserSeed(db)
	seed.FeatureSeed(db)

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
}

func TestGetAllHouse(t *testing.T) {
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
	houseRepo = NewHouseRepo(db)

	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.House{})
	db.AutoMigrate(&model.Feature{})
	db.AutoMigrate(&model.HouseHasFeatures{})
	db.AutoMigrate(&model.Transaction{})
	db.AutoMigrate(&model.Rating{})

	seed.UserSeed(db)
	seed.FeatureSeed(db)
	seed.HouseSeed(db)

	t.Run("Get All House", func(t *testing.T) {
		offset := 0
		pageSize := 10
		search := "rumah"
		city := "indonesia"

		_, err := houseRepo.GetAll(offset, pageSize, search, city)
		assert.Nil(t, err)
	})
}

func TestGetMyHouse(t *testing.T) {
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
	houseRepo = NewHouseRepo(db)

	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.House{})
	db.AutoMigrate(&model.Feature{})
	db.AutoMigrate(&model.HouseHasFeatures{})
	db.AutoMigrate(&model.Transaction{})
	db.AutoMigrate(&model.Rating{})

	seed.UserSeed(db)
	seed.FeatureSeed(db)
	seed.HouseSeed(db)

	t.Run("Get All My House", func(t *testing.T) {
		userId := 1
		_, err := houseRepo.GetAllMine(userId)
		assert.Nil(t, err)
	})
}

func TestGetHouse(t *testing.T) {
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
	houseRepo = NewHouseRepo(db)

	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.House{})
	db.AutoMigrate(&model.Feature{})
	db.AutoMigrate(&model.HouseHasFeatures{})
	db.AutoMigrate(&model.Transaction{})
	db.AutoMigrate(&model.Rating{})

	seed.UserSeed(db)
	seed.FeatureSeed(db)
	seed.HouseSeed(db)

	t.Run("Get House", func(t *testing.T) {
		houseId := 1
		_, err := houseRepo.Get(houseId)
		assert.Nil(t, err)
	})

	t.Run("Error Get House", func(t *testing.T) {
		houseId := 100
		_, err := houseRepo.Get(houseId)
		assert.NotNil(t, err)
	})
}

func TestUpdateHouse(t *testing.T) {
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
	houseRepo = NewHouseRepo(db)

	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.House{})
	db.AutoMigrate(&model.Feature{})
	db.AutoMigrate(&model.HouseHasFeatures{})
	db.AutoMigrate(&model.Transaction{})
	db.AutoMigrate(&model.Rating{})

	seed.UserSeed(db)
	seed.FeatureSeed(db)
	seed.HouseSeed(db)

	t.Run("Update House", func(t *testing.T) {
		var mockHouse model.House
		mockHouse.Title = "rumah2"
		mockHouse.Address = "jalan awal"
		mockHouse.City = "indonesia"
		mockHouse.Price = 100000

		houseId := 2
		userId := 4

		res, err := houseRepo.Update(mockHouse, houseId, userId)
		assert.Nil(t, err)
		assert.Equal(t, res.Title, "rumah2")
		assert.Equal(t, res.Address, "jalan awal")
	})

	t.Run("Error Update House", func(t *testing.T) {
		var mockHouse model.House
		mockHouse.Title = "rumah2"
		mockHouse.Address = "jalan awal"
		mockHouse.City = "indonesia"
		mockHouse.Price = 100000

		houseId := 100
		userId := 100

		_, err := houseRepo.Update(mockHouse, houseId, userId)
		assert.NotNil(t, err)
	})
}

func TestDeleteHouse(t *testing.T) {
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
	houseRepo = NewHouseRepo(db)

	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.House{})
	db.AutoMigrate(&model.Feature{})
	db.AutoMigrate(&model.HouseHasFeatures{})
	db.AutoMigrate(&model.Transaction{})
	db.AutoMigrate(&model.Rating{})

	seed.UserSeed(db)
	seed.FeatureSeed(db)
	seed.HouseSeed(db)

	t.Run("Delete House", func(t *testing.T) {
		houseId := 1
		userId := 2
		_, err := houseRepo.Delete(houseId, userId)
		assert.Nil(t, err)
	})

	t.Run("Error Delete House No ID Or UserID", func(t *testing.T) {
		houseId := 100
		userId := 100
		_, err := houseRepo.Delete(houseId, userId)
		assert.NotNil(t, err)
	})
}

func TestSaveHouseHasFeature(t *testing.T) {
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
	houseRepo = NewHouseRepo(db)

	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.House{})
	db.AutoMigrate(&model.Feature{})
	db.AutoMigrate(&model.HouseHasFeatures{})
	db.AutoMigrate(&model.Transaction{})
	db.AutoMigrate(&model.Rating{})

	seed.UserSeed(db)
	seed.FeatureSeed(db)
	seed.HouseSeed(db)

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
}

func TestDeleteHouseHasFeature(t *testing.T) {
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
	houseRepo = NewHouseRepo(db)

	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.House{})
	db.AutoMigrate(&model.Feature{})
	db.AutoMigrate(&model.HouseHasFeatures{})
	db.AutoMigrate(&model.Transaction{})
	db.AutoMigrate(&model.Rating{})

	seed.UserSeed(db)
	seed.FeatureSeed(db)
	seed.HouseSeed(db)

	t.Run("Delete House Has Feature", func(t *testing.T) {
		houseId := 1

		err := houseRepo.HouseHasFeatureDelete(houseId)
		assert.Equal(t, err, nil)
	})
}
