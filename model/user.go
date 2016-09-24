package model

import "github.com/jinzhu/gorm"

// User struct
type User struct {
	gorm.Model
	Name     string `gorm:"type:varchar(25);unique_index"`
	Password string `gorm:"not null" json:"-"`
	Level    int
}

// GetUserByID ...
func GetUserByID(id interface{}) (user User, err error) {
	err = DB.Where("id = ?", id).First(&user).Error
	return
}

// GetUser ...
func GetUser(name interface{}) (user User, err error) {
	err = DB.Where("name = ?", name).First(&user).Error
	return
}

// GetUsers ...
func GetUsers(self User, order ...string) (users []User, err error) {
	orderKey := "id"
	if len(order) > 0 && order[0] != "" {
		orderKey = order[0]
	}
	err = DB.Order(orderKey).Where("level < ?", self.Level).Find(&users).Error
	users = append([]User{self}, users...)
	return
}
