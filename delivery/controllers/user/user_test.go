package user

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
	mw "github.com/furqonzt99/airbnb/delivery/middleware"
	"github.com/furqonzt99/airbnb/model"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

var jwtToken string

func TestRegisterUser(t *testing.T) {
	t.Run("Test Register", func(t *testing.T) {
		e := echo.New()
		e.Validator = &UserValidator{Validator: validator.New()}

		requestBody, _ := json.Marshal(map[string]string{
			"email":    "test@gmail.com",
			"password": "test1234",
			"name":     "tester",
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)
		context.SetPath("/register")

		userController := NewUsersControllers(mockUserRepository{})
		userController.RegisterController()(context)

		response := RegisterUserResponseFormat{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Successful Operation", response.Message)
	})

	t.Run("Error Test Register Password Length Below 8", func(t *testing.T) {
		e := echo.New()
		e.Validator = &UserValidator{Validator: validator.New()}

		requestBody, _ := json.Marshal(map[string]string{
			"email":    "test@gmail.com",
			"password": "test",
			"name":     "tester",
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)
		context.SetPath("/register")

		userController := NewUsersControllers(mockFalseUserRepository{})
		userController.RegisterController()(context)

		response := RegisterUserResponseFormat{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Bad Request", response.Message)
	})

	t.Run("Error Test Email Already Exist", func(t *testing.T) {
		e := echo.New()
		e.Validator = &UserValidator{Validator: validator.New()}

		requestBody, _ := json.Marshal(map[string]string{
			"email":    "test@gmail.com",
			"password": "test1234",
			"name":     "tester",
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")

		context := e.NewContext(req, res)
		context.SetPath("/register")

		userController := NewUsersControllers(mockFalseUserRepository{})
		userController.RegisterController()(context)

		response := RegisterUserResponseFormat{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Email already exist", response.Message)
	})
}

func TestLoginUser(t *testing.T) {
	t.Run("Test Login", func(t *testing.T) {
		e := echo.New()
		e.Validator = &UserValidator{Validator: validator.New()}

		requestBody, _ := json.Marshal(map[string]string{
			"email":    "test@gmail.com",
			"password": "test1234",
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		context := e.NewContext(req, res)
		context.SetPath("/login")

		userController := NewUsersControllers(mockUserRepository{})
		userController.LoginController()(context)

		response := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		jwtToken = response.Data.(string)
		assert.Equal(t, "Successful Operation", response.Message)
		assert.NotNil(t, jwtToken)
	})

	t.Run("Error Test Login Password Length Below 8", func(t *testing.T) {
		e := echo.New()
		e.Validator = &UserValidator{Validator: validator.New()}

		requestBody, _ := json.Marshal(map[string]string{
			"email":    "test@gmail.com",
			"password": "test",
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		context := e.NewContext(req, res)
		context.SetPath("/login")

		userController := NewUsersControllers(mockFalseUserRepository{})
		userController.LoginController()(context)

		response := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Bad Request", response.Message)
	})

	t.Run("Error Test Login Wrong Password", func(t *testing.T) {
		e := echo.New()
		e.Validator = &UserValidator{Validator: validator.New()}

		requestBody, _ := json.Marshal(map[string]string{
			"email":    "test@gmail.com",
			"password": "test1234",
		})

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		context := e.NewContext(req, res)
		context.SetPath("/login")

		userController := NewUsersControllers(mockFalseUserRepository{})
		userController.LoginController()(context)

		response := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Wrong Password", response.Message)
	})
}

func TestGetUser(t *testing.T) {
	t.Run("Test Get User", func(t *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		res := httptest.NewRecorder()

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/profile")

		userController := NewUsersControllers(mockUserRepository{})
		if err := middleware.JWT([]byte(constant.JWT_SECRET_KEY))(userController.GetUserController())(context); err != nil {
			log.Fatal(err)
			return
		}

		response := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)
		assert.Equal(t, "Successful Operation", response.Message)
		assert.Equal(t, response.Data.(map[string]interface{})["name"], "tester")
	})
}

func TestUpdateUser(t *testing.T) {
	t.Run("Test Update", func(t *testing.T) {
		e := echo.New()
		e.Validator = &UserValidator{Validator: validator.New()}

		requestBody, _ := json.Marshal(map[string]string{
			"email":    "test2@gmail.com",
			"password": "test4321",
			"name":     "tester2",
		})

		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewBuffer(requestBody))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/users")

		userController := NewUsersControllers(mockUserRepository{})
		if err := middleware.JWT([]byte(constant.JWT_SECRET_KEY))(userController.UpdateUserController())(context); err != nil {
			log.Fatal(err)
			return
		}

		response := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Successful Operation", response.Message)
	})

	t.Run("Error Test Update Password Length Below 8", func(t *testing.T) {
		e := echo.New()
		e.Validator = &UserValidator{Validator: validator.New()}

		requestBody, _ := json.Marshal(map[string]string{
			"email":    "test2@gmail.com",
			"password": "test",
			"name":     "tester2",
		})

		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewBuffer(requestBody))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/users")

		userController := NewUsersControllers(mockFalseUserRepository{})
		if err := middleware.JWT([]byte(constant.JWT_SECRET_KEY))(userController.UpdateUserController())(context); err != nil {
			log.Fatal(err)
			return
		}

		response := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Bad Request", response.Message)
	})

	t.Run("Error Test Update User Not Found", func(t *testing.T) {
		e := echo.New()
		e.Validator = &UserValidator{Validator: validator.New()}

		requestBody, _ := json.Marshal(map[string]string{
			"email":    "test2@gmail.com",
			"password": "test1234",
			"name":     "tester2",
		})

		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewBuffer(requestBody))
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/users")

		userController := NewUsersControllers(mockFalseUserRepository{})
		if err := middleware.JWT([]byte(constant.JWT_SECRET_KEY))(userController.UpdateUserController())(context); err != nil {
			log.Fatal(err)
			return
		}

		response := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "User not found", response.Message)
	})
}

func TestDeleteUser(t *testing.T) {
	t.Run("Test Delete", func(t *testing.T) {
		e := echo.New()
		mw.LogMiddleware(e)
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		res := httptest.NewRecorder()

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken))

		context := e.NewContext(req, res)
		context.SetPath("/users")

		userController := NewUsersControllers(mockUserRepository{})
		if err := middleware.JWT([]byte(constant.JWT_SECRET_KEY))(userController.DeleteUserController())(context); err != nil {
			log.Fatal(err)
			return
		}

		response := common.ResponseSuccess{}
		json.Unmarshal([]byte(res.Body.Bytes()), &response)

		assert.Equal(t, "Successful Operation", response.Message)
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

type mockFalseUserRepository struct{}

func (m mockFalseUserRepository) Register(newUser model.User) (model.User, error) {
	hash, _ := bcrypt.GenerateFromPassword([]byte(newUser.Password), 14)
	return model.User{Email: newUser.Email, Password: string(hash), Name: newUser.Name}, errors.New("Email already exist")
}

func (m mockFalseUserRepository) Login(email string) (model.User, error) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("test4321"), 14)
	return model.User{Email: "test@gmail.com", Password: string(hash), Name: "tester"}, nil
}

func (m mockFalseUserRepository) Get(userid int) (model.User, error) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("test1234"), 14)
	return model.User{Email: "test@gmail.com", Password: string(hash), Name: "tester"}, errors.New("False Login Object")
}

func (m mockFalseUserRepository) Update(newUser model.User, userId int) (model.User, error) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("test4321"), 14)
	return model.User{Email: "test2@gmail.com", Password: string(hash), Name: "tester2"}, errors.New("False Login Object")
}

func (m mockFalseUserRepository) Delete(userId int) (model.User, error) {
	hash, _ := bcrypt.GenerateFromPassword([]byte("test4321"), 14)
	return model.User{Email: "test2@gmail.com", Password: string(hash), Name: "tester2"}, errors.New("False Login Object")
}
