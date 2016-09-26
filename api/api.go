package api

import "github.com/robertkrimen/otto"

// Option : exchange option
type Option struct {
	TraderID  uint
	Type      string // one of ["okcoin.cn", "huobi"]
	AccessKey string
	SecretKey string
	MainStock string
	Ctx       *otto.Otto
}

// Exchange interface
type Exchange interface {
	Log(...interface{})
	GetMainStock() string
	SetMainStock(stock string) string
	Simulate(balance, btc, ltc interface{}) bool
	GetAccount() interface{}
	Buy(stockType string, price, amount interface{}, msgs ...interface{}) interface{}
	Sell(stockType string, price, amount interface{}, msgs ...interface{}) interface{}
	GetOrder(stockType, id string) interface{}
	CancelOrder(order Order) bool
	GetOrders(stockType string) []Order
	GetTrades(stockType string) []Order
	GetTicker(stockType string, sizes ...int) interface{}
	GetRecords(stockType, period string, sizes ...int) []Record
}
