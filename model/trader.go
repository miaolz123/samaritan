package model

import (
	"fmt"
	"time"

	"github.com/miaolz123/samaritan/constant"
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

	Exchange `gorm:"-"`
}

// ListTrader ...
func (user User) ListTrader(algorithmID int64) (traders []Trader, err error) {
	err = DB.Where("user_id = ? AND algorithm_id = ?", user.ID, algorithmID).Find(&traders).Error
	for i, t := range traders {
		if err = DB.Raw(`SELECT e.* FROM exchanges e, trader_exchanges r WHERE r.trader_id
		= ? AND e.id = r.exchange_id`, t.ID).Scan(&traders[i].Exchanges).Error; err != nil {
			return
		}
	}
	return
}

// GetTrader ...
func (user User) GetTrader(id interface{}) (trader Trader, err error) {
	if err = DB.Where("id = ?", id).First(&trader).Error; err != nil {
		return
	}
	self, err := GetUserByID(trader.UserID)
	if err != nil {
		return
	}
	if user.Level < self.Level || user.ID != self.ID {
		err = fmt.Errorf(constant.ErrInsufficientPermissions)
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
func (user User) GetTraderExchanges(id interface{}) (traderExchanges []TraderExchange, err error) {
	if _, err = user.GetTrader(id); err != nil {
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

// UpdateTrader ...
func (user User) UpdateTrader(req Trader) (err error) {
	db, err := NewOrm()
	if err != nil {
		return err
	}
	defer db.Close()
	db = db.Begin()
	runner := Trader{}
	if err := db.First(&runner, req.ID).Error; err != nil {
		db.Rollback()
		return err
	}
	runner.Name = req.Name
	runner.Environment = req.Environment
	rs, err := user.GetTraderExchanges(runner.ID)
	if err != nil {
		db.Rollback()
		return err
	}
	for i, r := range rs {
		if i >= len(req.Exchanges) {
			if err := db.Delete(&r).Error; err != nil {
				db.Rollback()
				return err
			}
			continue
		}
		if r.Exchange.ID == req.Exchanges[i].ID {
			continue
		}
		r.ExchangeID = req.Exchanges[i].ID
		if err := db.Save(&r).Error; err != nil {
			db.Rollback()
			return err
		}
	}
	for i, e := range req.Exchanges {
		if i < len(rs) {
			continue
		}
		r := TraderExchange{
			TraderID:   runner.ID,
			ExchangeID: e.ID,
		}
		if err := db.Create(&r).Error; err != nil {
			db.Rollback()
			return err
		}
	}
	if err := db.Save(&runner).Error; err != nil {
		db.Rollback()
		return err
	}
	if err := db.Commit().Error; err != nil {
		db.Rollback()
		return err
	}
	return
}
