package model

import (
	"time"
)

// Algorithm struct
type Algorithm struct {
	ID          int64      `gorm:"primary_key" json:"id"`
	UserID      int64      `gorm:"index" json:"userId"`
	Name        string     `gorm:"type:varchar(200)" json:"name"`
	Description string     `gorm:"type:text" json:"description"`
	Script      string     `gorm:"type:text" json:"script"`
	EvnDefault  string     `gorm:"type:text" json:"evnDefault"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	DeletedAt   *time.Time `sql:"index" json:"-"`

	Traders []Trader `gorm:"-" json:"traders"`
}

// ListAlgorithm ...
func (user User) ListAlgorithm(size, page int64, order string) (total int64, algorithms []Algorithm, err error) {
	_, users, err := user.ListUser(-1, 1, "id")
	if err != nil {
		return
	}
	userIDs := []int64{}
	for _, u := range users {
		userIDs = append(userIDs, u.ID)
	}
	err = DB.Model(&Algorithm{}).Where("user_id in (?)", userIDs).Count(&total).Error
	if err != nil {
		return
	}
	err = DB.Where("user_id in (?)", userIDs).Order(toUnderScoreCase(order)).Limit(size).Offset((page - 1) * size).Find(&algorithms).Error
	return
}
