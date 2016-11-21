package model

import (
	"fmt"
	"time"
)

// Exchange struct
type Exchange struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	UserID    uint       `gorm:"index"`
	Name      string     `gorm:"type:varchar(50)"`
	Type      string     `gorm:"type:varchar(50)"`
	AccessKey string     `gorm:"type:varchar(200)"`
	SecretKey string     `gorm:"type:varchar(200)"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `sql:"index" json:"-"`
}

// GetExchanges ...
func GetExchanges(self User) (exchanges []Exchange, err error) {
	users, err := GetUsers(self)
	if err != nil {
		return
	}
	userIDs := []uint{}
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
