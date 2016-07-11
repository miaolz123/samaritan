package api

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/miaolz123/conver"
	"github.com/miaolz123/samaritan/log"
)

// Huobi : the exchange struct of okcoin.cn
type Huobi struct {
	stockMap     map[string]string
	orderTypeMap map[string]int
	periodMap    map[string]string
	records      map[string][]map[string]float64
	host         string
	log          log.Logger
	option       Option
}

// NewHuobi : create an exchange struct of okcoin.cn
func NewHuobi(opt Option) *Huobi {
	e := Huobi{
		stockMap:     map[string]string{"BTC": "btc", "LTC": "ltc"},
		orderTypeMap: map[string]int{"buy": 1, "sell": -1, "buy_market": 2, "sell_market": -2},
		periodMap:    map[string]string{"M": "1min", "M5": "5min", "M15": "15min", "M30": "30min", "H": "1hour", "D": "1day", "W": "1week"},
		records:      make(map[string][]map[string]float64),
		host:         "https://api.huobi.com/apiv3",
		log:          log.New(opt.Type),
		option:       opt,
	}
	if _, ok := e.stockMap[e.option.MainStock]; !ok {
		e.option.MainStock = "BTC"
	}
	return &e
}

// Log : print something to console
func (e *Huobi) Log(msgs ...interface{}) {
	e.log.Do("info", 0.0, 0.0, msgs...)
}

// GetMainStock : get the MainStock of this exchange
func (e *Huobi) GetMainStock() string {
	return e.option.MainStock
}

// SetMainStock : set the MainStock of this exchange
func (e *Huobi) SetMainStock(stock string) string {
	if _, ok := e.stockMap[stock]; ok {
		e.option.MainStock = stock
	}
	return e.option.MainStock
}

func (e *Huobi) getAuthJSON(url string, params []string) (json *simplejson.Json, err error) {
	sort.Strings(params)
	params = append(params, "secret_key="+e.option.SecretKey)
	params = append(params, "sign="+strings.ToUpper(signMd5(params)))
	resp, err := post(url, params)
	if err != nil {
		return
	}
	return simplejson.NewJson(resp)
}

// GetAccount : get the account detail of this exchange
func (e *Huobi) GetAccount() interface{} {
	account := make(map[string]float64)
	params := []string{
		"api_key=" + e.option.AccessKey,
	}
	json, err := e.getAuthJSON(e.host+"userinfo.do", params)
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

// Buy : buy stocks
func (e *Huobi) Buy(stockType string, price, amount float64, msgs ...interface{}) (id int) {
	if _, ok := e.stockMap[stockType]; !ok {
		e.log.Do("error", 0.0, 0.0, "Buy() error, unrecognized stockType: ", stockType)
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
	json, err := e.getAuthJSON(e.host+"trade.do", params)
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

// Sell : sell stocks
func (e *Huobi) Sell(stockType string, price, amount float64, msgs ...interface{}) (id int) {
	if _, ok := e.stockMap[stockType]; !ok {
		e.log.Do("error", 0.0, 0.0, "Sell() error, unrecognized stockType: ", stockType)
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
	json, err := e.getAuthJSON(e.host+"trade.do", params)
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

// GetOrder : get details of an order
func (e *Huobi) GetOrder(order map[string]interface{}) interface{} {
	params := []string{
		"api_key=" + e.option.AccessKey,
		"symbol=" + e.stockMap[fmt.Sprint(order["StockType"])] + "_cny",
		fmt.Sprint("order_id=", conver.IntMust(order["Id"])),
	}
	json, err := e.getAuthJSON(e.host+"order_info.do", params)
	if err != nil {
		e.log.Do("error", 0.0, 0.0, "GetOrders() error, ", err)
		return nil
	}
	if result := json.Get("result").MustBool(); !result {
		e.log.Do("error", 0.0, 0.0, "GetOrders() error, the error number is ", json.Get("error_code").MustInt())
		return nil
	}
	ordersJSON := json.Get("orders")
	if len(ordersJSON.MustArray()) > 0 {
		orderJSON := ordersJSON.GetIndex(0)
		return map[string]interface{}{
			"Id":         orderJSON.Get("order_id").MustInt(),
			"Price":      orderJSON.Get("price").MustFloat64(),
			"Amount":     orderJSON.Get("amount").MustFloat64(),
			"DealAmount": orderJSON.Get("deal_amount").MustFloat64(),
			"OrderType":  e.orderTypeMap[orderJSON.Get("type").MustString()],
			"StockType":  order["StockType"],
		}
	}
	return nil
}

// CancelOrder : cancel an order
func (e *Huobi) CancelOrder(order map[string]interface{}) bool {
	params := []string{
		"api_key=" + e.option.AccessKey,
		"symbol=" + e.stockMap[fmt.Sprint(order["StockType"])] + "_cny",
		fmt.Sprint("order_id=", conver.IntMust(order["Id"])),
	}
	json, err := e.getAuthJSON(e.host+"cancel_order.do", params)
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

// GetOrders : get all unfilled orders
func (e *Huobi) GetOrders(stockType string) (orders []map[string]interface{}) {
	if _, ok := e.stockMap[stockType]; !ok {
		e.log.Do("error", 0.0, 0.0, "GetOrders() error, unrecognized stockType: ", stockType)
		return
	}
	params := []string{
		"api_key=" + e.option.AccessKey,
		"symbol=" + e.stockMap[stockType] + "_cny",
		"order_id=-1",
	}
	json, err := e.getAuthJSON(e.host+"order_info.do", params)
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

// GetTrades : get all filled orders recently
func (e *Huobi) GetTrades(stockType string) (orders []map[string]interface{}) {
	if _, ok := e.stockMap[stockType]; !ok {
		e.log.Do("error", 0.0, 0.0, "GetTrades() error, unrecognized stockType: ", stockType)
		return
	}
	params := []string{
		"api_key=" + e.option.AccessKey,
		"symbol=" + e.stockMap[stockType] + "_cny",
		"status=1",
		"current_page=1",
		"page_length=200",
	}
	json, err := e.getAuthJSON(e.host+"order_history.do", params)
	if err != nil {
		e.log.Do("error", 0.0, 0.0, "GetTrades() error, ", err)
		return
	}
	if result := json.Get("result").MustBool(); !result {
		e.log.Do("error", 0.0, 0.0, "GetTrades() error, the error number is ", json.Get("error_code").MustInt())
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

// GetTicker : get market ticker & depth
func (e *Huobi) GetTicker(stockType string, sizes ...int) interface{} {
	if _, ok := e.stockMap[stockType]; !ok {
		e.log.Do("error", 0.0, 0.0, "GetTicker() error, unrecognized stockType: ", stockType)
		return nil
	}
	size := 20
	if len(sizes) > 0 && sizes[0] > 20 {
		size = sizes[0]
	}
	resp, err := get(fmt.Sprint(e.host, "depth.do?symbol=", e.stockMap[stockType], "_cny&size=", size))
	if err != nil {
		e.log.Do("error", 0.0, 0.0, "GetTicker() error, ", err)
		return nil
	}
	json, err := simplejson.NewJson(resp)
	if err != nil {
		e.log.Do("error", 0.0, 0.0, "GetTicker() error, ", err)
		return nil
	}
	bids := []map[string]float64{}
	depthsJSON := json.Get("bids")
	for i := 0; i < len(depthsJSON.MustArray()); i++ {
		depthJSON := depthsJSON.GetIndex(i)
		bids = append(bids, map[string]float64{
			"Price":  depthJSON.GetIndex(0).MustFloat64(),
			"Amount": depthJSON.GetIndex(1).MustFloat64(),
		})
	}
	asks := []map[string]float64{}
	depthsJSON = json.Get("asks")
	for i := len(depthsJSON.MustArray()); i > 0; i-- {
		depthJSON := depthsJSON.GetIndex(i - 1)
		asks = append(asks, map[string]float64{
			"Price":  depthJSON.GetIndex(0).MustFloat64(),
			"Amount": depthJSON.GetIndex(1).MustFloat64(),
		})
	}
	if len(bids) < 1 || len(asks) < 1 {
		e.log.Do("error", 0.0, 0.0, "GetTicker() error, can not get enough Bids or Asks")
		return nil
	}
	buy := bids[0]["Price"]
	sell := asks[0]["Price"]
	mid := (buy + sell) / 2
	return map[string]interface{}{
		"Mid":  mid,
		"Buy":  buy,
		"Sell": sell,
		"Bids": bids,
		"Asks": asks,
	}
}

// GetRecords : get candlestick data
func (e *Huobi) GetRecords(stockType, period string, sizes ...int) interface{} {
	if _, ok := e.stockMap[stockType]; !ok {
		e.log.Do("error", 0.0, 0.0, "GetRecords() error, unrecognized stockType: ", stockType)
		return nil
	}
	if _, ok := e.periodMap[period]; !ok {
		e.log.Do("error", 0.0, 0.0, "GetRecords() error, unrecognized period: ", period)
		return nil
	}
	size := 200
	if len(sizes) > 0 && sizes[0] > 200 {
		size = sizes[0]
	}
	resp, err := get(fmt.Sprint(e.host, "kline.do?symbol=", e.stockMap[stockType], "_cny&type=", e.periodMap[period], "&size=", size))
	if err != nil {
		e.log.Do("error", 0.0, 0.0, "GetRecords() error, ", err)
		return nil
	}
	json, err := simplejson.NewJson(resp)
	if err != nil {
		e.log.Do("error", 0.0, 0.0, "GetRecords() error, ", err)
		return nil
	}
	timeLast := 0.0
	if len(e.records[period]) > 0 {
		timeLast = e.records[period][len(e.records[period])-1]["Time"]
	}
	recordsNew := []map[string]float64{}
	for i := len(json.MustArray()); i > 0; i-- {
		recordJSON := json.GetIndex(i - 1)
		record := map[string]float64{
			"Time": conver.Float64Must(time.Unix(recordJSON.GetIndex(0).MustInt64()/1000, 0).Format("20060102150405")),
		}
		if record["Time"] > timeLast {
			record["Open"] = recordJSON.GetIndex(1).MustFloat64()
			record["High"] = recordJSON.GetIndex(2).MustFloat64()
			record["Low"] = recordJSON.GetIndex(3).MustFloat64()
			record["Close"] = recordJSON.GetIndex(4).MustFloat64()
			record["Volume"] = recordJSON.GetIndex(5).MustFloat64()
			recordsNew = append([]map[string]float64{record}, recordsNew...)
		} else if record["Time"] == timeLast {
			e.records[period][len(e.records[period])-1] = map[string]float64{
				"Time":   timeLast,
				"Open":   recordJSON.GetIndex(1).MustFloat64(),
				"High":   recordJSON.GetIndex(2).MustFloat64(),
				"Low":    recordJSON.GetIndex(3).MustFloat64(),
				"Close":  recordJSON.GetIndex(4).MustFloat64(),
				"Volume": recordJSON.GetIndex(5).MustFloat64(),
			}
		} else {
			break
		}
	}
	e.records[period] = append(e.records[period], recordsNew...)
	if len(e.records[period]) > size {
		e.records[period] = e.records[period][len(e.records[period])-size:]
	}
	return e.records[period]
}
