package rating

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
var ratingRepo *RatingRepository

func TestCreateRating(t *testing.T) {
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
	ratingRepo = NewRatingRepository(db)

	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.House{})
	db.AutoMigrate(&model.Feature{})
	db.AutoMigrate(&model.HouseHasFeatures{})
	db.AutoMigrate(&model.Transaction{})
	db.AutoMigrate(&model.Rating{})

	seed.UserSeed(db)
	seed.FeatureSeed(db)
	seed.HouseSeed(db)

	t.Run("Create Rating", func(t *testing.T) {
		var mockRating model.Rating
		mockRating.HouseID = 1
		mockRating.UserID = 1
		mockRating.Rating = 5
		mockRating.Comment = "mantap"

		res, err := ratingRepo.Create(mockRating)
		assert.Nil(t, err)
		assert.Equal(t, 1, int(res.HouseID))
		assert.Equal(t, 1, int(res.UserID))
	})

	t.Run("Error Create Rating", func(t *testing.T) {
		var mockRating model.Rating

		_, err := ratingRepo.Create(mockRating)
		assert.NotNil(t, err)
	})
}

func TestUpdateRating(t *testing.T) {
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
	ratingRepo = NewRatingRepository(db)

	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.House{})
	db.AutoMigrate(&model.Feature{})
	db.AutoMigrate(&model.HouseHasFeatures{})
	db.AutoMigrate(&model.Transaction{})
	db.AutoMigrate(&model.Rating{})

	seed.UserSeed(db)
	seed.FeatureSeed(db)
	seed.HouseSeed(db)

	dummyRating := model.Rating{
		HouseID: 1,
		UserID:  1,
		Rating:  5,
		Comment: "nyaman",
	}
	ratingRepo.Create(dummyRating)

	t.Run("Update Rating", func(t *testing.T) {
		var mockRating model.Rating
		mockRating.HouseID = 1
		mockRating.UserID = 1
		mockRating.Rating = 3
		mockRating.Comment = "biasa"

		res, err := ratingRepo.Update(mockRating)
		assert.Nil(t, err)
		assert.Equal(t, 3, res.Rating)
	})

	t.Run("Error Update Rating", func(t *testing.T) {
		var mockRating model.Rating
		mockRating.HouseID = 100
		mockRating.UserID = 100
		mockRating.Rating = 3
		mockRating.Comment = "biasa"

		_, err := ratingRepo.Update(mockRating)
		assert.NotNil(t, err)
	})
}

func TestDeleteRating(t *testing.T) {
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
	ratingRepo = NewRatingRepository(db)

	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.House{})
	db.AutoMigrate(&model.Feature{})
	db.AutoMigrate(&model.HouseHasFeatures{})
	db.AutoMigrate(&model.Transaction{})
	db.AutoMigrate(&model.Rating{})

	seed.UserSeed(db)
	seed.FeatureSeed(db)
	seed.HouseSeed(db)

	dummyRating := model.Rating{
		HouseID: 1,
		UserID:  1,
		Rating:  5,
		Comment: "nyaman",
	}
	ratingRepo.Create(dummyRating)

	t.Run("Delete Rating", func(t *testing.T) {
		houseId := 1
		userId := 1

		_, err := ratingRepo.Delete(userId, houseId)
		assert.Nil(t, err)
	})

	t.Run("Error Delete Rating", func(t *testing.T) {
		houseId := 100
		userId := 100

		_, err := ratingRepo.Delete(userId, houseId)
		assert.NotNil(t, err)
	})
}
