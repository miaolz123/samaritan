package api

import (
	"fmt"
	"strings"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/miaolz123/conver"
	"github.com/miaolz123/samaritan/constant"
	"github.com/miaolz123/samaritan/model"
)

// Chbtc : the exchange struct of chbtc.com
type Chbtc struct {
	stockTypeMap     map[string]string
	tradeTypeMap     map[int]string
	recordsPeriodMap map[string]string
	minAmountMap     map[string]float64
	records          map[string][]Record
	host             string
	logger           model.Logger
	option           Option

	simulate bool
	account  map[string]float64
	orders   map[string]Order

	limit     float64
	lastSleep int64
	lastTimes int64
}

func init() {
	constructor["chbtc"] = NewChbtc
}

// NewChbtc : create an exchange struct of chbtc.com
func NewChbtc(opt Option) Exchange {
	return &Chbtc{
		stockTypeMap: map[string]string{
			"BTC/CNY": "btc_cny",
			"LTC/CNY": "ltc_cny",
			"ETH/CNY": "eth_cny",
			"ETC/CNY": "etc_cny",
			"ETH/BTC": "eth_btc",
		},
		tradeTypeMap: map[int]string{
			1: constant.TradeTypeBuy,
			0: constant.TradeTypeSell,
		},
		recordsPeriodMap: map[string]string{
			"M":   "1min",
			"M3":  "3min",
			"M5":  "5min",
			"M15": "15min",
			"M30": "30min",
			"H":   "1hour",
			"H2":  "2hour",
			"H4":  "4hour",
			"H6":  "6hour",
			"H12": "12hour",
			"D":   "1day",
			"D3":  "3day",
			"W":   "1week",
		},
		minAmountMap: map[string]float64{
			"BTC/CNY": 0.001,
			"LTC/CNY": 0.001,
			"ETH/CNY": 0.001,
			"ETC/CNY": 0.001,
			"ETH/BTC": 0.001,
		},
		records: make(map[string][]Record),
		host:    "https://trade.chbtc.com/api",
		logger:  model.Logger{TraderID: opt.TraderID, ExchangeType: opt.Type},
		option:  opt,

		account: make(map[string]float64),

		limit:     10.0,
		lastSleep: time.Now().UnixNano(),
	}
}

// Log : print something to console
func (e *Chbtc) Log(msgs ...interface{}) {
	e.logger.Log(constant.INFO, "", 0.0, 0.0, msgs...)
}

// GetType : get the type of this exchange
func (e *Chbtc) GetType() string {
	return e.option.Type
}

// GetName : get the name of this exchange
func (e *Chbtc) GetName() string {
	return e.option.Name
}

// SetLimit : set the limit calls amount per second of this exchange
func (e *Chbtc) SetLimit(times interface{}) float64 {
	e.limit = conver.Float64Must(times)
	return e.limit
}

// AutoSleep : auto sleep to achieve the limit calls amount per second of this exchange
func (e *Chbtc) AutoSleep() {
	now := time.Now().UnixNano()
	interval := 1e+9/e.limit*conver.Float64Must(e.lastTimes) - conver.Float64Must(now-e.lastSleep)
	if interval > 0.0 {
		time.Sleep(time.Duration(conver.Int64Must(interval)))
	}
	e.lastTimes = 0
	e.lastSleep = now
}

// GetMinAmount : get the min trade amonut of this exchange
func (e *Chbtc) GetMinAmount(stock string) float64 {
	return e.minAmountMap[stock]
}

func (e *Chbtc) getAuthJSON(method string, params []string, optionals ...string) (json *simplejson.Json, err error) {
	e.lastTimes++
	params = append([]string{"method=" + method, "accesskey=" + e.option.AccessKey}, params...)
	params = append(params, "sign="+signChbtc(params, e.option.SecretKey), fmt.Sprint("reqTime=", time.Now().UnixNano()/1000000))
	resp, err := get(e.host + "/" + method + "?" + strings.Join(params, "&"))
	if err != nil {
		return
	}
	return simplejson.NewJson(resp)
}

// Simulate : set the account of simulation
func (e *Chbtc) Simulate(acc map[string]interface{}) bool {
	e.simulate = true
	// e.orders = make(map[string]Order)
	for k, v := range acc {
		e.account[k] = conver.Float64Must(v)
	}
	return true
}

// GetAccount : get the account detail of this exchange
func (e *Chbtc) GetAccount() interface{} {
	if e.simulate {
		return e.account
	}
	json, err := e.getAuthJSON("getAccountInfo", []string{})
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetAccount() error, ", err)
		return false
	}
	if code := json.Get("code").MustInt(); code != 0 && code > 1000 {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetAccount() error, ", json.Get("message").MustString())
		return false
	}
	return map[string]float64{
		"CNY":       json.GetPath("result", "balance", "CNY", "amount").MustFloat64(),
		"FrozenCNY": json.GetPath("result", "frozen", "CNY", "amount").MustFloat64(),
		"BTC":       json.GetPath("result", "balance", "BTC", "amount").MustFloat64(),
		"FrozenBTC": json.GetPath("result", "frozen", "BTC", "amount").MustFloat64(),
		"LTC":       json.GetPath("result", "balance", "LTC", "amount").MustFloat64(),
		"FrozenLTC": json.GetPath("result", "frozen", "LTC", "amount").MustFloat64(),
		"ETH":       json.GetPath("result", "balance", "ETH", "amount").MustFloat64(),
		"FrozenETH": json.GetPath("result", "frozen", "ETH", "amount").MustFloat64(),
		"ETC":       json.GetPath("result", "balance", "ETC", "amount").MustFloat64(),
		"FrozenETC": json.GetPath("result", "frozen", "ETC", "amount").MustFloat64(),
	}
}

// Trade : place an order
func (e *Chbtc) Trade(tradeType string, stockType string, _price, _amount interface{}, msgs ...interface{}) interface{} {
	stockType = strings.ToUpper(stockType)
	tradeType = strings.ToUpper(tradeType)
	price := conver.Float64Must(_price)
	amount := conver.Float64Must(_amount)
	if _, ok := e.stockTypeMap[stockType]; !ok {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "Trade() error, unrecognized stockType: ", stockType)
		return false
	}
	switch tradeType {
	case constant.TradeTypeBuy:
		return e.buy(stockType, price, amount, msgs...)
	case constant.TradeTypeSell:
		return e.sell(stockType, price, amount, msgs...)
	default:
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "Trade() error, unrecognized tradeType: ", tradeType)
		return false
	}
}

func (e *Chbtc) buy(stockType string, price, amount float64, msgs ...interface{}) interface{} {
	if e.simulate {
		currencies := strings.Split(stockType, "/")
		if len(currencies) < 2 {
			e.logger.Log(constant.ERROR, "", 0.0, 0.0, "Buy() error, unrecognized stockType: ", stockType)
			return false
		}
		ticker, err := e.getTicker(stockType, 10)
		if err != nil {
			e.logger.Log(constant.ERROR, "", 0.0, 0.0, "Buy() error, ", err)
			return false
		}
		total := simulateBuy(amount, ticker)
		if total > e.account[currencies[1]] {
			e.logger.Log(constant.ERROR, "", 0.0, 0.0, "Buy() error, ", currencies[1], " is not enough")
			return false
		}
		e.account[currencies[0]] += amount
		e.account[currencies[1]] -= total
		e.logger.Log(constant.BUY, stockType, price, amount, msgs...)
		return fmt.Sprint(time.Now().Unix())
	}
	params := []string{
		fmt.Sprintf("price=%f", price),
		fmt.Sprintf("amount=%f", amount),
		"tradeType=1",
		"currency=" + e.stockTypeMap[stockType],
	}
	json, err := e.getAuthJSON("order", params)
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "Buy() error, ", err)
		return false
	}
	if code := json.Get("code").MustInt(); code != 0 && code > 1000 {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "Buy() error, ", json.Get("message").MustString())
		return false
	}
	e.logger.Log(constant.BUY, stockType, price, amount, msgs...)
	return fmt.Sprint(json.Get("id").Interface())
}

func (e *Chbtc) sell(stockType string, price, amount float64, msgs ...interface{}) interface{} {
	if e.simulate {
		currencies := strings.Split(stockType, "/")
		if len(currencies) < 2 {
			e.logger.Log(constant.ERROR, "", 0.0, 0.0, "Sell() error, unrecognized stockType: ", stockType)
			return false
		}
		ticker, err := e.getTicker(stockType, 10)
		if err != nil {
			e.logger.Log(constant.ERROR, "", 0.0, 0.0, "Sell() error, ", err)
			return false
		}
		if amount > e.account[currencies[0]] {
			e.logger.Log(constant.ERROR, "", 0.0, 0.0, "Sell() error, ", currencies[0], " is not enough")
			return false
		}
		total := simulateSell(amount, ticker)
		e.account[currencies[0]] -= amount
		e.account[currencies[1]] += total
		e.logger.Log(constant.SELL, stockType, price, amount, msgs...)
		return fmt.Sprint(time.Now().Unix())
	}
	params := []string{
		fmt.Sprintf("price=%f", price),
		fmt.Sprintf("amount=%f", amount),
		"tradeType=0",
		"currency=" + e.stockTypeMap[stockType],
	}
	json, err := e.getAuthJSON("order", params)
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "Sell() error, ", err)
		return false
	}
	if code := json.Get("code").MustInt(); code != 0 && code > 1000 {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "Sell() error, ", json.Get("message").MustString())
		return false
	}
	e.logger.Log(constant.SELL, stockType, price, amount, msgs...)
	return fmt.Sprint(json.Get("id").Interface())
}

// GetOrder : get details of an order
func (e *Chbtc) GetOrder(stockType, id string) interface{} {
	stockType = strings.ToUpper(stockType)
	if _, ok := e.stockTypeMap[stockType]; !ok {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetOrder() error, unrecognized stockType: ", stockType)
		return false
	}
	if e.simulate {
		return Order{ID: id, StockType: stockType}
	}
	params := []string{
		"id=" + id,
		"currency=" + e.stockTypeMap[stockType],
	}
	json, err := e.getAuthJSON("getOrder", params)
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetOrder() error, ", err)
		return false
	}
	if code := json.Get("code").MustInt(); code != 0 && code > 1000 {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetOrder() error, ", json.Get("message").MustString())
		return false
	}
	return Order{
		ID:         fmt.Sprint(json.Get("id").Interface()),
		Price:      conver.Float64Must(json.Get("price").Interface()),
		Amount:     conver.Float64Must(json.Get("total_amount").Interface()),
		DealAmount: conver.Float64Must(json.Get("trade_amount").Interface()),
		TradeType:  e.tradeTypeMap[json.Get("type").MustInt()],
		StockType:  stockType,
	}
}

// GetOrders : get all unfilled orders
func (e *Chbtc) GetOrders(stockType string) interface{} {
	stockType = strings.ToUpper(stockType)
	orders := []Order{}
	if _, ok := e.stockTypeMap[stockType]; !ok {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetOrders() error, unrecognized stockType: ", stockType)
		return false
	}
	if e.simulate {
		return orders
	}
	params := []string{
		"currency=" + e.stockTypeMap[stockType],
		"pageIndex=1",
		"pageSize=100",
	}
	json, err := e.getAuthJSON("getUnfinishedOrdersIgnoreTradeType", params)
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetOrders() error, ", err)
		return false
	}
	if code := json.Get("code").MustInt(); code != 0 && code != 3001 && code > 1000 {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetOrders() error, ", json.Get("message").MustString())
		return false
	}
	count := len(json.MustArray())
	for i := 0; i < count; i++ {
		orderJSON := json.GetIndex(i)
		orders = append(orders, Order{
			ID:         fmt.Sprint(orderJSON.Get("id").Interface()),
			Price:      conver.Float64Must(orderJSON.Get("price").Interface()),
			Amount:     conver.Float64Must(orderJSON.Get("total_amount").Interface()),
			DealAmount: conver.Float64Must(orderJSON.Get("trade_amount").Interface()),
			TradeType:  e.tradeTypeMap[orderJSON.Get("type").MustInt()],
			StockType:  stockType,
		})
	}
	return orders
}

// GetTrades : get all filled orders recently
func (e *Chbtc) GetTrades(stockType string) interface{} {
	stockType = strings.ToUpper(stockType)
	orders := []Order{}
	if _, ok := e.stockTypeMap[stockType]; !ok {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetTrades() error, unrecognized stockType: ", stockType)
		return false
	}
	if e.simulate {
		return orders
	}
	params := []string{
		"currency=" + e.stockTypeMap[stockType],
		"pageIndex=1",
		"pageSize=100",
	}
	json, err := e.getAuthJSON("getOrdersIgnoreTradeType", params)
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetTrades() error, ", err)
		return false
	}
	if code := json.Get("code").MustInt(); code != 0 && code != 3001 && code > 1000 {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetTrades() error, ", json.Get("message").MustString())
		return false
	}
	count := len(json.MustArray())
	for i := 0; i < count; i++ {
		orderJSON := json.GetIndex(i)
		orders = append(orders, Order{
			ID:         fmt.Sprint(orderJSON.Get("id").Interface()),
			Price:      conver.Float64Must(orderJSON.Get("price").Interface()),
			Amount:     conver.Float64Must(orderJSON.Get("total_amount").Interface()),
			DealAmount: conver.Float64Must(orderJSON.Get("trade_amount").Interface()),
			TradeType:  e.tradeTypeMap[orderJSON.Get("type").MustInt()],
			StockType:  stockType,
		})
	}
	return orders
}

// CancelOrder : cancel an order
func (e *Chbtc) CancelOrder(order Order) bool {
	if e.simulate {
		e.logger.Log(constant.CANCEL, order.StockType, order.Price, order.Amount-order.DealAmount, order)
		return true
	}
	params := []string{
		"id=" + order.ID,
		"currency=" + e.stockTypeMap[order.StockType],
	}
	json, err := e.getAuthJSON("cancelOrder", params)
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "CancelOrder() error, ", err)
		return false
	}
	if code := json.Get("code").MustInt(); code != 0 && code > 1000 {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "CancelOrder() error, ", json.Get("message").MustString())
		return false
	}
	e.logger.Log(constant.CANCEL, order.StockType, order.Price, order.Amount-order.DealAmount, order)
	return true
}

// getTicker : get market ticker & depth
func (e *Chbtc) getTicker(stockType string, sizes ...interface{}) (ticker Ticker, err error) {
	stockType = strings.ToUpper(stockType)
	if _, ok := e.stockTypeMap[stockType]; !ok {
		err = fmt.Errorf("GetTicker() error, unrecognized stockType: %+v", stockType)
		return
	}
	size := 20
	if len(sizes) > 0 && conver.IntMust(sizes[0]) > 0 {
		size = conver.IntMust(sizes[0])
	}
	resp, err := get(fmt.Sprintf("http://api.chbtc.com/data/v1/depth?currency=%v", e.stockTypeMap[stockType]))
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
		if i >= size {
			break
		}
		depthJSON := depthsJSON.GetIndex(i)
		ticker.Bids = append(ticker.Bids, OrderBook{
			Price:  depthJSON.GetIndex(0).MustFloat64(),
			Amount: depthJSON.GetIndex(1).MustFloat64(),
		})
	}
	depthsJSON = json.Get("asks")
	length := len(depthsJSON.MustArray())
	for i := length; i > 0; i-- {
		if length-i >= size {
			break
		}
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
func (e *Chbtc) GetTicker(stockType string, sizes ...interface{}) interface{} {
	ticker, err := e.getTicker(stockType, sizes...)
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, err)
		return false
	}
	return ticker
}

// GetRecords : get candlestick data
func (e *Chbtc) GetRecords(stockType, period string, sizes ...interface{}) interface{} {
	stockType = strings.ToUpper(stockType)
	if _, ok := e.stockTypeMap[stockType]; !ok {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetRecords() error, unrecognized stockType: ", stockType)
		return false
	}
	if _, ok := e.recordsPeriodMap[period]; !ok {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetRecords() error, unrecognized period: ", period)
		return false
	}
	size := 200
	if len(sizes) > 0 && conver.IntMust(sizes[0]) > 0 {
		size = conver.IntMust(sizes[0])
	}
	resp, err := get(fmt.Sprintf("http://api.chbtc.com/data/v1/kline?currency=%v&type=%v&size=%v", e.stockTypeMap[stockType], e.recordsPeriodMap[period], size))
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetRecords() error, ", err)
		return false
	}
	json, err := simplejson.NewJson(resp)
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetRecords() error, ", err)
		return false
	}
	json = json.Get("data")
	timeLast := int64(0)
	if len(e.records[period]) > 0 {
		timeLast = e.records[period][len(e.records[period])-1].Time
	}
	recordsNew := []Record{}
	for i := len(json.MustArray()); i > 0; i-- {
		recordJSON := json.GetIndex(i - 1)
		recordTime := recordJSON.GetIndex(0).MustInt64() / 1000
		if recordTime > timeLast {
			recordsNew = append([]Record{{
				Time:   recordTime,
				Open:   recordJSON.GetIndex(1).MustFloat64(),
				High:   recordJSON.GetIndex(2).MustFloat64(),
				Low:    recordJSON.GetIndex(3).MustFloat64(),
				Close:  recordJSON.GetIndex(4).MustFloat64(),
				Volume: recordJSON.GetIndex(5).MustFloat64(),
			}}, recordsNew...)
		} else if timeLast > 0 && recordTime == timeLast {
			e.records[period][len(e.records[period])-1] = Record{
				Time:   recordTime,
				Open:   recordJSON.GetIndex(1).MustFloat64(),
				High:   recordJSON.GetIndex(2).MustFloat64(),
				Low:    recordJSON.GetIndex(3).MustFloat64(),
				Close:  recordJSON.GetIndex(4).MustFloat64(),
				Volume: recordJSON.GetIndex(5).MustFloat64(),
			}
		} else {
			break
		}
	}
	e.records[period] = append(e.records[period], recordsNew...)
	if len(e.records[period]) > size {
		e.records[period] = e.records[period][len(e.records[period])-size : len(e.records[period])]
	}
	return e.records[period]
}
