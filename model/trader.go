package model

import (
	"github.com/jinzhu/gorm"
	"github.com/miaolz123/samaritan/candyjs"
	"github.com/miaolz123/samaritan/constant"
)

// Trader struct
type Trader struct {
	gorm.Model
	UserID     uint       `gorm:"index"`
	StrategyID uint       `gorm:"index"`
	Name       string     `gorm:"type:varchar(200)"`
	Exchanges  []Exchange `gorm:"many2many:trader_exchanges"`

	Status   int              `gorm:"-"`
	Logger   Logger           `gorm:"-"`
	Strategy Strategy         `gorm:"-"`
	Ctx      *candyjs.Context `gorm:"-"`
}

// TraderExchange struct
type TraderExchange struct {
	gorm.Model
	TraderID   uint `gorm:"index"`
	ExchangeID uint `gorm:"index"`
}

// Run ...
func (t *Trader) Run() error {
	t.Status = 1
	defer t.stop()
	t.Logger.Log(constant.INFO, 0.0, 0.0, "Start Running")
	if err := t.Ctx.PevalString(t.Strategy.Script); err != nil {
		t.Logger.Log(constant.ERROR, 0.0, 0.0, err)
		return err
	}
	return nil
}

func (t *Trader) stop() {
	t.Ctx.Destroy()
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
		if TraderMap[t.ID] != nil {
			traders[i].Status = TraderMap[t.ID].Status
		}
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

// StopTrader ...
func StopTrader(trader Trader) {
	if TraderMap[trader.ID] != nil {
		TraderMap[trader.ID].stop()
	}
}
