package middleware

import (
	"errors"
	"time"

	"github.com/furqonzt99/airbnb/constant"
	"github.com/furqonzt99/airbnb/delivery/common"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func CreateToken(userId int, email string) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["userId"] = int(userId)
	claims["email"] = email
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(constant.JWT_SECRET_KEY))
}

func ExtractTokenUser(e echo.Context) (common.JWTPayload, error) {
	user := e.Get("user").(*jwt.Token)
	if user.Valid {
		claims := user.Claims.(jwt.MapClaims)
		userId := claims["userId"].(float64)
		email := claims["email"]
		return common.JWTPayload{
			UserID: int(userId),
			Email:  email.(string),
		}, nil
	}
	return common.JWTPayload{}, errors.New("invalid token")
}
