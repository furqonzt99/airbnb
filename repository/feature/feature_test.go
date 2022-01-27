package feature

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
var featureRepo *FeatureRepository

func TestGetAllFeature(t *testing.T) {
	configTest = config.GetConfig()
	db = util.InitDB(configTest)

	db.Migrator().DropTable(&model.User{})
	db.Migrator().DropTable(&model.House{})
	db.Migrator().DropTable(&model.Feature{})
	db.Migrator().DropTable(&model.HouseHasFeatures{})
	db.Migrator().DropTable(&model.Transaction{})
	db.Migrator().DropTable(&model.Rating{})

	featureRepo = NewFeatureRepo(db)

	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.House{})
	db.AutoMigrate(&model.Feature{})
	db.AutoMigrate(&model.HouseHasFeatures{})
	db.AutoMigrate(&model.Transaction{})
	db.AutoMigrate(&model.Rating{})

	seed.FeatureSeed(db)

	t.Run("Get All Features", func(t *testing.T) {
		res, err := featureRepo.GetAll()
		assert.Nil(t, err)
		assert.Equal(t, true, len(res) > 0)
	})
}
