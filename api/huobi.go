package api

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/miaolz123/conver"
	"github.com/miaolz123/samaritan/constant"
	"github.com/miaolz123/samaritan/model"
)

// Huobi : the exchange struct of okcoin.cn
type Huobi struct {
	stockMap     map[string]string
	orderTypeMap map[int]int
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

// NewHuobi : create an exchange struct of okcoin.cn
func NewHuobi(opt Option) *Huobi {
	opt.MainStock = constant.BTC
	e := Huobi{
		stockMap: map[string]string{
			constant.BTC: "1",
			constant.LTC: "2",
		},
		orderTypeMap: map[int]int{
			1: 1,
			2: -1,
			3: 2,
			4: -2,
		},
		periodMap: map[string]string{
			"M":   "001",
			"M5":  "005",
			"M15": "015",
			"M30": "030",
			"H":   "060",
			"D":   "100",
			"W":   "200",
		},
		minAmountMap: map[string]float64{
			constant.BTC: 0.001,
			constant.LTC: 0.01,
		},
		records: make(map[string][]Record),
		host:    "https://api.huobi.com/apiv3",
		logger:  model.Logger{TraderID: opt.TraderID, ExchangeType: opt.Type},
		option:  opt,

		limit:     10.0,
		lastSleep: time.Now().UnixNano(),
	}
	if _, ok := e.stockMap[e.option.MainStock]; !ok {
		e.option.MainStock = constant.BTC
	}
	return &e
}

// Log : print something to console
func (e *Huobi) Log(msgs ...interface{}) {
	e.logger.Log(constant.INFO, 0.0, 0.0, msgs...)
}

// GetType : get the type of this exchange
func (e *Huobi) GetType() string {
	return e.option.Type
}

// GetName : get the name of this exchange
func (e *Huobi) GetName() string {
	return e.option.Name
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

// SetLimit : set the limit calls amount per second of this exchange
func (e *Huobi) SetLimit(times interface{}) float64 {
	e.limit = conver.Float64Must(times)
	return e.limit
}

// AutoSleep : auto sleep to achieve the limit calls amount per second of this exchange
func (e *Huobi) AutoSleep() {
	now := time.Now().UnixNano()
	interval := 1e+9/e.limit*conver.Float64Must(e.lastTimes) - conver.Float64Must(now-e.lastSleep)
	if interval > 0.0 {
		time.Sleep(time.Duration(conver.Int64Must(interval)))
	}
	e.lastTimes = 0
	e.lastSleep = now
}

// GetMinAmount : get the min trade amonut of this exchange
func (e *Huobi) GetMinAmount(stock string) float64 {
	return e.minAmountMap[stock]
}

func (e *Huobi) getAuthJSON(url string, params []string, optionals ...string) (json *simplejson.Json, err error) {
	e.lastTimes++
	params = append(params, []string{
		"access_key=" + e.option.AccessKey,
		"secret_key=" + e.option.SecretKey,
		fmt.Sprint("created=", time.Now().Unix()),
	}...)
	sort.Strings(params)
	params = append(params, "sign="+signMd5(params))
	resp, err := post(url, append(params, optionals...))
	if err != nil {
		return
	}
	return simplejson.NewJson(resp)
}

// Simulate : set the account of simulation
func (e *Huobi) Simulate(balance, btc, ltc interface{}) bool {
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
func (e *Huobi) GetAccount() interface{} {
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
	params := []string{
		"method=get_account_info",
	}
	json, err := e.getAuthJSON(e.host, params, "market=cny")
	if err != nil {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "GetAccount() error, ", err)
		return false
	}
	if code := conver.IntMust(json.Get("code").Interface()); code > 0 {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "GetAccount() error, ", strings.TrimSpace(json.Get("msg").MustString()))
		return false
	}
	account := Account{
		Total:         conver.Float64Must(json.Get("total").Interface()),
		Net:           conver.Float64Must(json.Get("net_asset").Interface()),
		Balance:       conver.Float64Must(json.Get("available_cny_display").Interface()),
		FrozenBalance: conver.Float64Must(json.Get("frozen_cny_display").Interface()),
		BTC:           conver.Float64Must(json.Get("available_btc_display").Interface()),
		FrozenBTC:     conver.Float64Must(json.Get("frozen_btc_display").Interface()),
		LTC:           conver.Float64Must(json.Get("available_ltc_display").Interface()),
		FrozenLTC:     conver.Float64Must(json.Get("frozen_ltc_display").Interface()),
	}
	switch e.option.MainStock {
	case "BTC":
		account.Stock = account.BTC
		account.FrozenStock = account.FrozenBTC
	case "LTC":
		account.Stock = account.LTC
		account.FrozenStock = account.FrozenLTC
	}
	return account
}

// Buy : buy stocks
func (e *Huobi) Buy(stockType string, _price, _amount interface{}, msgs ...interface{}) interface{} {
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
	params := []string{
		"coin_type=" + e.stockMap[stockType],
		fmt.Sprint("amount=", amount),
	}
	methodParam := "method=buy_market"
	if price > 0 {
		methodParam = "method=buy"
		params = append(params, fmt.Sprint("price=", price))
	}
	params = append(params, methodParam)
	json, err := e.getAuthJSON(e.host, params, "market=cny")
	if err != nil {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "Buy() error, ", err)
		return false
	}
	if code := conver.IntMust(json.Get("code").Interface()); code > 0 {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "Buy() error, ", strings.TrimSpace(json.Get("msg").MustString()))
		return false
	}
	e.logger.Log(constant.BUY, price, amount, msgs...)
	return fmt.Sprint(json.Get("id").Interface())
}

// Sell : sell stocks
func (e *Huobi) Sell(stockType string, _price, _amount interface{}, msgs ...interface{}) interface{} {
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
	params := []string{
		"coin_type=" + e.stockMap[stockType],
		fmt.Sprint("amount=", amount),
	}
	methodParam := "method=sell_market"
	if price > 0 {
		methodParam = "method=sell"
		params = append(params, fmt.Sprint("price=", price))
	}
	params = append(params, methodParam)
	json, err := e.getAuthJSON(e.host, params, "market=cny")
	if err != nil {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "Sell() error, ", err)
		return false
	}
	if code := conver.IntMust(json.Get("code").Interface()); code > 0 {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "Sell() error, ", strings.TrimSpace(json.Get("msg").MustString()))
		return false
	}
	e.logger.Log(constant.SELL, price, amount, msgs...)
	return fmt.Sprint(json.Get("id").Interface())
}

// GetOrder : get details of an order
func (e *Huobi) GetOrder(stockType, id string) interface{} {
	if e.simulate {
		return Order{ID: id, StockType: stockType}
	}
	params := []string{
		"method=order_info",
		"coin_type=" + e.stockMap[stockType],
		"id=" + id,
	}
	json, err := e.getAuthJSON(e.host, params, "market=cny")
	if err != nil {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "GetOrder() error, ", err)
		return false
	}
	if code := conver.IntMust(json.Get("code").Interface()); code > 0 {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "GetOrder() error, ", strings.TrimSpace(json.Get("msg").MustString()))
		return false
	}
	return Order{
		ID:         fmt.Sprint(json.Get("id").Interface()),
		Price:      conver.Float64Must(json.Get("order_price").Interface()),
		Amount:     conver.Float64Must(json.Get("order_amount").Interface()),
		DealAmount: conver.Float64Must(json.Get("processed_amount").Interface()),
		OrderType:  e.orderTypeMap[json.Get("type").MustInt()],
		StockType:  stockType,
	}
}

// GetOrders : get all unfilled orders
func (e *Huobi) GetOrders(stockType string) interface{} {
	orders := []Order{}
	if _, ok := e.stockMap[stockType]; !ok {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "GetOrders() error, unrecognized stockType: ", stockType)
		return false
	}
	if e.simulate {
		return orders
	}
	params := []string{
		"method=get_orders",
		"coin_type=" + e.stockMap[stockType],
	}
	json, err := e.getAuthJSON(e.host+"order_info.do", params)
	if err != nil {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "GetOrders() error, ", err)
		return false
	}
	if code := conver.IntMust(json.Get("code").Interface()); code > 0 {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "GetOrders() error, ", strings.TrimSpace(json.Get("msg").MustString()))
		return false
	}
	count := len(json.MustArray())
	for i := 0; i < count; i++ {
		orderJSON := json.GetIndex(i)
		orders = append(orders, Order{
			ID:         fmt.Sprint(orderJSON.Get("id").Interface()),
			Price:      conver.Float64Must(orderJSON.Get("order_price").Interface()),
			Amount:     conver.Float64Must(orderJSON.Get("order_amount").Interface()),
			DealAmount: conver.Float64Must(orderJSON.Get("processed_amount").Interface()),
			OrderType:  e.orderTypeMap[orderJSON.Get("type").MustInt()],
			StockType:  stockType,
		})
	}
	return orders
}

// GetTrades : get all filled orders recently
func (e *Huobi) GetTrades(stockType string) interface{} {
	orders := []Order{}
	if _, ok := e.stockMap[stockType]; !ok {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "GetTrades() error, unrecognized stockType: ", stockType)
		return false
	}
	if e.simulate {
		return orders
	}
	params := []string{
		"method=get_new_deal_orders",
		"coin_type=" + e.stockMap[stockType],
	}
	json, err := e.getAuthJSON(e.host+"order_history.do", params)
	if err != nil {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "GetTrades() error, ", err)
		return false
	}
	if code := conver.IntMust(json.Get("code").Interface()); code > 0 {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "GetTrades() error, ", strings.TrimSpace(json.Get("msg").MustString()))
		return false
	}
	count := len(json.MustArray())
	for i := 0; i < count; i++ {
		orderJSON := json.GetIndex(i)
		orders = append(orders, Order{
			ID:         fmt.Sprint(orderJSON.Get("id").Interface()),
			Price:      conver.Float64Must(orderJSON.Get("order_price").Interface()),
			Amount:     conver.Float64Must(orderJSON.Get("order_amount").Interface()),
			DealAmount: conver.Float64Must(orderJSON.Get("processed_amount").Interface()),
			OrderType:  e.orderTypeMap[orderJSON.Get("type").MustInt()],
			StockType:  stockType,
		})
	}
	return orders
}

// CancelOrder : cancel an order
func (e *Huobi) CancelOrder(order Order) bool {
	if e.simulate {
		e.logger.Log(constant.CANCEL, order.Price, order.Amount-order.DealAmount, order)
		return true
	}
	params := []string{
		"method=cancel_order",
		"coin_type=" + e.stockMap[order.StockType],
		"id=" + order.ID,
	}
	json, err := e.getAuthJSON(e.host, params, "market=cny")
	if err != nil {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "CancelOrder() error, ", err)
		return false
	}
	if code := conver.IntMust(json.Get("code").Interface()); code > 0 {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "CancelOrder() error, ", strings.TrimSpace(json.Get("msg").MustString()))
		return false
	}
	if json.Get("result").MustString() == "success" {
		e.logger.Log(constant.CANCEL, order.Price, order.Amount-order.DealAmount, order)
		return true
	}
	e.logger.Log(constant.ERROR, 0.0, 0.0, "CancelOrder() error, ", json.Get("msg").Interface())
	return false
}

// getTicker : get market ticker & depth
func (e *Huobi) getTicker(stockType string, sizes ...interface{}) (ticker Ticker, err error) {
	if _, ok := e.stockMap[stockType]; !ok {
		err = fmt.Errorf("GetTicker() error, unrecognized stockType: %+v", stockType)
		return
	}
	size := 20
	if len(sizes) > 0 && conver.IntMust(sizes[0]) > 0 {
		size = conver.IntMust(sizes[0])
	}
	resp, err := get(fmt.Sprint("http://api.huobi.com/staticmarket/depth_", strings.ToLower(stockType), "_", size, ".js"))
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
	for i := 0; i < len(depthsJSON.MustArray()); i++ {
		depthJSON := depthsJSON.GetIndex(i)
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
func (e *Huobi) GetTicker(stockType string, sizes ...interface{}) interface{} {
	ticker, err := e.getTicker(stockType, sizes...)
	if err != nil {
		e.logger.Log(constant.ERROR, 0.0, 0.0, err)
		return false
	}
	return ticker
}

// GetRecords : get candlestick data
func (e *Huobi) GetRecords(stockType, period string, sizes ...interface{}) interface{} {
	if _, ok := e.stockMap[stockType]; !ok {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "GetRecords() error, unrecognized stockType: ", stockType)
		return false
	}
	if _, ok := e.periodMap[period]; !ok {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "GetRecords() error, unrecognized period: ", period)
		return false
	}
	size := 200
	if len(sizes) > 0 && conver.IntMust(sizes[0]) > 0 {
		size = conver.IntMust(sizes[0])
	}
	resp, err := get(fmt.Sprint("http://api.huobi.com/staticmarket/", strings.ToLower(stockType), "_kline_", e.periodMap[period], "_json.js"))
	if err != nil {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "GetRecords() error, ", err)
		return false
	}
	json, err := simplejson.NewJson(resp)
	if err != nil {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "GetRecords() error, ", err)
		return false
	}
	timeLast := int64(0)
	if len(e.records[period]) > 0 {
		timeLast = e.records[period][len(e.records[period])-1].Time
	}
	recordsNew := []Record{}
	for i := len(json.MustArray()); i > 0; i-- {
		recordJSON := json.GetIndex(i - 1)
		t, _ := time.Parse("20060102150405000", recordJSON.GetIndex(0).MustString("19700101000000000"))
		recordTime := t.Unix()
		if recordTime > timeLast {
			recordsNew = append([]Record{Record{
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
