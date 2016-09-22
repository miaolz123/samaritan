package model

import (
	"github.com/jinzhu/gorm"
	"github.com/robertkrimen/otto"
)

// Trader struct
type Trader struct {
	gorm.Model
	UserID     uint       `gorm:"index"`
	StrategyID uint       `gorm:"index"`
	Name       string     `gorm:"type:varchar(200)"`
	Exchanges  []Exchange `gorm:"many2many:trader_exchanges"`

	Status   int        `gorm:"-"`
	Logger   Logger     `gorm:"-" json:"-"`
	Strategy Strategy   `gorm:"-"`
	Ctx      *otto.Otto `gorm:"-" json:"-"`
}

// TraderExchange struct
type TraderExchange struct {
	gorm.Model
	TraderID   uint `gorm:"index"`
	ExchangeID uint `gorm:"index"`
}

// GetTrader ...
func GetTrader(id interface{}) (trader Trader, err error) {
	db, err := NewOrm()
	if err != nil {
		return
	}
	err = db.Where("id = ?", id).First(&trader).Error
	return
}

// GetTraders ...
func GetTraders(self User) (traders []Trader, err error) {
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
	if err = db.Where("user_id in (?)", userIDs).Find(&traders).Error; err != nil {
		return
	}
	for i, t := range traders {
		if t.StrategyID > 0 {
			if err = db.Where("id = ?", t.StrategyID).First(&traders[i].Strategy).Error; err != nil {
				return
			}
		}
		if err = db.Model(&t).Association("Exchanges").Find(&traders[i].Exchanges).Error; err != nil {
			return
		}
	}
	return
}
