package api

import "github.com/robertkrimen/otto"

// Option : exchange option
type Option struct {
	TraderID  uint
	Type      string // one of ["okcoin.cn", "huobi"]
	Name      string
	AccessKey string
	SecretKey string
	MainStock string
	Ctx       *otto.Otto
}

// Exchange interface
type Exchange interface {
	Log(...interface{})
	GetType() string
	GetName() string
	GetMainStock() string
	SetMainStock(stock string) string
	SetLimit(times interface{}) float64
	AutoSleep()
	GetMinAmount(stock string) float64
	Simulate(balance, btc, ltc interface{}) bool
	GetAccount() interface{}
	Buy(stockType string, price, amount interface{}, msgs ...interface{}) interface{}
	Sell(stockType string, price, amount interface{}, msgs ...interface{}) interface{}
	GetOrder(stockType, id string) interface{}
	GetOrders(stockType string) []Order
	GetTrades(stockType string) []Order
	CancelOrder(order Order) bool
	GetTicker(stockType string, sizes ...interface{}) interface{}
	GetRecords(stockType, period string, sizes ...interface{}) []Record
}
