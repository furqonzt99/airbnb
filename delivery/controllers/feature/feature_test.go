package feature

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/furqonzt99/airbnb/delivery/common"
	"github.com/furqonzt99/airbnb/model"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestFeatureGetAll(t *testing.T) {
	t.Run("Test Register", func(t *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()

		context := e.NewContext(req, res)
		context.SetPath("/features")

		featureController := NewFeatureControllers(mockFeatureRepository{})
		featureController.GetAllFeatureController()(context)

		response := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)
		assert.Equal(t, "Successful Operation", response.Message)
	})
}

type mockFeatureRepository struct{}

func (m mockFeatureRepository) GetAll() ([]model.Feature, error) {
	return []model.Feature{{Name: "wifi"}, {Name: "pool"}}, nil
}
