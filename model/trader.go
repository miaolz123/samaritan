package model

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/miaolz123/samaritan/api"
	"github.com/miaolz123/samaritan/candyjs"
	"github.com/miaolz123/samaritan/log"
	"github.com/miaolz123/samaritan/task"
)

// Trader struct
type Trader struct {
	gorm.Model
	UserID     uint       `gorm:"index"`
	StrategyID uint       `gorm:"index"`
	Name       string     `gorm:"type:varchar(200)"`
	Exchanges  []Exchange `gorm:"many2many:trader_exchanges"`

	Status   int `gorm:"-"`
	strategy Strategy
	log      log.Logger
	ctx      *candyjs.Context
	runner   *task.Task
}

func (t *Trader) run() error {
	t.Status = 1
	t.runner.Add(1)
	defer t.stop()
	t.log.Do("info", 0.0, 0.0, "Start Running")
	if err := t.ctx.PevalString(t.strategy.Script); err != nil {
		t.log.Do("error", 0.0, 0.0, err)
		return err
	}
	return nil
}

func (t *Trader) stop() bool {
	if t.runner.AllDone() {
		t.log.Do("info", 0.0, 0.0, "Stop Running")
		t.Status = 0
		return true
	}
	return false
}

// TraderExchange struct
type TraderExchange struct {
	gorm.Model
	TraderID   uint `gorm:"index"`
	ExchangeID uint `gorm:"index"`
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
		if err = db.Model(&t).Association("Exchanges").Find(&traders[i].Exchanges).Error; err != nil {
			return
		}
	}
	return
}

// RunTrader ...
func RunTrader(trader Trader) (err error) {
	db, err := NewOrm()
	if err != nil {
		return
	}
	if err = db.First(&trader, trader.ID).Error; err != nil {
		return
	}
	if err = db.First(&trader.strategy, trader.StrategyID).Error; err != nil {
		return
	}
	if err = db.Model(&trader).Association("Exchanges").Find(&trader.Exchanges).Error; err != nil {
		return
	}
	trader.log = log.New("global")
	trader.ctx = candyjs.NewContext()
	trader.runner = task.New()
	constants := []string{
		"BTC",
		"LTC",
		"M",
		"M5",
		"M15",
		"M30",
		"H",
		"D",
		"W",
	}
	exchanges := []interface{}{}
	for _, e := range trader.Exchanges {
		opt := api.Option{
			Type:      e.Type,
			AccessKey: e.AccessKey,
			SecretKey: e.SecretKey,
			MainStock: "BTC",
		}
		switch opt.Type {
		case "okcoin.cn":
			exchanges = append(exchanges, api.NewOKCoinCn(opt))
		case "huobi":
			exchanges = append(exchanges, api.NewHuobi(opt))
		}
	}
	if len(exchanges) == 0 {
		trader.log.Do("error", 0.0, 0.0, "Please add at least one exchange")
	}
	for _, c := range constants {
		trader.ctx.PushGlobalInterface(c, c)
	}
	trader.ctx.PushGlobalGoFunction("Log", func(msgs ...interface{}) {
		trader.log.Do("info", 0.0, 0.0, msgs...)
	})
	trader.ctx.PushGlobalGoFunction("Sleep", func(t float64) {
		time.Sleep(time.Duration(t * 1000000))
	})
	trader.ctx.PushGlobalInterface("exchange", exchanges[0])
	trader.ctx.PushGlobalInterface("exchanges", exchanges)
	TraderMap[trader.ID] = &trader
	go TraderMap[trader.ID].run()
	return
}

// StopTrader ...
func StopTrader(trader Trader) bool {
	if TraderMap[trader.ID] != nil {
		return TraderMap[trader.ID].stop()
	}
	return true
}
