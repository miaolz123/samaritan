package model

import (
	"fmt"
	"time"
)

// Exchange struct
type Exchange struct {
	ID        int64      `gorm:"primary_key" json:"id"`
	UserID    int64      `gorm:"index" json:"userID"`
	Name      string     `gorm:"type:varchar(50)" json:"name"`
	Type      string     `gorm:"type:varchar(50)" json:"type"`
	AccessKey string     `gorm:"type:varchar(200)" json:"accessKey"`
	SecretKey string     `gorm:"type:varchar(200)" json:"secretKey"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `sql:"index" json:"-"`
}

// ExchangeList ...
func (user User) ExchangeList(size, page int64) (total int64, exchanges []Exchange, err error) {
	_, users, err := user.UserList(-1, 1)
	if err != nil {
		return
	}
	userIDs := []int64{}
	for _, u := range users {
		userIDs = append(userIDs, u.ID)
	}
	err = DB.Model(&Exchange{}).Where("user_id in (?)", userIDs).Count(&total).Error
	if err != nil {
		return
	}
	err = DB.Where("user_id in (?)", userIDs).Order("id").Limit(size).Offset((page - 1) * size).Find(&exchanges).Error
	return
}

// GetExchanges ...
func GetExchanges(self User) (exchanges []Exchange, err error) {
	users, err := GetUsers(self)
	if err != nil {
		return
	}
	userIDs := []int64{}
	for _, u := range users {
		userIDs = append(userIDs, u.ID)
	}
	err = DB.Where("user_id in (?)", userIDs).Order("id").Find(&exchanges).Error
	return
}

// GetExchange ...
func GetExchange(self User, id interface{}) (exchange Exchange, err error) {
	if err = DB.First(&exchange, id).Error; err != nil {
		return
	}
	user, err := GetUserByID(exchange.UserID)
	if err != nil {
		return
	}
	if user.Level >= self.Level && user.ID != self.ID {
		err = fmt.Errorf("Insufficient permissions")
	}
	return
}
