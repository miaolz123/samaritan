package model

import "github.com/jinzhu/gorm"

// Exchange struct
type Exchange struct {
	gorm.Model
	UserID    uint   `gorm:"index"`
	Name      string `gorm:"type:varchar(50);unique_index"`
	Type      string `gorm:"type:varchar(50)"`
	AccessKey string `gorm:"type:varchar(200)"`
	SecretKey string `gorm:"type:varchar(200)"`
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
	err = DB.Where("user_id in (?)", userIDs).Find(&exchanges).Error
	return
}
