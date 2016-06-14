package api

import (
	"fmt"
	"strings"

	"github.com/bitly/go-simplejson"
	"github.com/miaolz123/conver"
	"github.com/miaolz123/samaritan/log"
)

// OKCoinCn : the exchange struct of okcoin.cn
type OKCoinCn struct {
	stockMap     map[string]string
	orderTypeMap map[string]int
	host         string
	log          log.Logger
	option       Option
}

// NewOKCoinCn : create an exchange struct of okcoin.cn
func NewOKCoinCn(opt Option) *OKCoinCn {
	e := OKCoinCn{
		stockMap:     map[string]string{"BTC": "btc", "LTC": "ltc"},
		orderTypeMap: map[string]int{"buy": 1, "sell": -1, "buy_market": 2, "sell_market": -2},
		host:         "https://www.okcoin.cn/api/v1/",
		log:          log.New(opt.Type),
		option:       opt,
	}
	if _, ok := e.stockMap[e.option.MainStock]; !ok {
		e.option.MainStock = "BTC"
	}
	return &e
}

// Log : print something to console
func (e *OKCoinCn) Log(msgs ...interface{}) {
	e.log.Do("info", 0.0, 0.0, msgs...)
}

// GetMainStock : get the MainStock of this exchange
func (e *OKCoinCn) GetMainStock() string {
	return e.option.MainStock
}

// SetMainStock : set the MainStock of this exchange
func (e *OKCoinCn) SetMainStock(stock string) string {
	if _, ok := e.stockMap[stock]; ok {
		e.option.MainStock = stock
	}
	return e.option.MainStock
}

// GetAccount : get the account detail of this exchange
func (e *OKCoinCn) GetAccount() interface{} {
	account := make(map[string]float64)
	params := []string{
		"api_key=" + e.option.AccessKey,
		"secret_key=" + e.option.SecretKey,
	}
	params = append(params, "sign="+strings.ToUpper(signMd5(params)))
	resp, err := post(e.host+"userinfo.do", params)
	if err != nil {
		e.log.Do("error", 0.0, 0.0, "GetAccount() error, ", err)
		return nil
	}
	json, err := simplejson.NewJson(resp)
	if err != nil {
		e.log.Do("error", 0.0, 0.0, "GetAccount() error, ", err)
		return nil
	}

	if result := json.Get("result").MustBool(); !result {
		err = fmt.Errorf("GetAccount() error, the error number is %v", json.Get("error_code").MustInt())
		e.log.Do("error", 0.0, 0.0, "GetAccount() error, ", err)
		return nil
	}
	account["Total"] = conver.Float64Must(json.GetPath("info", "funds", "asset", "total").Interface())
	account["Net"] = conver.Float64Must(json.GetPath("info", "funds", "asset", "net").Interface())
	account["Balance"] = conver.Float64Must(json.GetPath("info", "funds", "free", "cny").Interface())
	account["FrozenBalance"] = conver.Float64Must(json.GetPath("info", "funds", "freezed", "cny").Interface())
	account["BTC"] = conver.Float64Must(json.GetPath("info", "funds", "free", "btc").Interface())
	account["FrozenBTC"] = conver.Float64Must(json.GetPath("info", "funds", "freezed", "btc").Interface())
	account["LTC"] = conver.Float64Must(json.GetPath("info", "funds", "free", "ltc").Interface())
	account["FrozenLTC"] = conver.Float64Must(json.GetPath("info", "funds", "freezed", "ltc").Interface())
	account["Stocks"] = account[e.option.MainStock]
	account["FrozenStocks"] = account["Frozen"+e.option.MainStock]
	return account
}

// Buy ...
func (e *OKCoinCn) Buy(stockType string, price, amount float64, msgs ...interface{}) (id int) {
	if _, ok := e.stockMap[stockType]; !ok {
		e.log.Do("error", 0.0, 0.0, "Buy() error, unrecognized stockType")
		return
	}
	params := []string{
		"api_key=" + e.option.AccessKey,
		"symbol=" + e.stockMap[stockType] + "_cny",
	}
	typeParam := "type=buy_market"
	amountParam := fmt.Sprint("price=", amount)
	if price > 0 {
		typeParam = "type=buy"
		amountParam = fmt.Sprint("amount=", amount)
		params = append(params, fmt.Sprint("price=", price))
	}
	params = append(params, typeParam, amountParam)
	params = append(params, "sign="+strings.ToUpper(signMd5(params)))
	resp, err := post(e.host+"trade.do", params)
	if err != nil {
		e.log.Do("error", 0.0, 0.0, "Buy() error, ", err)
		return
	}
	json, err := simplejson.NewJson(resp)
	if err != nil {
		e.log.Do("error", 0.0, 0.0, "Buy() error, ", err)
		return
	}
	if result := json.Get("result").MustBool(); !result {
		e.log.Do("error", 0.0, 0.0, "Buy() error, the error number is ", json.Get("error_code").MustInt())
		return
	}
	e.log.Do("buy", price, amount, msgs...)
	id = json.Get("order_id").MustInt()
	return
}

// Sell ...
func (e *OKCoinCn) Sell(stockType string, price, amount float64, msgs ...interface{}) (id int) {
	if _, ok := e.stockMap[stockType]; !ok {
		e.log.Do("error", 0.0, 0.0, "Sell() error, unrecognized stockType")
		return
	}
	params := []string{
		"api_key=" + e.option.AccessKey,
		"symbol=" + e.stockMap[stockType] + "_cny",
		fmt.Sprint("amount=", amount),
	}
	typeParam := "type=sell_market"
	if price > 0 {
		typeParam = "type=sell"
		params = append(params, fmt.Sprint("price=", price))
	}
	params = append(params, typeParam)
	params = append(params, "sign="+strings.ToUpper(signMd5(params)))
	resp, err := post(e.host+"trade.do", params)
	if err != nil {
		e.log.Do("error", 0.0, 0.0, "Sell() error, ", err)
		return
	}
	json, err := simplejson.NewJson(resp)
	if err != nil {
		e.log.Do("error", 0.0, 0.0, "Sell() error, ", err)
		return
	}
	if result := json.Get("result").MustBool(); !result {
		e.log.Do("error", 0.0, 0.0, "Sell() error, the error number is ", json.Get("error_code").MustInt())
		return
	}
	e.log.Do("sell", price, amount, msgs...)
	id = json.Get("order_id").MustInt()
	return
}

// CancelOrder ...
func (e *OKCoinCn) CancelOrder(order map[string]interface{}) bool {
	params := []string{
		"api_key=" + e.option.AccessKey,
		"symbol=" + e.stockMap[fmt.Sprint(order["StockType"])] + "_cny",
		fmt.Sprint("order_id=", conver.IntMust(order["Id"])),
	}
	params = append(params, "sign="+strings.ToUpper(signMd5(params)))
	resp, err := post(e.host+"cancel_order.do", params)
	if err != nil {
		e.log.Do("error", 0.0, 0.0, "CancelOrder() error, ", err)
		return false
	}
	json, err := simplejson.NewJson(resp)
	if err != nil {
		e.log.Do("error", 0.0, 0.0, "CancelOrder() error, ", err)
		return false
	}
	if result := json.Get("result").MustBool(); !result {
		e.log.Do("error", 0.0, 0.0, "CancelOrder() error, the error number is ", json.Get("error_code").MustInt())
		return false
	}
	e.log.Do("cancel", 0.0, 0.0, fmt.Sprintf("%v", order))
	return true
}

// GetOrders ...
func (e *OKCoinCn) GetOrders(stockType string) (orders []map[string]interface{}) {
	params := []string{
		"api_key=" + e.option.AccessKey,
		"symbol=" + e.stockMap[stockType] + "_cny",
		"status=0",
		"current_page=1",
		"page_length=200",
	}
	params = append(params, "sign="+strings.ToUpper(signMd5(params)))
	resp, err := post(e.host+"order_history.do", params)
	if err != nil {
		e.log.Do("error", 0.0, 0.0, "GetOrders() error, ", err)
		return
	}
	json, err := simplejson.NewJson(resp)
	if err != nil {
		e.log.Do("error", 0.0, 0.0, "GetOrders() error, ", err)
		return
	}
	if result := json.Get("result").MustBool(); !result {
		e.log.Do("error", 0.0, 0.0, "GetOrders() error, the error number is ", json.Get("error_code").MustInt())
		return
	}
	ordersJSON := json.Get("orders")
	count := len(ordersJSON.MustArray())
	for i := 0; i < count; i++ {
		orderJSON := ordersJSON.GetIndex(i)
		order := map[string]interface{}{
			"Id":         orderJSON.Get("order_id").MustInt(),
			"Price":      orderJSON.Get("price").MustFloat64(),
			"Amount":     orderJSON.Get("amount").MustFloat64(),
			"DealAmount": orderJSON.Get("deal_amount").MustFloat64(),
			"OrderType":  e.orderTypeMap[orderJSON.Get("type").MustString()],
			"StockType":  stockType,
		}
		orders = append(orders, order)
	}
	return orders
}

// GetOrder ...
func (e *OKCoinCn) GetOrder() map[string]interface{} {
	return map[string]interface{}{
		"Id":        123456789,
		"StockType": "BTC",
	}
}
