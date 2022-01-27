package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/furqonzt99/airbnb/constant"
	"github.com/furqonzt99/airbnb/delivery/common"
	"github.com/furqonzt99/airbnb/delivery/controllers/house"
	"github.com/furqonzt99/airbnb/delivery/controllers/user"
	"github.com/furqonzt99/airbnb/model"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/assert"
)

func TestHouse(t *testing.T) {
	t.Run("Test Login", func(t *testing.T) {
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
		assert.Equal(t, "Successful Operation", response.Message)
		assert.NotNil(t, jwtToken)
	})

	t.Run("Test Create House", func(t *testing.T) {
		e := echo.New()
		e.Validator = &house.HouseValidator{Validator: validator.New()}

		requestBody, _ := json.Marshal(map[string]interface{}{
			"title":    "Rumah Bagus",
			"address":  "Jalan Ujung",
			"city":     "Indonesia",
			"price":    100000,
			"features": []int{1, 2},
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/houses")

		houseController := house.NewHouseControllers(mockHouseRepository{})
		if err := middleware.JWT([]byte(constant.JWT_SECRET_KEY))(houseController.CreateHouseController())(context); err != nil {
			log.Fatal(err)
			return
		}

		response := house.CreateHouseResponseFormat{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Successful Operation", response.Message)
	})

	t.Run("Test False Create House", func(t *testing.T) {
		e := echo.New()
		e.Validator = &house.HouseValidator{Validator: validator.New()}

		requestBody, _ := json.Marshal(map[string]interface{}{
			"house_id":   1,
			"feature_id": 2,
			"city":       "Indonesia",
			"price":      100000,
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/houses")

		houseController := house.NewHouseControllers(mockFalseHouseRepository{})
		if err := middleware.JWT([]byte(constant.JWT_SECRET_KEY))(houseController.CreateHouseController())(context); err != nil {
			log.Fatal(err)
			return
		}

		response := house.CreateHouseResponseFormat{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Bad Request", response.Message)
	})

	t.Run("Test Get All House", func(t *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(http.MethodPost, "/", nil)
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)
		context.SetPath("/houses")
		context.SetParamNames("name")
		context.SetParamValues("Rumah")

		houseController := house.NewHouseControllers(mockHouseRepository{})
		houseController.GetAllHouseController()(context)

		response := house.GetAllHouseResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)
		assert.Equal(t, "Successful Operation", response.Message)
		assert.Equal(t, "Rumah Bagus", response.Data[0].Title)
	})

	t.Run("Error Test Get All House", func(t *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(http.MethodPost, "/", nil)
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)
		context.SetPath("/houses")
		context.SetParamNames("name")
		context.SetParamValues("Rumah")

		houseController := house.NewHouseControllers(mockFalseHouseRepository{})
		houseController.GetAllHouseController()(context)

		response := house.GetAllHouseResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)
		assert.Equal(t, "Successful Operation", response.Message)
		assert.Equal(t, "Rumah Bagus", response.Data[0].Title)
	})

	t.Run("Test Get My House", func(t *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(http.MethodPost, "/", nil)
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/myhouses")

		houseController := house.NewHouseControllers(mockHouseRepository{})
		if err := middleware.JWT([]byte(constant.JWT_SECRET_KEY))(houseController.GetMyHouseController())(context); err != nil {
			log.Fatal(err)
			return
		}

		response := house.GetAllHouseResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)
		assert.Equal(t, "Successful Operation", response.Message)
		assert.Equal(t, "Rumah Bagus", response.Data[0].Title)
	})

	t.Run("Test Get House", func(t *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(http.MethodPost, "/", nil)
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)
		context.SetPath("/houses/:id")

		houseController := house.NewHouseControllers(mockHouseRepository{})
		houseController.GetHouseController()(context)

		response := house.GetHouseResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)
		assert.Equal(t, "Successful Operation", response.Message)
		assert.Equal(t, "Rumah Bagus", response.Data.Title)
	})

	t.Run("Error Test Get House", func(t *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(http.MethodPost, "/", nil)
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)
		context.SetPath("/houses/:id")

		houseController := house.NewHouseControllers(mockFalseHouseRepository{})
		houseController.GetHouseController()(context)

		response := house.GetHouseResponseFormat{}

		json.Unmarshal([]byte(res.Body.Bytes()), &response)
		assert.Equal(t, "Not Found", response.Message)
	})

	t.Run("Test Update House", func(t *testing.T) {
		e := echo.New()
		e.Validator = &house.HouseValidator{Validator: validator.New()}

		requestBody, _ := json.Marshal(map[string]interface{}{
			"title":    "Rumah Jelek",
			"address":  "Jalan Awal",
			"city":     "Bikini Bottom",
			"price":    200000,
			"features": []int{1, 2},
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/houses/:id")

		houseController := house.NewHouseControllers(mockHouseRepository{})
		if err := middleware.JWT([]byte(constant.JWT_SECRET_KEY))(houseController.UpdateHouseController())(context); err != nil {
			log.Fatal(err)
			return
		}

		response := house.CreateHouseResponseFormat{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Successful Operation", response.Message)
	})

	t.Run("Error Test Update House", func(t *testing.T) {
		e := echo.New()
		e.Validator = &house.HouseValidator{Validator: validator.New()}

		requestBody, _ := json.Marshal(map[string]interface{}{
			"title":    "Rumah Jelek",
			"address":  "Jalan Awal",
			"city":     "Bikini Bottom",
			"price":    200000,
			"features": []int{1, 2},
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/houses/:id")

		houseController := house.NewHouseControllers(mockFalseHouseRepository{})
		if err := middleware.JWT([]byte(constant.JWT_SECRET_KEY))(houseController.UpdateHouseController())(context); err != nil {
			log.Fatal(err)
			return
		}

		response := house.CreateHouseResponseFormat{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Bad Request", response.Message)
	})

	t.Run("Test Delete House", func(t *testing.T) {
		e := echo.New()
		e.Validator = &house.HouseValidator{Validator: validator.New()}

		req := httptest.NewRequest(http.MethodPost, "/", nil)
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/houses/:id")

		houseController := house.NewHouseControllers(mockHouseRepository{})
		if err := middleware.JWT([]byte(constant.JWT_SECRET_KEY))(houseController.DeleteHouseController())(context); err != nil {
			log.Fatal(err)
			return
		}

		response := house.CreateHouseResponseFormat{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Successful Operation", response.Message)
	})

	t.Run("Error Test Delete House", func(t *testing.T) {
		e := echo.New()
		e.Validator = &house.HouseValidator{Validator: validator.New()}

		req := httptest.NewRequest(http.MethodPost, "/", nil)
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/houses/:id")

		houseController := house.NewHouseControllers(mockFalseHouseRepository{})
		if err := middleware.JWT([]byte(constant.JWT_SECRET_KEY))(houseController.DeleteHouseController())(context); err != nil {
			log.Fatal(err)
			return
		}

		response := house.CreateHouseResponseFormat{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Not Found", response.Message)
	})
}

type mockHouseRepository struct{}

func (m mockHouseRepository) Create(newHouse model.House) (model.House, error) {
	return model.House{UserID: 1, Title: "Rumah Bagus", Address: "Jalan Ujung", City: "Indonesia", Price: 100000, Status: "open"}, nil
}

func (m mockHouseRepository) GetAll(offset, pageSize int, search, city string) ([]model.House, error) {
	return []model.House{{UserID: 1, Title: "Rumah Bagus", Address: "Jalan Ujung", City: "Indonesia", Price: 100000, Status: "open", Features: []model.Feature{{Name: "wifi"}}}}, nil
}

func (m mockHouseRepository) GetAllMine(userId int) ([]model.House, error) {
	return []model.House{{UserID: 1, Title: "Rumah Bagus", Address: "Jalan Ujung", City: "Indonesia", Price: 100000, Status: "open", Features: []model.Feature{{Name: "wifi"}}}}, nil
}

func (m mockHouseRepository) Get(houseId int) (model.House, error) {
	return model.House{UserID: 1, Title: "Rumah Bagus", Address: "Jalan Ujung", City: "Indonesia", Price: 100000, Status: "open", Features: []model.Feature{{Name: "wifi"}}}, nil
}

func (m mockHouseRepository) Update(newHouse model.House, houseId, userId int) (model.House, error) {
	return model.House{UserID: 1, Title: "Rumah Jelek", Address: "Jalan Awal", City: "Bikini Bottom", Price: 200000, Status: "open", Features: []model.Feature{{Name: "wifi"}}}, nil
}

func (m mockHouseRepository) Delete(houseId, userId int) (model.House, error) {
	return model.House{UserID: 1, Title: "Rumah Jelek", Address: "Jalan Awal", City: "Bikini Bottom", Price: 200000, Status: "open"}, nil
}

func (m mockHouseRepository) HouseHasFeature(houseHasFeature model.HouseHasFeatures) error {
	return nil
}

func (m mockHouseRepository) HouseHasFeatureDelete(houseId int) error {
	return nil
}

type mockFalseHouseRepository struct{}

func (m mockFalseHouseRepository) Create(newHouse model.House) (model.House, error) {
	return model.House{UserID: 1, Title: "Rumah Bagus", Address: "Jalan Ujung", City: "Indonesia", Price: 100000, Status: "open"}, errors.New("Error")
}

func (m mockFalseHouseRepository) GetAll(offset, pageSize int, search, city string) ([]model.House, error) {
	return []model.House{{UserID: 1, Title: "Rumah Bagus", Address: "Jalan Ujung", City: "Indonesia", Price: 100000, Status: "open", Features: []model.Feature{{Name: "wifi"}}}}, errors.New("Error")
}

func (m mockFalseHouseRepository) GetAllMine(userId int) ([]model.House, error) {
	return []model.House{{UserID: 1, Title: "Rumah Bagus", Address: "Jalan Ujung", City: "Indonesia", Price: 100000, Status: "open", Features: []model.Feature{{Name: "wifi"}}}}, errors.New("Error")
}

func (m mockFalseHouseRepository) Get(houseId int) (model.House, error) {
	return model.House{UserID: 1, Title: "Rumah Bagus", Address: "Jalan Ujung", City: "Indonesia", Price: 100000, Status: "open", Features: []model.Feature{{Name: "wifi"}}}, errors.New("Error")
}

func (m mockFalseHouseRepository) Update(newHouse model.House, houseId, userId int) (model.House, error) {
	return model.House{UserID: 1, Title: "Rumah Jelek", Address: "Jalan Awal", City: "Bikini Bottom", Price: 200000, Status: "open", Features: []model.Feature{{Name: "wifi"}}}, errors.New("Error")
}

func (m mockFalseHouseRepository) Delete(houseId, userId int) (model.House, error) {
	return model.House{UserID: 1, Title: "Rumah Jelek", Address: "Jalan Awal", City: "Bikini Bottom", Price: 200000, Status: "open"}, errors.New("Error")
}

func (m mockFalseHouseRepository) HouseHasFeature(houseHasFeature model.HouseHasFeatures) error {
	return nil
}

func (m mockFalseHouseRepository) HouseHasFeatureDelete(houseId int) error {
	return nil
}
