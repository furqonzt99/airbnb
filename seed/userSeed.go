package seed

import (
	"fmt"

	"github.com/furqonzt99/airbnb/helper"
	"github.com/furqonzt99/airbnb/model"
	"gorm.io/gorm"
)

func UserSeed(db *gorm.DB) {
	password, _ := helper.Hashpwd("1234qwer")
	for i := 1; i <= 5; i++ {
		user := model.User{
			Name:     "User " + fmt.Sprint(i),
			Email:    fmt.Sprintf("user%v@gmail.com", i),
			Password: password,
		}
		db.Create(&user)
	}
}
