package model

import "github.com/jinzhu/gorm"

// Strategy struct
type Strategy struct {
	gorm.Model
	UserID      uint   `gorm:"index"`
	Name        string `gorm:"type:varchar(200);unique_index"`
	Description string `gorm:"type:text"`
	Script      string `gorm:"type:text"`
}

// GetStrategies ...
func GetStrategies(self User) (strategies []Strategy, err error) {
	users, err := GetUsers(self)
	if err != nil {
		return
	}
	userIDs := []uint{}
	for _, u := range users {
		userIDs = append(userIDs, u.ID)
	}
	db, err := NewOrm()
	if err != nil {
		return
	}
	err = db.Where("user_id in (?)", userIDs).Find(&strategies).Error
	return
}
