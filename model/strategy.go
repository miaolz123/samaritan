package model

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

// Strategy struct
type Strategy struct {
	gorm.Model
	UserID      int64  `gorm:"index"`
	Name        string `gorm:"type:varchar(200)"`
	Description string `gorm:"type:text"`
	Script      string `gorm:"type:text"`
}

// GetStrategies ...
func GetStrategies(self User) (strategies []Strategy, err error) {
	users, err := GetUsers(self)
	if err != nil {
		return
	}
	userIDs := []int64{}
	for _, u := range users {
		userIDs = append(userIDs, u.ID)
	}
	err = DB.Where("user_id in (?)", userIDs).Order("id").Find(&strategies).Error
	return
}

// GetStrategy ...
func GetStrategy(self User, id interface{}) (strategy Strategy, err error) {
	if err = DB.First(&strategy, id).Error; err != nil {
		return
	}
	user, err := GetUserByID(strategy.UserID)
	if err != nil {
		return
	}
	if user.Level >= self.Level && user.ID != self.ID {
		err = fmt.Errorf("Insufficient permissions")
	}
	return
}
