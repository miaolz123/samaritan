package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/miaolz123/conver"
	"github.com/miaolz123/samaritan/constant"
	"github.com/miaolz123/samaritan/model"
)

// Btcc : the exchange struct of okcoin.cn
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
			constant.BTC: "btc",
			constant.LTC: "ltc",
		},
		orderTypeMap: map[string]int{
			"buy":         1,
			"sell":        -1,
			"buy_market":  2,
			"sell_market": -2,
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
			constant.BTC: 0.01,
			constant.LTC: 0.1,
		},
		records: make(map[string][]Record),
		host:    "https://api.btcc.com/api_trade_v1.php",
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

func (e *Btcc) getAuthJSON(method string, params []string) (jsoner *simplejson.Json, err error) {
	e.lastTimes++
	tonce := time.Now().UnixNano() / 1000000
	allParams := []string{
		fmt.Sprint("nonce=", tonce),
		"accesskey=" + e.option.AccessKey,
		"requestmethod=post",
		fmt.Sprint("id=", tonce),
		"method=" + method,
		"params=" + strings.Join(params, ","),
	}
	postData := struct {
		ID     int64         `json:"id"`
		Method string        `json:"method"`
		Params []interface{} `json:"params"`
	}{
		ID:     tonce,
		Method: method,
		Params: make([]interface{}, 0),
	}
	for _, p := range params {
		if p != "null" {
			postData.Params = append(postData.Params, p)
		} else {
			postData.Params = append(postData.Params, nil)
		}
	}
	postDatas, err := json.Marshal(postData)
	if err != nil {
		return
	}
	fmt.Println("POST DATA: ", string(postDatas))
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
	json, err := e.getAuthJSON("getAccountInfo", []string{})
	if err != nil {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "GetAccount() error, ", err)
		return false
	}
	fmt.Println(231231, json)
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
func (e *Btcc) GetOrder(stockType, id string) interface{} {
	if e.simulate {
		return Order{ID: id, StockType: stockType}
	}
	params := []string{
		"symbol=" + e.stockMap[stockType] + "_cny",
		"order_id=" + id,
	}
	json, err := e.getAuthJSON(e.host+"order_info.do", params)
	if err != nil {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "GetOrder() error, ", err)
		return false
	}
	if result := json.Get("result").MustBool(); !result {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "GetOrder() error, the error number is ", json.Get("error_code").MustInt())
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
	params := []string{
		"symbol=" + e.stockMap[stockType] + "_cny",
		"order_id=-1",
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
func (e *Btcc) GetTrades(stockType string) interface{} {
	orders := []Order{}
	if _, ok := e.stockMap[stockType]; !ok {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "GetTrades() error, unrecognized stockType: ", stockType)
		return false
	}
	if e.simulate {
		return orders
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
		return false
	}
	if result := json.Get("result").MustBool(); !result {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "GetTrades() error, the error number is ", json.Get("error_code").MustInt())
		return false
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

// CancelOrder : cancel an order
func (e *Btcc) CancelOrder(order Order) bool {
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
	resp, err := get(fmt.Sprintf("%vdepth.do?symbol=%v_cny&size=%v", e.host, e.stockMap[stockType], size))
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
	resp, err := get(fmt.Sprintf("%vkline.do?symbol=%v_cny&type=%v&size=%v", e.host, e.stockMap[stockType], e.periodMap[period], size))
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
