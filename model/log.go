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
	ID           int64   `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	TraderID     int64   `gorm:"index" json:"-"`
	Timestamp    int64   `json:"-"`
	ExchangeType string  `gorm:"type:varchar(50)" json:"exchangeType"`
	Type         string  `json:"type"` // [-1"error", 0"info", 1"profit", 2"buy", 3"sell", 4"cancel", 5"long", 6"short", 7"long_close", 8"short_close"]
	StockType    string  `gorm:"type:varchar(20)" json:"stockType"`
	Price        float64 `json:"price"`
	Amount       float64 `json:"amount"`
	Message      string  `gorm:"type:text" json:"message"`

	Time time.Time `gorm:"-" json:"time"`
}

// ListLog ...
func (user User) ListLog(id, size, page int64) (total int64, logs []Log, err error) {
	err = DB.Model(&Log{}).Where("trader_id = ?", id).Count(&total).Error
	if err != nil {
		return
	}
	err = DB.Where("trader_id = ?", id).Order("timestamp desc").Limit(size).Offset((page - 1) * size).Find(&logs).Error
	for i, l := range logs {
		logs[i].Time = time.Unix(l.Timestamp, 0)
	}
	return
}

// Logger struct
type Logger struct {
	TraderID     int64
	ExchangeType string
}

// Log ...
func (l Logger) Log(method string, stockType string, price, amount float64, messages ...interface{}) {
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
