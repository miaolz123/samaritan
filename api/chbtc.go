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

// NewChbtc : create an exchange struct of chbtc.com
func NewChbtc(opt Option) *Chbtc {
	opt.MainStock = constant.BTC
	e := Chbtc{
		stockMap: map[string]string{
			constant.BTC: "btc_cny",
			constant.LTC: "ltc_cny",
			"ETH":        "eth_cny",
			"ETC":        "etc_cny",
			"ETH_BTC":    "eth_btc",
		},
		orderTypeMap: map[int]int{
			1: 1,
			0: -1,
		},
		periodMap: map[string]string{
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
			constant.BTC: 0.001,
			constant.LTC: 0.001,
			"ETH":        0.001,
			"ETC":        0.001,
			"ETH_BTC":    0.001,
		},
		records: make(map[string][]Record),
		host:    "https://trade.chbtc.com/api",
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
func (e *Chbtc) Log(msgs ...interface{}) {
	e.logger.Log(constant.INFO, 0.0, 0.0, msgs...)
}

// GetType : get the type of this exchange
func (e *Chbtc) GetType() string {
	return e.option.Type
}

// GetName : get the name of this exchange
func (e *Chbtc) GetName() string {
	return e.option.Name
}

// GetMainStock : get the MainStock of this exchange
func (e *Chbtc) GetMainStock() string {
	return e.option.MainStock
}

// SetMainStock : set the MainStock of this exchange
func (e *Chbtc) SetMainStock(stock string) string {
	if _, ok := e.stockMap[stock]; ok {
		e.option.MainStock = stock
	}
	return e.option.MainStock
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
func (e *Chbtc) Simulate(balance, btc, ltc interface{}) bool {
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
func (e *Chbtc) GetAccount() interface{} {
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
	json, err := e.getAuthJSON("getAccountInfo", []string{})
	if err != nil {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "GetAccount() error, ", err)
		return false
	}
	if code := json.Get("code").MustInt(); code != 0 && code > 1000 {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "GetAccount() error, ", json.Get("message").MustString())
		return false
	}
	mainStock := strings.Split(e.option.MainStock, "_")[0]
	account := map[string]float64{
		"Total":         0.0,
		"Net":           0.0,
		"BTC":           json.GetPath("result", "balance", "BTC", "amount").MustFloat64(),
		"FrozenBTC":     json.GetPath("result", "frozen", "BTC", "amount").MustFloat64(),
		"LTC":           json.GetPath("result", "balance", "LTC", "amount").MustFloat64(),
		"FrozenLTC":     json.GetPath("result", "frozen", "LTC", "amount").MustFloat64(),
		"ETH":           json.GetPath("result", "balance", "ETH", "amount").MustFloat64(),
		"FrozenETH":     json.GetPath("result", "frozen", "ETH", "amount").MustFloat64(),
		"ETC":           json.GetPath("result", "balance", "ETC", "amount").MustFloat64(),
		"FrozenETC":     json.GetPath("result", "frozen", "ETC", "amount").MustFloat64(),
		"Balance":       json.GetPath("result", "balance", "CNY", "amount").MustFloat64(),
		"FrozenBalance": json.GetPath("result", "frozen", "CNY", "amount").MustFloat64(),
		"Stock":         0.0,
		"FrozenStock":   0.0,
	}
	account["Stock"] = account[mainStock]
	account["FrozenStock"] = account["Frozen"+mainStock]
	return account
}

// Buy : buy stocks
func (e *Chbtc) Buy(stockType string, _price, _amount interface{}, msgs ...interface{}) interface{} {
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
		fmt.Sprintf("price=%f", price),
		fmt.Sprintf("amount=%f", amount),
		"tradeType=1",
		"currency=" + e.stockMap[stockType],
	}
	json, err := e.getAuthJSON("order", params)
	if err != nil {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "Buy() error, ", err)
		return false
	}
	if code := json.Get("code").MustInt(); code != 0 && code > 1000 {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "Buy() error, ", json.Get("message").MustString())
		return false
	}
	e.logger.Log(constant.BUY, price, amount, msgs...)
	return fmt.Sprint(json.Get("id").Interface())
}

// Sell : sell stocks
func (e *Chbtc) Sell(stockType string, _price, _amount interface{}, msgs ...interface{}) interface{} {
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
		fmt.Sprintf("price=%f", price),
		fmt.Sprintf("amount=%f", amount),
		"tradeType=0",
		"currency=" + e.stockMap[stockType],
	}
	json, err := e.getAuthJSON("order", params)
	if err != nil {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "Sell() error, ", err)
		return false
	}
	if code := json.Get("code").MustInt(); code != 0 && code > 1000 {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "Sell() error, ", json.Get("message").MustString())
		return false
	}
	e.logger.Log(constant.SELL, price, amount, msgs...)
	return fmt.Sprint(json.Get("id").Interface())
}

// GetOrder : get details of an order
func (e *Chbtc) GetOrder(stockType, id string) interface{} {
	if _, ok := e.stockMap[stockType]; !ok {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "GetOrder() error, unrecognized stockType: ", stockType)
		return false
	}
	if e.simulate {
		return Order{ID: id, StockType: stockType}
	}
	params := []string{
		"id=" + id,
		"currency=" + e.stockMap[stockType],
	}
	json, err := e.getAuthJSON("getOrder", params)
	if err != nil {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "GetOrder() error, ", err)
		return false
	}
	if code := json.Get("code").MustInt(); code != 0 && code > 1000 {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "GetOrder() error, ", json.Get("message").MustString())
		return false
	}
	return Order{
		ID:         fmt.Sprint(json.Get("id").Interface()),
		Price:      conver.Float64Must(json.Get("price").Interface()),
		Amount:     conver.Float64Must(json.Get("total_amount").Interface()),
		DealAmount: conver.Float64Must(json.Get("trade_amount").Interface()),
		OrderType:  e.orderTypeMap[json.Get("type").MustInt()],
		StockType:  stockType,
	}
}

// GetOrders : get all unfilled orders
func (e *Chbtc) GetOrders(stockType string) interface{} {
	orders := []Order{}
	if _, ok := e.stockMap[stockType]; !ok {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "GetOrders() error, unrecognized stockType: ", stockType)
		return false
	}
	if e.simulate {
		return orders
	}
	params := []string{
		"currency=" + e.stockMap[stockType],
		"pageIndex=1",
		"pageSize=100",
	}
	json, err := e.getAuthJSON("getUnfinishedOrdersIgnoreTradeType", params)
	if err != nil {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "GetOrders() error, ", err)
		return false
	}
	if code := json.Get("code").MustInt(); code != 0 && code != 3001 && code > 1000 {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "GetOrders() error, ", json.Get("message").MustString())
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
			OrderType:  e.orderTypeMap[orderJSON.Get("type").MustInt()],
			StockType:  stockType,
		})
	}
	return orders
}

// GetTrades : get all filled orders recently
func (e *Chbtc) GetTrades(stockType string) interface{} {
	orders := []Order{}
	if _, ok := e.stockMap[stockType]; !ok {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "GetTrades() error, unrecognized stockType: ", stockType)
		return false
	}
	if e.simulate {
		return orders
	}
	params := []string{
		"currency=" + e.stockMap[stockType],
		"pageIndex=1",
		"pageSize=100",
	}
	json, err := e.getAuthJSON("getOrdersIgnoreTradeType", params)
	if err != nil {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "GetTrades() error, ", err)
		return false
	}
	if code := json.Get("code").MustInt(); code != 0 && code != 3001 && code > 1000 {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "GetTrades() error, ", json.Get("message").MustString())
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
			OrderType:  e.orderTypeMap[orderJSON.Get("type").MustInt()],
			StockType:  stockType,
		})
	}
	return orders
}

// CancelOrder : cancel an order
func (e *Chbtc) CancelOrder(order Order) bool {
	if e.simulate {
		e.logger.Log(constant.CANCEL, order.Price, order.Amount-order.DealAmount, order)
		return true
	}
	params := []string{
		"id=" + order.ID,
		"currency=" + e.stockMap[order.StockType],
	}
	json, err := e.getAuthJSON("cancelOrder", params)
	if err != nil {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "CancelOrder() error, ", err)
		return false
	}
	if code := json.Get("code").MustInt(); code != 0 && code > 1000 {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "CancelOrder() error, ", json.Get("message").MustString())
		return false
	}
	e.logger.Log(constant.CANCEL, order.Price, order.Amount-order.DealAmount, order)
	return true
}

// getTicker : get market ticker & depth
func (e *Chbtc) getTicker(stockType string, sizes ...interface{}) (ticker Ticker, err error) {
	if _, ok := e.stockMap[stockType]; !ok {
		err = fmt.Errorf("GetTicker() error, unrecognized stockType: %+v", stockType)
		return
	}
	size := 20
	if len(sizes) > 0 && conver.IntMust(sizes[0]) > 0 {
		size = conver.IntMust(sizes[0])
	}
	resp, err := get(fmt.Sprintf("http://api.chbtc.com/data/v1/depth?currency=%v", e.stockMap[stockType]))
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
		e.logger.Log(constant.ERROR, 0.0, 0.0, err)
		return false
	}
	return ticker
}

// GetRecords : get candlestick data
func (e *Chbtc) GetRecords(stockType, period string, sizes ...interface{}) interface{} {
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
	resp, err := get(fmt.Sprintf("http://api.chbtc.com/data/v1/kline?currency=%v&type=%v&size=%v", e.stockMap[stockType], e.periodMap[period], size))
	if err != nil {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "GetRecords() error, ", err)
		return false
	}
	json, err := simplejson.NewJson(resp)
	if err != nil {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "GetRecords() error, ", err)
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
