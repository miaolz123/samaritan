package model

import (
	"fmt"
	"time"
)

// Trader struct
type Trader struct {
	ID          int64      `gorm:"primary_key" json:"id"`
	UserID      int64      `gorm:"index" json:"userId"`
	AlgorithmID int64      `gorm:"index" json:"algorithmId"`
	Name        string     `gorm:"type:varchar(200)" json:"name"`
	Environment string     `gorm:"type:text" json:"environment"`
	LastRunAt   time.Time  `json:"lastRunAt"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	DeletedAt   *time.Time `sql:"index" json:"-"`

	Exchanges []Exchange `gorm:"-" json:"exchanges"`
	Status    int64      `gorm:"-" json:"status"`
	Algorithm Algorithm  `gorm:"-" json:"algorithm"`
}

// TraderExchange struct
type TraderExchange struct {
	ID         int64 `gorm:"primary_key"`
	TraderID   int64 `gorm:"index"`
	ExchangeID int64 `gorm:"index"`
	Exchange   `gorm:"-"`
}

// GetTraders ...
func (user User) GetTraders(algorithmID int64) (traders []Trader, err error) {
	err = DB.Where("algorithm_id = ?", algorithmID).Order("id desc").Find(&traders).Error
	return
}

// GetTraders ...
func GetTraders(self User) (traders []Trader, err error) {
	users, err := GetUsers(self)
	if err != nil {
		return
	}
	userIDs := []int64{}
	for _, u := range users {
		userIDs = append(userIDs, u.ID)
	}
	if err = DB.Where("user_id in (?)", userIDs).Order("id").Find(&traders).Error; err != nil {
		return
	}
	for i, t := range traders {
		if t.AlgorithmID > 0 {
			if err = DB.Where("id = ?", t.AlgorithmID).First(&traders[i].Algorithm).Error; err != nil {
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
	if trader.AlgorithmID > 0 {
		if err = DB.Where("id = ?", trader.AlgorithmID).First(&trader.Algorithm).Error; err != nil {
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
