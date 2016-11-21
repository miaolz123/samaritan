package model

import (
	"time"
)

// User struct
type User struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	Username  string     `gorm:"type:varchar(25);unique_index" json:"username"`
	Password  string     `gorm:"not null" json:"-"`
	Level     int        `json:"level"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `sql:"index" json:"-"`
}

// GetUserByID ...
func GetUserByID(id interface{}) (user User, err error) {
	err = DB.First(&user, id).Error
	return
}

// GetUser ...
func GetUser(username interface{}) (user User, err error) {
	err = DB.Where("username = ?", username).First(&user).Error
	return
}

// UserList ...
func (user User) UserList(size, page int64) (total int64, users []User, err error) {
	err = DB.Model(&User{}).Where("level < ? OR id = ?", user.Level, user.ID).Count(&total).Error
	if err != nil {
		return
	}
	err = DB.Order("id").Limit(size).Offset((page-1)*size).Where("level < ? OR id = ?", user.Level, user.ID).Find(&users).Error
	return
}

// GetUsers ...
func GetUsers(self User, order ...string) (users []User, err error) {
	orderKey := "id"
	if len(order) > 0 && order[0] != "" {
		orderKey = order[0]
	}
	err = DB.Order(orderKey).Where("level < ?", self.Level).Order("id").Find(&users).Error
	users = append([]User{self}, users...)
	return
}
