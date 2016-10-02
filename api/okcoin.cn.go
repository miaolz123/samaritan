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

// OKCoinCn : the exchange struct of okcoin.cn
type OKCoinCn struct {
	stockMap     map[string]string
	orderTypeMap map[string]int
	periodMap    map[string]string
	records      map[string][]Record
	host         string
	logger       model.Logger
	option       Option

	simulate bool
	account  Account
	orders   map[string]Order
}

// NewOKCoinCn : create an exchange struct of okcoin.cn
func NewOKCoinCn(opt Option) *OKCoinCn {
	opt.MainStock = constant.BTC
	e := OKCoinCn{
		stockMap:     map[string]string{"BTC": "btc", "LTC": "ltc"},
		orderTypeMap: map[string]int{"buy": 1, "sell": -1, "buy_market": 2, "sell_market": -2},
		periodMap:    map[string]string{"M": "1min", "M5": "5min", "M15": "15min", "M30": "30min", "H": "1hour", "D": "1day", "W": "1week"},
		records:      make(map[string][]Record),
		host:         "https://www.okcoin.cn/api/v1/",
		logger:       model.Logger{TraderID: opt.TraderID, ExchangeType: opt.Type},
		option:       opt,
	}
	if _, ok := e.stockMap[e.option.MainStock]; !ok {
		e.option.MainStock = constant.BTC
	}
	return &e
}

// Log : print something to console
func (e *OKCoinCn) Log(msgs ...interface{}) {
	e.logger.Log(constant.INFO, 0.0, 0.0, msgs...)
}

// GetType : get the type of this exchange
func (e *OKCoinCn) GetType() string {
	return e.option.Type
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

func (e *OKCoinCn) getAuthJSON(url string, params []string) (json *simplejson.Json, err error) {
	params = append(params, "api_key="+e.option.AccessKey)
	sort.Strings(params)
	params = append(params, "secret_key="+e.option.SecretKey)
	params = append(params, "sign="+strings.ToUpper(signMd5(params)))
	resp, err := post(url, params)
	if err != nil {
		return
	}
	return simplejson.NewJson(resp)
}

// Simulate : set the account of simulation
func (e *OKCoinCn) Simulate(balance, btc, ltc interface{}) bool {
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
func (e *OKCoinCn) GetAccount() interface{} {
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
	json, err := e.getAuthJSON(e.host+"userinfo.do", []string{})
	if err != nil {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "GetAccount() error, ", err)
		return false
	}
	if result := json.Get("result").MustBool(); !result {
		err = fmt.Errorf("GetAccount() error, the error number is %v", json.Get("error_code").MustInt())
		e.logger.Log(constant.ERROR, 0.0, 0.0, "GetAccount() error, ", err)
		return false
	}
	return Account{
		Total:         conver.Float64Must(json.GetPath("info", "funds", "asset", "total").Interface()),
		Net:           conver.Float64Must(json.GetPath("info", "funds", "asset", "net").Interface()),
		Balance:       conver.Float64Must(json.GetPath("info", "funds", "free", "cny").Interface()),
		FrozenBalance: conver.Float64Must(json.GetPath("info", "funds", "freezed", "cny").Interface()),
		BTC:           conver.Float64Must(json.GetPath("info", "funds", "free", "btc").Interface()),
		FrozenBTC:     conver.Float64Must(json.GetPath("info", "funds", "freezed", "btc").Interface()),
		LTC:           conver.Float64Must(json.GetPath("info", "funds", "free", "ltc").Interface()),
		FrozenLTC:     conver.Float64Must(json.GetPath("info", "funds", "freezed", "ltc").Interface()),
		Stock:         conver.Float64Must(json.GetPath("info", "funds", "free", e.stockMap[e.option.MainStock]).Interface()),
		FrozenStock:   conver.Float64Must(json.GetPath("info", "funds", "freezed", e.stockMap[e.option.MainStock]).Interface()),
	}
}

// Buy : buy stocks
func (e *OKCoinCn) Buy(stockType string, _price, _amount interface{}, msgs ...interface{}) interface{} {
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
		if price < ticker.Sell {
			e.logger.Log(constant.ERROR, 0.0, 0.0, "Buy() error, order price must be greater than market sell price")
			return false
		}
		if price*amount > e.account.Balance {
			e.logger.Log(constant.ERROR, 0.0, 0.0, "Buy() error, balance is not enough")
			return false
		}
		e.account.Balance -= ticker.Sell * amount
		if stockType == constant.LTC {
			e.account.LTC += amount
		} else {
			e.account.BTC += amount
		}
		e.logger.Log(constant.BUY, price, amount, msgs...)
		return fmt.Sprint(time.Now().Unix())
	}
	params := []string{
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
		e.logger.Log(constant.ERROR, 0.0, 0.0, "Buy() error, ", err)
		return false
	}
	if result := json.Get("result").MustBool(); !result {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "Buy() error, the error number is ", json.Get("error_code").MustInt())
		return false
	}
	e.logger.Log(constant.BUY, price, amount, msgs...)
	return json.Get("order_id").MustString()
}

// Sell : sell stocks
func (e *OKCoinCn) Sell(stockType string, _price, _amount interface{}, msgs ...interface{}) interface{} {
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
		e.account.Balance += ticker.Buy * amount
		e.logger.Log(constant.SELL, price, amount, msgs...)
		return fmt.Sprint(time.Now().Unix())
	}
	params := []string{
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
		e.logger.Log(constant.ERROR, 0.0, 0.0, "Sell() error, ", err)
		return false
	}
	if result := json.Get("result").MustBool(); !result {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "Sell() error, the error number is ", json.Get("error_code").MustInt())
		return false
	}
	e.logger.Log(constant.SELL, price, amount, msgs...)
	return json.Get("order_id").MustString()
}

// GetOrder : get details of an order
func (e *OKCoinCn) GetOrder(stockType, id string) interface{} {
	if e.simulate {
		return Order{ID: id, StockType: stockType}
	}
	params := []string{
		"symbol=" + e.stockMap[stockType] + "_cny",
		"order_id=" + id,
	}
	json, err := e.getAuthJSON(e.host+"order_info.do", params)
	if err != nil {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "GetOrders() error, ", err)
		return false
	}
	if result := json.Get("result").MustBool(); !result {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "GetOrders() error, the error number is ", json.Get("error_code").MustInt())
		return false
	}
	ordersJSON := json.Get("orders")
	if len(ordersJSON.MustArray()) > 0 {
		orderJSON := ordersJSON.GetIndex(0)
		return Order{
			ID:         fmt.Sprint(orderJSON.Get("order_id").Interface()),
			Price:      orderJSON.Get("price").MustFloat64(),
			Amount:     orderJSON.Get("amount").MustFloat64(),
			DealAmount: orderJSON.Get("deal_amount").MustFloat64(),
			OrderType:  e.orderTypeMap[orderJSON.Get("type").MustString()],
			StockType:  stockType,
		}
	}
	return false
}

// CancelOrder : cancel an order
func (e *OKCoinCn) CancelOrder(order Order) bool {
	if e.simulate {
		e.logger.Log(constant.CANCEL, order.Price, order.Amount-order.DealAmount, order)
		return true
	}
	params := []string{
		"symbol=" + e.stockMap[order.StockType] + "_cny",
		"order_id=" + order.ID,
	}
	json, err := e.getAuthJSON(e.host+"cancel_order.do", params)
	if err != nil {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "CancelOrder() error, ", err)
		return false
	}
	if result := json.Get("result").MustBool(); !result {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "CancelOrder() error, the error number is ", json.Get("error_code").MustInt())
		return false
	}
	e.logger.Log(constant.CANCEL, order.Price, order.Amount-order.DealAmount, order)
	return true
}

// GetOrders : get all unfilled orders
func (e *OKCoinCn) GetOrders(stockType string) (orders []Order) {
	if _, ok := e.stockMap[stockType]; !ok {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "GetOrders() error, unrecognized stockType: ", stockType)
		return
	}
	if e.simulate {
		return
	}
	params := []string{
		"symbol=" + e.stockMap[stockType] + "_cny",
		"order_id=-1",
	}
	json, err := e.getAuthJSON(e.host+"order_info.do", params)
	if err != nil {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "GetOrders() error, ", err)
		return
	}
	if result := json.Get("result").MustBool(); !result {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "GetOrders() error, the error number is ", json.Get("error_code").MustInt())
		return
	}
	ordersJSON := json.Get("orders")
	count := len(ordersJSON.MustArray())
	for i := 0; i < count; i++ {
		orderJSON := ordersJSON.GetIndex(i)
		orders = append(orders, Order{
			ID:         fmt.Sprint(orderJSON.Get("order_id").Interface()),
			Price:      orderJSON.Get("price").MustFloat64(),
			Amount:     orderJSON.Get("amount").MustFloat64(),
			DealAmount: orderJSON.Get("deal_amount").MustFloat64(),
			OrderType:  e.orderTypeMap[orderJSON.Get("type").MustString()],
			StockType:  stockType,
		})
	}
	return orders
}

// GetTrades : get all filled orders recently
func (e *OKCoinCn) GetTrades(stockType string) (orders []Order) {
	if _, ok := e.stockMap[stockType]; !ok {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "GetTrades() error, unrecognized stockType: ", stockType)
		return
	}
	if e.simulate {
		return
	}
	params := []string{
		"symbol=" + e.stockMap[stockType] + "_cny",
		"status=1",
		"current_page=1",
		"page_length=200",
	}
	json, err := e.getAuthJSON(e.host+"order_history.do", params)
	if err != nil {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "GetTrades() error, ", err)
		return
	}
	if result := json.Get("result").MustBool(); !result {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "GetTrades() error, the error number is ", json.Get("error_code").MustInt())
		return
	}
	ordersJSON := json.Get("orders")
	count := len(ordersJSON.MustArray())
	for i := 0; i < count; i++ {
		orderJSON := ordersJSON.GetIndex(i)
		orders = append(orders, Order{
			ID:         fmt.Sprint(orderJSON.Get("order_id").Interface()),
			Price:      orderJSON.Get("price").MustFloat64(),
			Amount:     orderJSON.Get("amount").MustFloat64(),
			DealAmount: orderJSON.Get("deal_amount").MustFloat64(),
			OrderType:  e.orderTypeMap[orderJSON.Get("type").MustString()],
			StockType:  stockType,
		})
	}
	return orders
}

// getTicker : get market ticker & depth
func (e *OKCoinCn) getTicker(stockType string, sizes ...int) (ticker Ticker, err error) {
	if _, ok := e.stockMap[stockType]; !ok {
		err = fmt.Errorf("GetTicker() error, unrecognized stockType: %+v", stockType)
		return
	}
	size := 20
	if len(sizes) > 0 && sizes[0] > 20 {
		size = sizes[0]
	}
	resp, err := get(fmt.Sprint(e.host, "depth.do?symbol=", e.stockMap[stockType], "_cny&size=", size))
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
		ticker.Bids = append(ticker.Bids, MarketOrder{
			Price:  depthJSON.GetIndex(0).MustFloat64(),
			Amount: depthJSON.GetIndex(1).MustFloat64(),
		})
	}
	depthsJSON = json.Get("asks")
	for i := len(depthsJSON.MustArray()); i > 0; i-- {
		depthJSON := depthsJSON.GetIndex(i - 1)
		ticker.Asks = append(ticker.Asks, MarketOrder{
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
func (e *OKCoinCn) GetTicker(stockType string, sizes ...int) interface{} {
	ticker, err := e.getTicker(stockType, sizes...)
	if err != nil {
		e.logger.Log(constant.ERROR, 0.0, 0.0, err)
		return false
	}
	return ticker
}

// GetRecords : get candlestick data
func (e *OKCoinCn) GetRecords(stockType, period string, sizes ...int) (records []Record) {
	if _, ok := e.stockMap[stockType]; !ok {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "GetRecords() error, unrecognized stockType: ", stockType)
		return
	}
	if _, ok := e.periodMap[period]; !ok {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "GetRecords() error, unrecognized period: ", period)
		return
	}
	size := 200
	if len(sizes) > 0 {
		size = sizes[0]
	}
	resp, err := get(fmt.Sprint(e.host, "kline.do?symbol=", e.stockMap[stockType], "_cny&type=", e.periodMap[period], "&size=", size))
	if err != nil {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "GetRecords() error, ", err)
		return
	}
	json, err := simplejson.NewJson(resp)
	if err != nil {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "GetRecords() error, ", err)
		return
	}
	timeLast := int64(0)
	if len(e.records[period]) > 0 {
		timeLast = e.records[period][len(e.records[period])-1].Time
	}
	recordsNew := []Record{}
	for i := len(json.MustArray()); i > 0; i-- {
		recordJSON := json.GetIndex(i - 1)
		recordTime := conver.Int64Must(time.Unix(recordJSON.GetIndex(0).MustInt64()/1000, 0).Format("200601021504"))
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
