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
	ID           int64 `gorm:"primary_key;AUTO_INCREMENT"`
	TraderID     int64 `gorm:"index"`
	Timestamp    int64
	ExchangeType string `gorm:"type:varchar(50)"`
	Type         int    // [-1"error", 0"info", 1"profit", 2"buy", 3"sell", 4"cancel", 5"long", 6"short", 7"long_close", 8"short_close"]
	StockType    string `gorm:"type:varchar(20)"`
	Price        float64
	Amount       float64
	Message      string `gorm:"type:text"`

	Time string `gorm:"-"`
}

// Logger struct
type Logger struct {
	TraderID     int64
	ExchangeType string
}

// Log ...
func (l Logger) Log(method int, stockType string, price, amount float64, messages ...interface{}) {
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
			StockType:    stockType,
			Price:        price,
			Amount:       amount,
			Message:      message,
		}
		DB.Create(&log)
	}()
}
