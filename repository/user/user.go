package user

import (
	"github.com/furqonzt99/airbnb/model"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (ur *UserRepository) Register(newUser model.User) (model.User, error) {
	err := ur.db.Save(&newUser).Error
	if err != nil {
		return newUser, err
	}
	return newUser, nil
}

func (ur *UserRepository) Login(email string) (model.User, error) {
	var user model.User
	var err = ur.db.First(&user, "email = ?", email).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

func (ur *UserRepository) Get(userId int) (model.User, error) {
	user := model.User{}
	if err := ur.db.First(&user, userId).Error; err != nil {
		return user, err
	}
	return user, nil
}

func (ur *UserRepository) Update(newUser model.User, userId int) (model.User, error) {
	user := model.User{}
	if err := ur.db.First(&user, "id=?", userId).Error; err != nil {
		return newUser, err
	}
	ur.db.Model(&user).Updates(newUser)
	return newUser, nil
}

func (ur *UserRepository) Delete(userId int) (model.User, error) {
	user := model.User{}
	if err := ur.db.First(&user, "id=?", userId).Error; err != nil {
		return user, err
	}
	ur.db.Delete(&user)
	return user, nil
}
