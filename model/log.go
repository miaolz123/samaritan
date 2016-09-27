package model

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/miaolz123/samaritan/constant"
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

// Log ...
func (l Logger) Log(method int, price, amount float64, messages ...interface{}) {
	go func() {
		message := ""
		for _, m := range messages {
			if method != constant.ERROR {
				v := reflect.ValueOf(m)
				switch v.Kind() {
				case reflect.Struct, reflect.Map, reflect.Slice:
					if bs, err := json.Marshal(m); err == nil {
						message += string(bs)
						continue
					}
				}
			}
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
