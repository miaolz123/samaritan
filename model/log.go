package model

import (
	"fmt"
	"time"
)

// Log struct
type Log struct {
	ID           uint64 `gorm:"primary_key;AUTO_INCREMENT"`
	TraderID     uint   `gorm:"index"`
	Timestamp    int64
	ExchangeType string `gorm:"type:varchar(50)"`
	Type         int    // ["info", "profit", "buy", "sell", "cancel"]
	Price        float64
	Amount       float64
	Message      string `gorm:"type:text"`

	Time string `gorm:"-"`
}

// Logger struct
type Logger struct {
	TraderID     uint
	ExchangeType string
}

// GetLogs ...
func GetLogs(self User, traderID interface{}, page, amount int64) (logs []Log, err error) {
	trader, err := GetTrader(self, traderID)
	if err != nil {
		return
	}
	user, err := GetUserByID(trader.UserID)
	if err != nil {
		return
	}
	if user.ID != self.ID && user.Level >= self.Level {
		return
	}
	if amount < 1 {
		amount = 20
	} else if amount > 1000 {
		amount = 1000
	}
	err = DB.Where("trader_id = ?", traderID).Order("timestamp DESC").Limit(amount).Offset(page * amount).Find(&logs).Error
	return
}

// Log ...
func (l Logger) Log(method int, price, amount float64, messages ...interface{}) {
	go func() {
		message := ""
		for _, m := range messages {
			message += fmt.Sprintf("%+v", m)
		}
		log := Log{
			TraderID:     l.TraderID,
			Timestamp:    time.Now().Unix(),
			ExchangeType: l.ExchangeType,
			Type:         method,
			Price:        price,
			Amount:       amount,
			Message:      message,
		}
		DB.Create(&log)
	}()
}
