package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/miaolz123/conver"
	"github.com/miaolz123/samaritan/constant"
	"github.com/miaolz123/samaritan/model"
)

// Btcc : the exchange struct of btcc.com
type Btcc struct {
	stockMap     map[string]string
	orderTypeMap map[string]int
	periodMap    map[string]string
	minAmountMap map[string]float64
	records      map[string][]Record
	host         string
	logger       model.Logger
	option       Option

	simulate bool
	account  Account
	orders   map[string]Order

	limit     float64
	lastSleep int64
	lastTimes int64
}

// NewBtcc : create an exchange struct of okcoin.cn
func NewBtcc(opt Option) *Btcc {
	opt.MainStock = constant.BTC
	e := Btcc{
		stockMap: map[string]string{
			constant.BTC: "BTCCNY",
			constant.LTC: "BTCCNY",
			"LTCBTC":     "LTCBTC",
		},
		orderTypeMap: map[string]int{
			"bid": 1,
			"ask": -1,
		},
		periodMap: map[string]string{
			"M":   "1min",
			"M5":  "5min",
			"M15": "15min",
			"M30": "30min",
			"H":   "1hour",
			"D":   "1day",
			"W":   "1week",
		},
		minAmountMap: map[string]float64{
			constant.BTC: 0.001,
			constant.LTC: 0.01,
			"LTCBTC":     0.01,
		},
		records: make(map[string][]Record),
		host:    "https://api.btcc.com/api_trade_v1.php",
		logger:  model.Logger{TraderID: opt.TraderID, ExchangeType: opt.Type},
		option:  opt,

		limit:     5.0,
		lastSleep: time.Now().UnixNano(),
	}
	if _, ok := e.stockMap[e.option.MainStock]; !ok {
		e.option.MainStock = constant.BTC
	}
	return &e
}

// Log : print something to console
func (e *Btcc) Log(msgs ...interface{}) {
	e.logger.Log(constant.INFO, 0.0, 0.0, msgs...)
}

// GetType : get the type of this exchange
func (e *Btcc) GetType() string {
	return e.option.Type
}

// GetName : get the name of this exchange
func (e *Btcc) GetName() string {
	return e.option.Name
}

// GetMainStock : get the MainStock of this exchange
func (e *Btcc) GetMainStock() string {
	return e.option.MainStock
}

// SetMainStock : set the MainStock of this exchange
func (e *Btcc) SetMainStock(stock string) string {
	if _, ok := e.stockMap[stock]; ok {
		e.option.MainStock = stock
	}
	return e.option.MainStock
}

// SetLimit : set the limit calls amount per second of this exchange
func (e *Btcc) SetLimit(times interface{}) float64 {
	e.limit = conver.Float64Must(times)
	return e.limit
}

// AutoSleep : auto sleep to achieve the limit calls amount per second of this exchange
func (e *Btcc) AutoSleep() {
	now := time.Now().UnixNano()
	interval := 1e+9/e.limit*conver.Float64Must(e.lastTimes) - conver.Float64Must(now-e.lastSleep)
	if interval > 0.0 {
		time.Sleep(time.Duration(conver.Int64Must(interval)))
	}
	e.lastTimes = 0
	e.lastSleep = now
}

// GetMinAmount : get the min trade amonut of this exchange
func (e *Btcc) GetMinAmount(stock string) float64 {
	return e.minAmountMap[stock]
}

func (e *Btcc) getAuthJSON(method string, params ...interface{}) (jsoner *simplejson.Json, err error) {
	e.lastTimes++
	tonce := time.Now().UnixNano() / 1000
	param := ""
	for _, p := range params {
		if p != nil {
			param += fmt.Sprint(p, ",")
		} else {
			param += ","
		}
	}
	param = strings.TrimSuffix(param, ",")
	allParams := []string{
		fmt.Sprint("tonce=", tonce),
		"accesskey=" + e.option.AccessKey,
		"requestmethod=post",
		fmt.Sprint("id=", tonce),
		"method=" + method,
		"params=" + param,
	}
	postData := struct {
		ID     int64         `json:"id"`
		Method string        `json:"method"`
		Params []interface{} `json:"params"`
	}{
		ID:     tonce,
		Method: method,
		Params: params,
	}
	if len(postData.Params) == 0 {
		postData.Params = make([]interface{}, 0)
	}
	postDatas, err := json.Marshal(postData)
	if err != nil {
		return
	}
	req, err := http.NewRequest("POST", e.host, bytes.NewReader(postDatas))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json-rpc")
	req.Header.Set("Authorization", "Basic "+base64Encode(e.option.AccessKey+":"+signSha1(allParams, e.option.SecretKey)))
	req.Header.Set("Json-Rpc-Tonce", fmt.Sprint(tonce))
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return
	}
	data := []byte{}
	if resp == nil {
		err = fmt.Errorf("[POST %s] HTTP Error", method)
	} else if resp.StatusCode == 200 {
		data, _ = ioutil.ReadAll(resp.Body)
		resp.Body.Close()
	} else {
		data, _ = ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		err = fmt.Errorf("[POST %s] HTTP Status: %d, Info: %v", method, resp.StatusCode, string(data))
	}
	if err != nil {
		return
	}
	return simplejson.NewJson(data)
}

// Simulate : set the account of simulation
func (e *Btcc) Simulate(balance, btc, ltc interface{}) bool {
	e.simulate = true
	// e.orders = make(map[string]Order)
	e.account = Account{
		Balance: conver.Float64Must(balance),
		BTC:     conver.Float64Must(btc),
		LTC:     conver.Float64Must(ltc),
	}
	return true
}

// GetAccount : get the account detail of this exchange
func (e *Btcc) GetAccount() interface{} {
	if e.simulate {
		e.account.Total = e.account.Balance + e.account.FrozenBalance
		ticker, err := e.getTicker(constant.BTC, 10)
		if err != nil {
			e.logger.Log(constant.ERROR, 0.0, 0.0, "GetAccount() error, ", err)
			return false
		}
		e.account.Total += ticker.Mid * (e.account.BTC + e.account.FrozenBTC)
		ticker, err = e.getTicker(constant.LTC, 10)
		if err != nil {
			e.logger.Log(constant.ERROR, 0.0, 0.0, "GetAccount() error, ", err)
			return false
		}
		e.account.Total += ticker.Mid * (e.account.LTC + e.account.FrozenLTC)
		e.account.Net = e.account.Total
		if e.option.MainStock == constant.LTC {
			e.account.Stock = e.account.LTC
			e.account.FrozenStock = e.account.FrozenLTC
		} else {
			e.account.Stock = e.account.BTC
			e.account.FrozenStock = e.account.FrozenBTC
		}
		return e.account
	}
	json, err := e.getAuthJSON("getAccountInfo")
	if err != nil {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "GetAccount() error, ", err)
		return false
	}
	if errMsg := json.GetPath("error", "message").MustString(); errMsg != "" {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "GetAccount() error, ", errMsg)
		return false
	}
	return Account{
		Balance:       conver.Float64Must(json.GetPath("result", "balance", "cny", "amount").Interface()),
		FrozenBalance: conver.Float64Must(json.GetPath("result", "frozen", "cny", "amount").Interface()),
		BTC:           conver.Float64Must(json.GetPath("result", "balance", "btc", "amount").Interface()),
		FrozenBTC:     conver.Float64Must(json.GetPath("result", "frozen", "btc", "amount").Interface()),
		LTC:           conver.Float64Must(json.GetPath("result", "balance", "ltc", "amount").Interface()),
		FrozenLTC:     conver.Float64Must(json.GetPath("result", "frozen", "ltc", "amount").Interface()),
		Stock:         conver.Float64Must(json.GetPath("result", "balance", strings.ToLower(e.option.MainStock), "amount").Interface()),
		FrozenStock:   conver.Float64Must(json.GetPath("result", "frozen", strings.ToLower(e.option.MainStock), "amount").Interface()),
	}
}

// Buy : buy stocks
func (e *Btcc) Buy(stockType string, _price, _amount interface{}, msgs ...interface{}) interface{} {
	price := conver.Float64Must(_price)
	amount := conver.Float64Must(_amount)
	if _, ok := e.stockMap[stockType]; !ok {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "Buy() error, unrecognized stockType: ", stockType)
		return false
	}
	if e.simulate {
		ticker, err := e.getTicker(stockType, 10)
		if err != nil {
			e.logger.Log(constant.ERROR, 0.0, 0.0, "Buy() error, ", err)
			return false
		}
		total := simulateBuy(amount, ticker)
		if total > e.account.Balance {
			e.logger.Log(constant.ERROR, 0.0, 0.0, "Buy() error, balance is not enough")
			return false
		}
		e.account.Balance -= total
		if stockType == constant.LTC {
			e.account.LTC += amount
		} else {
			e.account.BTC += amount
		}
		e.logger.Log(constant.BUY, price, amount, msgs...)
		return fmt.Sprint(time.Now().Unix())
	}
	params := []interface{}{fmt.Sprintf("%f", price), fmt.Sprintf("%f", amount), e.stockMap[stockType]}
	if price <= 0 {
		ticker, err := e.getTicker(stockType, 5)
		if err != nil {
			e.logger.Log(constant.ERROR, 0.0, 0.0, "Buy() error, ", err)
			return false
		}
		precision := 0.01
		if e.minAmountMap[stockType] > 0 {
			precision = e.minAmountMap[stockType]
		}
		amountNew := math.Floor(amount/ticker.Sell/precision) * precision
		if amountNew <= precision {
			e.logger.Log(constant.ERROR, 0.0, 0.0, "Buy() error, amount less than min trade amount")
			return false
		}
		params[0] = nil
		params[1] = fmt.Sprintf("%f", amountNew)
	}
	json, err := e.getAuthJSON("buyOrder2", params...)
	if err != nil {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "Buy() error, ", err)
		return false
	}
	if errMsg := json.GetPath("error", "message").MustString(); errMsg != "" {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "Buy() error, ", errMsg)
		return false
	}
	e.logger.Log(constant.BUY, price, amount, msgs...)
	return fmt.Sprint(json.Get("result").Interface())
}

// Sell : sell stocks
func (e *Btcc) Sell(stockType string, _price, _amount interface{}, msgs ...interface{}) interface{} {
	price := conver.Float64Must(_price)
	amount := conver.Float64Must(_amount)
	if _, ok := e.stockMap[stockType]; !ok {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "Sell() error, unrecognized stockType: ", stockType)
		return false
	}
	if e.simulate {
		ticker, err := e.getTicker(stockType, 10)
		if err != nil {
			e.logger.Log(constant.ERROR, 0.0, 0.0, "Sell() error, ", err)
			return false
		}
		if price > ticker.Buy {
			e.logger.Log(constant.ERROR, 0.0, 0.0, "Sell() error, order price must be lesser than market buy price")
			return false
		}
		if stockType == constant.LTC {
			if amount > e.account.LTC {
				e.logger.Log(constant.ERROR, 0.0, 0.0, "Sell() error, stock is not enough")
				return false
			}
			e.account.LTC -= amount
		} else {
			if amount > e.account.BTC {
				e.logger.Log(constant.ERROR, 0.0, 0.0, "Sell() error, stock is not enough")
				return false
			}
			e.account.BTC -= amount
		}
		e.account.Balance += simulateSell(amount, ticker)
		e.logger.Log(constant.SELL, price, amount, msgs...)
		return fmt.Sprint(time.Now().Unix())
	}
	params := []interface{}{fmt.Sprintf("%f", price), fmt.Sprintf("%f", amount), e.stockMap[stockType]}
	if price <= 0 {
		params[0] = nil
	}
	json, err := e.getAuthJSON("sellOrder2", params...)
	if err != nil {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "Sell() error, ", err)
		return false
	}
	if errMsg := json.GetPath("error", "message").MustString(); errMsg != "" {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "Sell() error, ", errMsg)
		return false
	}
	e.logger.Log(constant.BUY, price, amount, msgs...)
	return fmt.Sprint(json.Get("result").Interface())
}

// GetOrder : get details of an order
func (e *Btcc) GetOrder(stockType, id string) interface{} {
	if _, ok := e.stockMap[stockType]; !ok {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "GetOrder() error, unrecognized stockType: ", stockType)
		return false
	}
	if e.simulate {
		return Order{ID: id, StockType: stockType}
	}
	json, err := e.getAuthJSON("getOrder", e.stockMap[stockType], conver.Int64Must(id))
	if err != nil {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "GetOrder() error, ", err)
		return false
	}
	if errMsg := json.GetPath("error", "message").MustString(); errMsg != "" {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "GetOrder() error, ", errMsg)
		return false
	}
	return Order{
		ID:         fmt.Sprint(json.GetPath("result", "order", "id").Interface()),
		Price:      conver.Float64Must(json.GetPath("result", "order", "price").Interface()),
		Amount:     conver.Float64Must(json.GetPath("result", "order", "amount_original").Interface()),
		DealAmount: conver.Float64Must(json.GetPath("result", "order", "amount").Interface()),
		OrderType:  e.orderTypeMap[json.GetPath("result", "order", "type").MustString()],
		StockType:  stockType,
	}
}

// GetOrders : get all unfilled orders
func (e *Btcc) GetOrders(stockType string) interface{} {
	orders := []Order{}
	if _, ok := e.stockMap[stockType]; !ok {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "GetOrders() error, unrecognized stockType: ", stockType)
		return false
	}
	if e.simulate {
		return orders
	}
	json, err := e.getAuthJSON("getOrders", true, e.stockMap[stockType])
	if err != nil {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "GetOrders() error, ", err)
		return false
	}
	if errMsg := json.GetPath("error", "message").MustString(); errMsg != "" {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "GetOrders() error, ", errMsg)
		return false
	}
	ordersJSON := json.GetPath("result", "order")
	count := len(ordersJSON.MustArray())
	for i := 0; i < count; i++ {
		orderJSON := ordersJSON.GetIndex(i)
		orders = append(orders, Order{
			ID:         fmt.Sprint(orderJSON.Get("id").Interface()),
			Price:      conver.Float64Must(orderJSON.Get("price").Interface()),
			Amount:     conver.Float64Must(orderJSON.Get("amount_original").Interface()),
			DealAmount: conver.Float64Must(orderJSON.Get("amount").Interface()),
			OrderType:  e.orderTypeMap[orderJSON.Get("type").MustString()],
			StockType:  stockType,
		})
	}
	return orders
}

// GetTrades : get all filled orders recently
func (e *Btcc) GetTrades(stockType string) interface{} {
	orders := []Order{}
	if _, ok := e.stockMap[stockType]; !ok {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "GetTrades() error, unrecognized stockType: ", stockType)
		return false
	}
	if e.simulate {
		return orders
	}
	json, err := e.getAuthJSON("getOrders", false, e.stockMap[stockType])
	if err != nil {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "GetTrades() error, ", err)
		return false
	}
	if errMsg := json.GetPath("error", "message").MustString(); errMsg != "" {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "GetTrades() error, ", errMsg)
		return false
	}
	ordersJSON := json.GetPath("result", "order")
	count := len(ordersJSON.MustArray())
	for i := 0; i < count; i++ {
		orderJSON := ordersJSON.GetIndex(i)
		order := Order{
			ID:         fmt.Sprint(orderJSON.Get("id").Interface()),
			Price:      conver.Float64Must(orderJSON.Get("price").Interface()),
			Amount:     conver.Float64Must(orderJSON.Get("amount_original").Interface()),
			DealAmount: conver.Float64Must(orderJSON.Get("amount").Interface()),
			OrderType:  e.orderTypeMap[orderJSON.Get("type").MustString()],
			StockType:  stockType,
		}
		if order.DealAmount == order.Amount {
			orders = append(orders, order)
		}
	}
	return orders
}

// CancelOrder : cancel an order
func (e *Btcc) CancelOrder(order Order) bool {
	if e.simulate {
		e.logger.Log(constant.CANCEL, order.Price, order.Amount-order.DealAmount, order)
		return true
	}
	json, err := e.getAuthJSON("cancelOrder", conver.Int64Must(order.ID), e.stockMap[order.StockType])
	if err != nil {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "CancelOrder() error, ", err)
		return false
	}
	if errMsg := json.GetPath("error", "message").MustString(); errMsg != "" {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "CancelOrder() error, ", errMsg)
		return false
	}
	if result := json.Get("result").MustBool(); !result {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "CancelOrder() error")
		return false
	}
	e.logger.Log(constant.CANCEL, order.Price, order.Amount-order.DealAmount, order)
	return true
}

// getTicker : get market ticker & depth
func (e *Btcc) getTicker(stockType string, sizes ...interface{}) (ticker Ticker, err error) {
	if _, ok := e.stockMap[stockType]; !ok {
		err = fmt.Errorf("GetTicker() error, unrecognized stockType: %+v", stockType)
		return
	}
	size := 20
	if len(sizes) > 0 && conver.IntMust(sizes[0]) > 0 {
		size = conver.IntMust(sizes[0])
	}
	resp, err := get(fmt.Sprintf("https://data.btcchina.com/data/orderbook?market=%v&limit=%v", strings.ToLower(e.stockMap[stockType]), size))
	if err != nil {
		err = fmt.Errorf("GetTicker() error, %+v", err)
		return
	}
	json, err := simplejson.NewJson(resp)
	if err != nil {
		err = fmt.Errorf("GetTicker() error, %+v", err)
		return
	}
	depthsJSON := json.Get("bids")
	for i := 0; i < len(depthsJSON.MustArray()); i++ {
		depthJSON := depthsJSON.GetIndex(i)
		ticker.Bids = append(ticker.Bids, OrderBook{
			Price:  depthJSON.GetIndex(0).MustFloat64(),
			Amount: depthJSON.GetIndex(1).MustFloat64(),
		})
	}
	depthsJSON = json.Get("asks")
	for i := len(depthsJSON.MustArray()); i > 0; i-- {
		depthJSON := depthsJSON.GetIndex(i - 1)
		ticker.Asks = append(ticker.Asks, OrderBook{
			Price:  depthJSON.GetIndex(0).MustFloat64(),
			Amount: depthJSON.GetIndex(1).MustFloat64(),
		})
	}
	if len(ticker.Bids) < 1 || len(ticker.Asks) < 1 {
		err = fmt.Errorf("GetTicker() error, can not get enough Bids or Asks")
		return
	}
	ticker.Buy = ticker.Bids[0].Price
	ticker.Sell = ticker.Asks[0].Price
	ticker.Mid = (ticker.Buy + ticker.Sell) / 2
	return
}

// GetTicker : get market ticker & depth
func (e *Btcc) GetTicker(stockType string, sizes ...interface{}) interface{} {
	ticker, err := e.getTicker(stockType, sizes...)
	if err != nil {
		e.logger.Log(constant.ERROR, 0.0, 0.0, err)
		return false
	}
	return ticker
}

// GetRecords : get candlestick data
func (e *Btcc) GetRecords(stockType, period string, sizes ...interface{}) interface{} {
	if _, ok := e.periodMap[period]; !ok {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "GetRecords() error, unrecognized period: ", period)
		return false
	}
	size := 200
	if len(sizes) > 0 && conver.IntMust(sizes[0]) > 0 {
		size = conver.IntMust(sizes[0])
	}
	records, err := getSosobtcRecords(e.records[period], e.option.Type, stockType, period, size)
	if err != nil {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "GetRecords() error, ", err)
		return false
	}
	e.records[period] = records
	return e.records[period]
}
