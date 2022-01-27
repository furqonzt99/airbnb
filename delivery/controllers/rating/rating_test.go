package rating

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/furqonzt99/airbnb/constant"
	"github.com/furqonzt99/airbnb/delivery/common"
	"github.com/furqonzt99/airbnb/delivery/controllers/user"
	"github.com/furqonzt99/airbnb/model"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

var jwtToken string

func TestCreateRating(t *testing.T) {
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

	t.Run("Test Create Rating", func(t *testing.T) {
		e := echo.New()
		e.Validator = &RatingValidator{Validator: validator.New()}

		requestBody, _ := json.Marshal(map[string]interface{}{
			"house_id": 1,
			"rating":   5,
			"comment":  "nyaman",
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/ratings")

		ratingController := NewRatingController(mockRatingRepository{})
		if err := middleware.JWT([]byte(constant.JWT_SECRET_KEY))(ratingController.Create)(context); err != nil {
			log.Fatal(err)
			return
		}

		response := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Successful Operation", response.Message)
	})

	t.Run("Error Test Create Rating", func(t *testing.T) {
		e := echo.New()
		e.Validator = &RatingValidator{Validator: validator.New()}

		requestBody, _ := json.Marshal(map[string]interface{}{
			"house_id": 1,
			"rating":   5,
			"comment":  "nyaman",
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/ratings")

		ratingController := NewRatingController(mockFalseRatingRepository{})
		if err := middleware.JWT([]byte(constant.JWT_SECRET_KEY))(ratingController.Create)(context); err != nil {
			log.Fatal(err)
			return
		}

		response := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Bad Request", response.Message)
	})

	t.Run("Error Test Create Rating Bind Error", func(t *testing.T) {
		e := echo.New()
		e.Validator = &RatingValidator{Validator: validator.New()}

		requestBody, _ := json.Marshal(map[string]interface{}{})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/ratings")

		ratingController := NewRatingController(mockFalseRatingRepository{})
		if err := middleware.JWT([]byte(constant.JWT_SECRET_KEY))(ratingController.Create)(context); err != nil {
			log.Fatal(err)
			return
		}

		response := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Bad Request", response.Message)
	})

	t.Run("Error Test Create Rating Validator Error", func(t *testing.T) {
		e := echo.New()
		e.Validator = &RatingValidator{Validator: validator.New()}

		requestBody, _ := json.Marshal(map[string]interface{}{
			"house_id": 1,
			"rating":   10,
			"comment":  "nyaman",
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/ratings")

		ratingController := NewRatingController(mockFalseRatingRepository{})
		if err := middleware.JWT([]byte(constant.JWT_SECRET_KEY))(ratingController.Create)(context); err != nil {
			log.Fatal(err)
			return
		}

		response := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Bad Request", response.Message)
	})
}

func TestUpdateRating(t *testing.T) {
	t.Run("Test Update Rating", func(t *testing.T) {
		e := echo.New()
		e.Validator = &RatingValidator{Validator: validator.New()}

		requestBody, _ := json.Marshal(map[string]interface{}{
			"rating":  5,
			"comment": "nyaman",
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/ratings/:houseId")
		context.SetParamNames("houseId")
		context.SetParamValues("1")

		ratingController := NewRatingController(mockRatingRepository{})
		if err := middleware.JWT([]byte(constant.JWT_SECRET_KEY))(ratingController.Update)(context); err != nil {
			log.Fatal(err)
			return
		}

		response := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Successful Operation", response.Message)
	})

	t.Run("Error Test Update Rating", func(t *testing.T) {
		e := echo.New()
		e.Validator = &RatingValidator{Validator: validator.New()}

		requestBody, _ := json.Marshal(map[string]interface{}{
			"rating":  5,
			"comment": "nyaman",
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/ratings/:houseId")
		context.SetParamNames("houseId")
		context.SetParamValues("1")

		ratingController := NewRatingController(mockFalseRatingRepository{})
		if err := middleware.JWT([]byte(constant.JWT_SECRET_KEY))(ratingController.Update)(context); err != nil {
			log.Fatal(err)
			return
		}

		response := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Not Found", response.Message)
	})

	t.Run("Error Test Update Rating Bind Error", func(t *testing.T) {
		e := echo.New()
		e.Validator = &RatingValidator{Validator: validator.New()}

		requestBody, _ := json.Marshal(map[string]interface{}{})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/ratings/:houseId")
		context.SetParamNames("houseId")
		context.SetParamValues("1")

		ratingController := NewRatingController(mockFalseRatingRepository{})
		if err := middleware.JWT([]byte(constant.JWT_SECRET_KEY))(ratingController.Update)(context); err != nil {
			log.Fatal(err)
			return
		}

		response := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Bad Request", response.Message)
	})

	t.Run("Error Test Update Rating Validator Error", func(t *testing.T) {
		e := echo.New()
		e.Validator = &RatingValidator{Validator: validator.New()}

		requestBody, _ := json.Marshal(map[string]interface{}{
			"rating":  10,
			"comment": "nyaman",
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/ratings/:houseId")
		context.SetParamNames("houseId")
		context.SetParamValues("1")

		ratingController := NewRatingController(mockFalseRatingRepository{})
		if err := middleware.JWT([]byte(constant.JWT_SECRET_KEY))(ratingController.Update)(context); err != nil {
			log.Fatal(err)
			return
		}

		response := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Bad Request", response.Message)
	})
}

func TestDeleteRating(t *testing.T) {
	t.Run("Test Delete Rating", func(t *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(http.MethodPost, "/", nil)
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/ratings/:houseId")
		context.SetParamNames("houseId")
		context.SetParamValues("1")

		ratingController := NewRatingController(mockRatingRepository{})
		if err := middleware.JWT([]byte(constant.JWT_SECRET_KEY))(ratingController.Delete)(context); err != nil {
			log.Fatal(err)
			return
		}

		response := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Successful Operation", response.Message)
	})

	t.Run("Error Test Delete Rating", func(t *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(http.MethodPost, "/", nil)
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/ratings/:houseId")
		context.SetParamNames("houseId")
		context.SetParamValues("1")

		ratingController := NewRatingController(mockFalseRatingRepository{})
		if err := middleware.JWT([]byte(constant.JWT_SECRET_KEY))(ratingController.Delete)(context); err != nil {
			log.Fatal(err)
			return
		}

		response := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Not Found", response.Message)
	})

	t.Run("Error Test Delete Rating No HouseId Parameter", func(t *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(http.MethodPost, "/", nil)
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/ratings/:houseId")

		ratingController := NewRatingController(mockFalseRatingRepository{})
		if err := middleware.JWT([]byte(constant.JWT_SECRET_KEY))(ratingController.Delete)(context); err != nil {
			log.Fatal(err)
			return
		}

		response := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Bad Request", response.Message)
	})
}

type mockUserRepository struct{}

func (m mockUserRepository) Register(newUser model.User) (model.User, error) {
	hash, _ := bcrypt.GenerateFromPassword([]byte(newUser.Password), 14)
	return model.User{Email: newUser.Email, Password: string(hash), Name: newUser.Name}, nil
}

func (m mockUserRepository) Login(email string) (model.User, error) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("test1234"), 14)
	return model.User{Email: "test@gmail.com", Password: string(hash), Name: "tester"}, nil
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

type mockRatingRepository struct{}

func (m mockRatingRepository) Create(model.Rating) (model.Rating, error) {
	return model.Rating{HouseID: 1, UserID: 1, Rating: 5, Comment: "nyaman"}, nil
}

func (m mockRatingRepository) Update(model.Rating) (model.Rating, error) {
	return model.Rating{HouseID: 1, UserID: 1, Rating: 5, Comment: "nyaman"}, nil
}

func (m mockRatingRepository) Delete(userId, houseId int) (model.Rating, error) {
	return model.Rating{HouseID: 1, UserID: 1, Rating: 5, Comment: "nyaman"}, nil
}

type mockFalseRatingRepository struct{}

func (m mockFalseRatingRepository) Create(model.Rating) (model.Rating, error) {
	return model.Rating{HouseID: 1, UserID: 1, Rating: 5, Comment: "nyaman"}, errors.New("Error")
}

func (m mockFalseRatingRepository) Update(model.Rating) (model.Rating, error) {
	return model.Rating{HouseID: 1, UserID: 1, Rating: 5, Comment: "nyaman"}, errors.New("Error")
}

func (m mockFalseRatingRepository) Delete(userId, houseId int) (model.Rating, error) {
	return model.Rating{HouseID: 1, UserID: 1, Rating: 5, Comment: "nyaman"}, errors.New("Error")
}
