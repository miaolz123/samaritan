package api

import (
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

// Poloniex : the exchange struct of poloniex
type Poloniex struct {
	stockMap     map[string]string
	orderTypeMap map[string]int
	periodMap    map[string]string
	minAmountMap map[string]float64
	records      map[string][]Record
	host         string
	logger       model.Logger
	option       Option

	simulate bool
	account  map[string]float64
	orders   map[string]Order

	limit     float64
	lastSleep int64
	lastTimes int64
}

// NewPoloniex : create an exchange struct of poloniex
func NewPoloniex(opt Option) *Poloniex {
	opt.MainStock = "BTC_XMR"
	e := Poloniex{
		stockMap: map[string]string{
			"BTC_1CR":    "BTC_1CR",
			"BTC_BBR":    "BTC_BBR",
			"BTC_BCN":    "BTC_BCN",
			"BTC_BELA":   "BTC_BELA",
			"BTC_BITS":   "BTC_BITS",
			"BTC_BLK":    "BTC_BLK",
			"BTC_BLOCK":  "BTC_BLOCK",
			"BTC_BTCD":   "BTC_BTCD",
			"BTC_BTM":    "BTC_BTM",
			"BTC_BTS":    "BTC_BTS",
			"BTC_BURST":  "BTC_BURST",
			"BTC_C2":     "BTC_C2",
			"BTC_CGA":    "BTC_CGA",
			"BTC_CLAM":   "BTC_CLAM",
			"BTC_CURE":   "BTC_CURE",
			"BTC_DASH":   "BTC_DASH",
			"BTC_DGB":    "BTC_DGB",
			"BTC_DIEM":   "BTC_DIEM",
			"BTC_DOGE":   "BTC_DOGE",
			"BTC_EMC2":   "BTC_EMC2",
			"BTC_FLDC":   "BTC_FLDC",
			"BTC_FLO":    "BTC_FLO",
			"BTC_GEO":    "BTC_GEO",
			"BTC_GAME":   "BTC_GAME",
			"BTC_GRC":    "BTC_GRC",
			"BTC_HUC":    "BTC_HUC",
			"BTC_HZ":     "BTC_HZ",
			"BTC_LTBC":   "BTC_LTBC",
			"BTC_LTC":    "BTC_LTC",
			"BTC_MAID":   "BTC_MAID",
			"BTC_MMNXT":  "BTC_MMNXT",
			"BTC_OMNI":   "BTC_OMNI",
			"BTC_MYR":    "BTC_MYR",
			"BTC_NAUT":   "BTC_NAUT",
			"BTC_NAV":    "BTC_NAV",
			"BTC_NBT":    "BTC_NBT",
			"BTC_NEOS":   "BTC_NEOS",
			"BTC_NMC":    "BTC_NMC",
			"BTC_NOBL":   "BTC_NOBL",
			"BTC_NOTE":   "BTC_NOTE",
			"BTC_NSR":    "BTC_NSR",
			"BTC_NXT":    "BTC_NXT",
			"BTC_PINK":   "BTC_PINK",
			"BTC_POT":    "BTC_POT",
			"BTC_PPC":    "BTC_PPC",
			"BTC_QBK":    "BTC_QBK",
			"BTC_QORA":   "BTC_QORA",
			"BTC_QTL":    "BTC_QTL",
			"BTC_RBY":    "BTC_RBY",
			"BTC_RDD":    "BTC_RDD",
			"BTC_RIC":    "BTC_RIC",
			"BTC_SDC":    "BTC_SDC",
			"BTC_SJCX":   "BTC_SJCX",
			"BTC_STR":    "BTC_STR",
			"BTC_SYNC":   "BTC_SYNC",
			"BTC_SYS":    "BTC_SYS",
			"BTC_UNITY":  "BTC_UNITY",
			"BTC_VIA":    "BTC_VIA",
			"BTC_XVC":    "BTC_XVC",
			"BTC_VRC":    "BTC_VRC",
			"BTC_VTC":    "BTC_VTC",
			"BTC_XBC":    "BTC_XBC",
			"BTC_XCN":    "BTC_XCN",
			"BTC_XCP":    "BTC_XCP",
			"BTC_XDN":    "BTC_XDN",
			"BTC_XEM":    "BTC_XEM",
			"BTC_XMG":    "BTC_XMG",
			"BTC_XMR":    "BTC_XMR",
			"BTC_XPM":    "BTC_XPM",
			"BTC_XRP":    "BTC_XRP",
			"BTC_XST":    "BTC_XST",
			"USDT_BTC":   "USDT_BTC",
			"USDT_DASH":  "USDT_DASH",
			"USDT_LTC":   "USDT_LTC",
			"USDT_NXT":   "USDT_NXT",
			"USDT_STR":   "USDT_STR",
			"USDT_XMR":   "USDT_XMR",
			"USDT_XRP":   "USDT_XRP",
			"XMR_BBR":    "XMR_BBR",
			"XMR_BCN":    "XMR_BCN",
			"XMR_BLK":    "XMR_BLK",
			"XMR_BTCD":   "XMR_BTCD",
			"XMR_DASH":   "XMR_DASH",
			"XMR_DIEM":   "XMR_DIEM",
			"XMR_LTC":    "XMR_LTC",
			"XMR_MAID":   "XMR_MAID",
			"XMR_NXT":    "XMR_NXT",
			"XMR_QORA":   "XMR_QORA",
			"XMR_XDN":    "XMR_XDN",
			"BTC_IOC":    "BTC_IOC",
			"BTC_ETH":    "BTC_ETH",
			"USDT_ETH":   "USDT_ETH",
			"BTC_SC":     "BTC_SC",
			"BTC_BCY":    "BTC_BCY",
			"BTC_EXP":    "BTC_EXP",
			"BTC_FCT":    "BTC_FCT",
			"BTC_BITCNY": "BTC_BITCNY",
			"BTC_RADS":   "BTC_RADS",
			"BTC_AMP":    "BTC_AMP",
			"BTC_VOX":    "BTC_VOX",
			"BTC_DCR":    "BTC_DCR",
			"BTC_LSK":    "BTC_LSK",
			"ETH_LSK":    "ETH_LSK",
			"BTC_LBC":    "BTC_LBC",
			"BTC_STEEM":  "BTC_STEEM",
			"ETH_STEEM":  "ETH_STEEM",
			"BTC_SBD":    "BTC_SBD",
			"BTC_ETC":    "BTC_ETC",
			"ETH_ETC":    "ETH_ETC",
			"USDT_ETC":   "USDT_ETC",
			"BTC_REP":    "BTC_REP",
			"USDT_REP":   "USDT_REP",
			"ETH_REP":    "ETH_REP",
		},
		orderTypeMap: map[string]int{
			"buy":  1,
			"sell": -1,
		},
		periodMap: map[string]string{
			"M5":  "300",
			"M15": "900",
			"M30": "1800",
			"H2":  "7200",
			"H4":  "14400",
			"D":   "86400",
		},
		minAmountMap: map[string]float64{
			constant.BTC: 0.01,
			constant.LTC: 0.1,
		},
		records: make(map[string][]Record),
		host:    "https://poloniex.com/",
		logger:  model.Logger{TraderID: opt.TraderID, ExchangeType: opt.Type},
		option:  opt,

		account: make(map[string]float64),

		limit:     10.0,
		lastSleep: time.Now().UnixNano(),
	}
	if _, ok := e.stockMap[e.option.MainStock]; !ok {
		e.option.MainStock = "BTC_XMR"
	}
	return &e
}

// Log : print something to console
func (e *Poloniex) Log(msgs ...interface{}) {
	e.logger.Log(constant.INFO, 0.0, 0.0, msgs...)
}

// GetType : get the type of this exchange
func (e *Poloniex) GetType() string {
	return e.option.Type
}

// GetName : get the name of this exchange
func (e *Poloniex) GetName() string {
	return e.option.Name
}

// GetMainStock : get the MainStock of this exchange
func (e *Poloniex) GetMainStock() string {
	return e.option.MainStock
}

// SetMainStock : set the MainStock of this exchange
func (e *Poloniex) SetMainStock(stock string) string {
	if _, ok := e.stockMap[stock]; ok {
		e.option.MainStock = stock
	}
	return e.option.MainStock
}

// SetLimit : set the limit calls amount per second of this exchange
func (e *Poloniex) SetLimit(times interface{}) float64 {
	e.limit = conver.Float64Must(times)
	return e.limit
}

// AutoSleep : auto sleep to achieve the limit calls amount per second of this exchange
func (e *Poloniex) AutoSleep() {
	now := time.Now().UnixNano()
	interval := 1e+9/e.limit*conver.Float64Must(e.lastTimes) - conver.Float64Must(now-e.lastSleep)
	if interval > 0.0 {
		time.Sleep(time.Duration(conver.Int64Must(interval)))
	}
	e.lastTimes = 0
	e.lastSleep = now
}

// GetMinAmount : get the min trade amonut of this exchange
func (e *Poloniex) GetMinAmount(stock string) float64 {
	return e.minAmountMap[stock]
}

func (e *Poloniex) getAuthJSON(url string, params []string) (data []byte, json *simplejson.Json, err error) {
	e.lastTimes++
	params = append(params, fmt.Sprint("nonce=", time.Now().UnixNano()))
	req, err := http.NewRequest("POST", url, strings.NewReader(strings.Join(params, "&")))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Key", e.option.AccessKey)
	req.Header.Set("Sign", signSha512(params, e.option.SecretKey))
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return
	}
	if resp == nil {
		err = fmt.Errorf("[POST %s] HTTP Error Info: %v", url, err)
	} else if resp.StatusCode == 200 {
		data, _ = ioutil.ReadAll(resp.Body)
		resp.Body.Close()
	} else {
		err = fmt.Errorf("[POST %s] HTTP Status: %d, Info: %v", url, resp.StatusCode, err)
	}
	if err != nil {
		return
	}
	json, err = simplejson.NewJson(data)
	return
}

// Simulate : set the account of simulation
func (e *Poloniex) Simulate(acc map[string]interface{}) bool {
	e.simulate = true
	// e.orders = make(map[string]Order)
	for k, v := range acc {
		e.account[k] = conver.Float64Must(v)
	}
	return true
}

// GetAccount : get the account detail of this exchange
func (e *Poloniex) GetAccount() interface{} {
	if e.simulate {
		return e.account
	}
	data, jsoner, err := e.getAuthJSON(e.host+"tradingApi", []string{
		"command=returnCompleteBalances",
		"account=all",
	})
	if err != nil {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "GetAccount() error, ", err)
		return false
	}
	if errMsg := jsoner.Get("error").MustString(); errMsg != "" {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "GetAccount() error, ", errMsg)
		return false
	}
	resp := map[string]struct {
		Available string
		OnOrders  string
		BtcValue  string
	}{}
	if err = json.Unmarshal(data, &resp); err != nil {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "GetAccount() error, ", err)
		return false
	}
	account := map[string]float64{
		"Total":         0.0,
		"Net":           0.0,
		"Balance":       0.0,
		"FrozenBalance": 0.0,
		"Stock":         0.0,
		"FrozenStock":   0.0,
	}
	for k, v := range resp {
		account[k] = conver.Float64Must(v.Available)
		account["Frozen"+k] = conver.Float64Must(v.OnOrders)
		account["Total"] += conver.Float64Must(v.BtcValue)
	}
	account["Net"] = account["Total"]
	return account
}

// Buy : buy stocks
func (e *Poloniex) Buy(stockType string, _price, _amount interface{}, msgs ...interface{}) interface{} {
	price := conver.Float64Must(_price)
	amount := conver.Float64Must(_amount)
	if _, ok := e.stockMap[stockType]; !ok {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "Buy() error, unrecognized stockType: ", stockType)
		return false
	}
	if e.simulate {
		types := strings.Split(stockType, "_")
		if len(types) < 2 {
			e.logger.Log(constant.ERROR, 0.0, 0.0, "Buy() error, unrecognized stockType: ", stockType)
		}
		ticker, err := e.getTicker(stockType, 10)
		if err != nil {
			e.logger.Log(constant.ERROR, 0.0, 0.0, "Buy() error, ", err)
			return false
		}
		total := simulateBuy(amount, ticker)
		if total > e.account[types[0]] {
			e.logger.Log(constant.ERROR, 0.0, 0.0, "Buy() error, ", types[0], " is not enough")
			return false
		}
		e.account[types[0]] -= total
		e.account[types[1]] += amount
		e.logger.Log(constant.BUY, price, amount, msgs...)
		return fmt.Sprint(time.Now().Unix())
	}
	_, json, err := e.getAuthJSON(e.host+"tradingApi", []string{
		"command=buy",
		"currencyPair=" + stockType,
		fmt.Sprint("rate=", price),
		fmt.Sprint("amount=", amount),
	})
	if err != nil {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "Buy() error, ", err)
		return false
	}
	if errMsg := json.Get("error").MustString(); errMsg != "" {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "Buy() error, ", errMsg)
		return false
	}
	e.logger.Log(constant.BUY, price, amount, msgs...)
	return fmt.Sprint(json.Get("orderNumber").Interface())
}

// Sell : sell stocks
func (e *Poloniex) Sell(stockType string, _price, _amount interface{}, msgs ...interface{}) interface{} {
	price := conver.Float64Must(_price)
	amount := conver.Float64Must(_amount)
	if _, ok := e.stockMap[stockType]; !ok {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "Sell() error, unrecognized stockType: ", stockType)
		return false
	}
	if e.simulate {
		types := strings.Split(stockType, "_")
		if len(types) < 2 {
			e.logger.Log(constant.ERROR, 0.0, 0.0, "Buy() error, unrecognized stockType: ", stockType)
		}
		ticker, err := e.getTicker(stockType, 10)
		if err != nil {
			e.logger.Log(constant.ERROR, 0.0, 0.0, "Sell() error, ", err)
			return false
		}
		if price > ticker.Buy {
			e.logger.Log(constant.ERROR, 0.0, 0.0, "Sell() error, order price must be lesser than market buy price")
			return false
		}
		if amount > e.account[types[1]] {
			e.logger.Log(constant.ERROR, 0.0, 0.0, "Sell() error, ", e.account[types[1]], " is not enough")
			return false
		}
		e.account[types[1]] -= amount
		e.account[types[0]] += simulateSell(amount, ticker)
		e.logger.Log(constant.SELL, price, amount, msgs...)
		return fmt.Sprint(time.Now().Unix())
	}
	_, json, err := e.getAuthJSON(e.host+"tradingApi", []string{
		"command=sell",
		"currencyPair=" + stockType,
		fmt.Sprint("rate=", price),
		fmt.Sprint("amount=", amount),
	})
	if err != nil {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "Sell() error, ", err)
		return false
	}
	if errMsg := json.Get("error").MustString(); errMsg != "" {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "Sell() error, ", errMsg)
		return false
	}
	e.logger.Log(constant.BUY, price, amount, msgs...)
	return fmt.Sprint(json.Get("orderNumber").Interface())
}

// GetOrder : get details of an order
func (e *Poloniex) GetOrder(stockType, id string) interface{} {
	return Order{ID: id, StockType: stockType}
}

// GetOrders : get all unfilled orders
func (e *Poloniex) GetOrders(stockType string) interface{} {
	orders := []Order{}
	if _, ok := e.stockMap[stockType]; !ok {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "GetOrders() error, unrecognized stockType: ", stockType)
		return false
	}
	if e.simulate {
		return orders
	}
	_, json, err := e.getAuthJSON(e.host+"tradingApi", []string{
		"command=returnOpenOrders",
		"currencyPair=" + stockType,
	})
	if err != nil {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "GetOrders() error, ", err)
		return false
	}
	if errMsg := json.Get("error").MustString(); errMsg != "" {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "GetOrders() error, ", errMsg)
		return false
	}
	count := len(json.MustArray())
	for i := 0; i < count; i++ {
		orderJSON := json.GetIndex(i)
		orders = append(orders, Order{
			ID:         fmt.Sprint(orderJSON.Get("orderNumber").Interface()),
			Price:      conver.Float64Must(orderJSON.Get("rate").Interface()),
			Amount:     conver.Float64Must(orderJSON.Get("amount").Interface()),
			DealAmount: 0.0,
			OrderType:  e.orderTypeMap[orderJSON.Get("type").MustString()],
			StockType:  stockType,
		})
	}
	return orders
}

// GetTrades : get all filled orders recently
func (e *Poloniex) GetTrades(stockType string) interface{} {
	orders := []Order{}
	if _, ok := e.stockMap[stockType]; !ok {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "GetTrades() error, unrecognized stockType: ", stockType)
		return false
	}
	if e.simulate {
		return orders
	}
	_, json, err := e.getAuthJSON(e.host+"tradingApi", []string{
		"command=returnTradeHistory",
		"currencyPair=" + stockType,
	})
	if err != nil {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "GetTrades() error, ", err)
		return false
	}
	if errMsg := json.Get("error").MustString(); errMsg != "" {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "GetTrades() error, ", errMsg)
		return false
	}
	count := len(json.MustArray())
	for i := 0; i < count; i++ {
		orderJSON := json.GetIndex(i)
		orders = append(orders, Order{
			ID:         fmt.Sprint(orderJSON.Get("orderNumber").Interface()),
			Price:      conver.Float64Must(orderJSON.Get("rate").Interface()),
			Amount:     conver.Float64Must(orderJSON.Get("amount").Interface()),
			DealAmount: 0.0,
			OrderType:  e.orderTypeMap[orderJSON.Get("type").MustString()],
			StockType:  stockType,
		})
	}
	return orders
}

// CancelOrder : cancel an order
func (e *Poloniex) CancelOrder(order Order) bool {
	if e.simulate {
		e.logger.Log(constant.CANCEL, order.Price, order.Amount-order.DealAmount, order)
		return true
	}
	_, json, err := e.getAuthJSON(e.host+"tradingApi", []string{
		"command=cancelOrder",
		"orderNumber=" + order.ID,
	})
	if err != nil {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "CancelOrder() error, ", err)
		return false
	}
	if errMsg := json.Get("error").MustString(); errMsg != "" {
		e.logger.Log(constant.ERROR, 0.0, 0.0, "CancelOrder() error, ", errMsg)
		return false
	}
	e.logger.Log(constant.CANCEL, order.Price, order.Amount-order.DealAmount, order)
	return true
}

// getTicker : get market ticker & depth
func (e *Poloniex) getTicker(stockType string, sizes ...interface{}) (ticker Ticker, err error) {
	e.lastTimes++
	if _, ok := e.stockMap[stockType]; !ok {
		err = fmt.Errorf("GetTicker() error, unrecognized stockType: %+v", stockType)
		return
	}
	size := 20
	if len(sizes) > 0 && conver.IntMust(sizes[0]) > 0 {
		size = conver.IntMust(sizes[0])
	}
	resp, err := get(fmt.Sprintf("%vpublic?command=returnOrderBook&currencyPair=%v&depth=%v", e.host, e.stockMap[stockType], size))
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
			Price:  conver.Float64Must(depthJSON.GetIndex(0).Interface()),
			Amount: depthJSON.GetIndex(1).MustFloat64(),
		})
	}
	depthsJSON = json.Get("asks")
	for i := 0; i < len(depthsJSON.MustArray()); i++ {
		depthJSON := depthsJSON.GetIndex(i)
		ticker.Asks = append(ticker.Asks, OrderBook{
			Price:  conver.Float64Must(depthJSON.GetIndex(0).Interface()),
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
func (e *Poloniex) GetTicker(stockType string, sizes ...interface{}) interface{} {
	ticker, err := e.getTicker(stockType, sizes...)
	if err != nil {
		e.logger.Log(constant.ERROR, 0.0, 0.0, err)
		return false
	}
	return ticker
}

// GetRecords : get candlestick data
func (e *Poloniex) GetRecords(stockType, period string, sizes ...interface{}) interface{} {
	e.lastTimes++
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
	interval := conver.Int64Must(e.periodMap[period])
	start := time.Now().Unix() - interval*int64(size)
	if start < 0 {
		start = 0
	}
	resp, err := get(fmt.Sprintf("%vpublic?command=returnChartData&currencyPair=%v&start=%v&end=9999999999&period=%v", e.host, e.stockMap[stockType], start, e.periodMap[period]))
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
		recordTime := recordJSON.Get("date").MustInt64()
		if recordTime > timeLast {
			recordsNew = append([]Record{Record{
				Time:   recordTime,
				Open:   recordJSON.Get("open").MustFloat64(),
				High:   recordJSON.Get("high").MustFloat64(),
				Low:    recordJSON.Get("low").MustFloat64(),
				Close:  recordJSON.Get("close").MustFloat64(),
				Volume: recordJSON.Get("volume").MustFloat64(),
			}}, recordsNew...)
		} else if timeLast > 0 && recordTime == timeLast {
			e.records[period][len(e.records[period])-1] = Record{
				Time:   recordTime,
				Open:   recordJSON.Get("open").MustFloat64(),
				High:   recordJSON.Get("high").MustFloat64(),
				Low:    recordJSON.Get("low").MustFloat64(),
				Close:  recordJSON.Get("close").MustFloat64(),
				Volume: recordJSON.Get("volume").MustFloat64(),
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
