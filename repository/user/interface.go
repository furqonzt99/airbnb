package user

import "github.com/furqonzt99/airbnb/model"

type UserInterface interface {
	Register(newUser model.User) (model.User, error)
	Login(email string) (model.User, error)
	Get(userId int) (model.User, error)
	Update(newUser model.User, userId int) (model.User, error)
	Delete(userId int) (model.User, error)
}
