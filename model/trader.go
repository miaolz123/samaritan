package model

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/robertkrimen/otto"
)

// Trader struct
type Trader struct {
	gorm.Model
	UserID     uint   `gorm:"index"`
	StrategyID uint   `gorm:"index"`
	Name       string `gorm:"type:varchar(200)"`

	Exchanges []Exchange `gorm:"-"`
	Status    int        `gorm:"-"`
	Logger    Logger     `gorm:"-" json:"-"`
	Strategy  Strategy   `gorm:"-"`
	Ctx       *otto.Otto `gorm:"-" json:"-"`
}

// TraderExchange struct
type TraderExchange struct {
	ID         int64 `gorm:"primary_key;AUTO_INCREMENT"`
	TraderID   uint  `gorm:"index"`
	ExchangeID uint  `gorm:"index"`
	Exchange   `gorm:"-"`
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
	if err = DB.Where("user_id in (?)", userIDs).Order("id").Find(&traders).Error; err != nil {
		return
	}
	for i, t := range traders {
		if t.StrategyID > 0 {
			if err = DB.Where("id = ?", t.StrategyID).First(&traders[i].Strategy).Error; err != nil {
				return
			}
		}
		if err = DB.Raw(`SELECT e.* FROM exchanges e, trader_exchanges r WHERE r.trader_id
		= ? AND e.id = r.exchange_id`, t.ID).Scan(&traders[i].Exchanges).Error; err != nil {
			return
		}
	}
	return
}

// GetTrader ...
func GetTrader(self User, id interface{}) (trader Trader, err error) {
	if err = DB.Where("id = ?", id).First(&trader).Error; err != nil {
		return
	}
	user, err := GetUserByID(trader.UserID)
	if err != nil {
		return
	}
	if user.Level >= self.Level && user.ID != self.ID {
		err = fmt.Errorf("Insufficient permissions")
	}
	if trader.StrategyID > 0 {
		if err = DB.Where("id = ?", trader.StrategyID).First(&trader.Strategy).Error; err != nil {
			return
		}
	}
	err = DB.Raw(`SELECT e.* FROM exchanges e, trader_exchanges r WHERE r.trader_id
		= ? AND e.id = r.exchange_id`, trader.ID).Scan(&trader.Exchanges).Error
	return
}

// GetTraderExchanges ...
func GetTraderExchanges(self User, id interface{}) (traderExchanges []TraderExchange, err error) {
	if _, err = GetTrader(self, id); err != nil {
		return
	}
	if err = DB.Where("trader_id = ?", id).Find(&traderExchanges).Error; err != nil {
		return
	}
	for i, r := range traderExchanges {
		if err = DB.Where("id = ?", r.ExchangeID).Find(&traderExchanges[i].Exchange).Error; err != nil {
			return
		}
	}
	return
}
