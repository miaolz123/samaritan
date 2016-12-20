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
	stockTypeMap     map[string]string
	tradeTypeMap     map[string]string
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

// NewPoloniex : create an exchange struct of poloniex
func NewPoloniex(opt Option) Exchange {
	return &Poloniex{
		stockTypeMap: map[string]string{
			"BTC/1CR":    "BTC_1CR",
			"BTC/BBR":    "BTC_BBR",
			"BTC/BCN":    "BTC_BCN",
			"BTC/BELA":   "BTC_BELA",
			"BTC/BITS":   "BTC_BITS",
			"BTC/BLK":    "BTC_BLK",
			"BTC/BLOCK":  "BTC_BLOCK",
			"BTC/BTCD":   "BTC_BTCD",
			"BTC/BTM":    "BTC_BTM",
			"BTC/BTS":    "BTC_BTS",
			"BTC/BURST":  "BTC_BURST",
			"BTC/C2":     "BTC_C2",
			"BTC/CGA":    "BTC_CGA",
			"BTC/CLAM":   "BTC_CLAM",
			"BTC/CURE":   "BTC_CURE",
			"BTC/DASH":   "BTC_DASH",
			"BTC/DGB":    "BTC_DGB",
			"BTC/DIEM":   "BTC_DIEM",
			"BTC/DOGE":   "BTC_DOGE",
			"BTC/EMC2":   "BTC_EMC2",
			"BTC/FLDC":   "BTC_FLDC",
			"BTC/FLO":    "BTC_FLO",
			"BTC/GEO":    "BTC_GEO",
			"BTC/GAME":   "BTC_GAME",
			"BTC/GRC":    "BTC_GRC",
			"BTC/HUC":    "BTC_HUC",
			"BTC/HZ":     "BTC_HZ",
			"BTC/LTBC":   "BTC_LTBC",
			"BTC/LTC":    "BTC_LTC",
			"BTC/MAID":   "BTC_MAID",
			"BTC/MMNXT":  "BTC_MMNXT",
			"BTC/OMNI":   "BTC_OMNI",
			"BTC/MYR":    "BTC_MYR",
			"BTC/NAUT":   "BTC_NAUT",
			"BTC/NAV":    "BTC_NAV",
			"BTC/NBT":    "BTC_NBT",
			"BTC/NEOS":   "BTC_NEOS",
			"BTC/NMC":    "BTC_NMC",
			"BTC/NOBL":   "BTC_NOBL",
			"BTC/NOTE":   "BTC_NOTE",
			"BTC/NSR":    "BTC_NSR",
			"BTC/NXT":    "BTC_NXT",
			"BTC/PINK":   "BTC_PINK",
			"BTC/POT":    "BTC_POT",
			"BTC/PPC":    "BTC_PPC",
			"BTC/QBK":    "BTC_QBK",
			"BTC/QORA":   "BTC_QORA",
			"BTC/QTL":    "BTC_QTL",
			"BTC/RBY":    "BTC_RBY",
			"BTC/RDD":    "BTC_RDD",
			"BTC/RIC":    "BTC_RIC",
			"BTC/SDC":    "BTC_SDC",
			"BTC/SJCX":   "BTC_SJCX",
			"BTC/STR":    "BTC_STR",
			"BTC/SYNC":   "BTC_SYNC",
			"BTC/SYS":    "BTC_SYS",
			"BTC/UNITY":  "BTC_UNITY",
			"BTC/VIA":    "BTC_VIA",
			"BTC/XVC":    "BTC_XVC",
			"BTC/VRC":    "BTC_VRC",
			"BTC/VTC":    "BTC_VTC",
			"BTC/XBC":    "BTC_XBC",
			"BTC/XCN":    "BTC_XCN",
			"BTC/XCP":    "BTC_XCP",
			"BTC/XDN":    "BTC_XDN",
			"BTC/XEM":    "BTC_XEM",
			"BTC/XMG":    "BTC_XMG",
			"BTC/XMR":    "BTC_XMR",
			"BTC/XPM":    "BTC_XPM",
			"BTC/XRP":    "BTC_XRP",
			"BTC/XST":    "BTC_XST",
			"USDT/BTC":   "USDT_BTC",
			"USDT/DASH":  "USDT_DASH",
			"USDT/LTC":   "USDT_LTC",
			"USDT/NXT":   "USDT_NXT",
			"USDT/STR":   "USDT_STR",
			"USDT/XMR":   "USDT_XMR",
			"USDT/XRP":   "USDT_XRP",
			"XMR/BBR":    "XMR_BBR",
			"XMR/BCN":    "XMR_BCN",
			"XMR/BLK":    "XMR_BLK",
			"XMR/BTCD":   "XMR_BTCD",
			"XMR/DASH":   "XMR_DASH",
			"XMR/DIEM":   "XMR_DIEM",
			"XMR/LTC":    "XMR_LTC",
			"XMR/MAID":   "XMR_MAID",
			"XMR/NXT":    "XMR_NXT",
			"XMR/QORA":   "XMR_QORA",
			"XMR/XDN":    "XMR_XDN",
			"BTC/IOC":    "BTC_IOC",
			"BTC/ETH":    "BTC_ETH",
			"USDT/ETH":   "USDT_ETH",
			"BTC/SC":     "BTC_SC",
			"BTC/BCY":    "BTC_BCY",
			"BTC/EXP":    "BTC_EXP",
			"BTC/FCT":    "BTC_FCT",
			"BTC/BITCNY": "BTC_BITCNY",
			"BTC/RADS":   "BTC_RADS",
			"BTC/AMP":    "BTC_AMP",
			"BTC/VOX":    "BTC_VOX",
			"BTC/DCR":    "BTC_DCR",
			"BTC/LSK":    "BTC_LSK",
			"ETH/LSK":    "ETH_LSK",
			"BTC/LBC":    "BTC_LBC",
			"BTC/STEEM":  "BTC_STEEM",
			"ETH/STEEM":  "ETH_STEEM",
			"BTC/SBD":    "BTC_SBD",
			"BTC/ETC":    "BTC_ETC",
			"ETH/ETC":    "ETH_ETC",
			"USDT/ETC":   "USDT_ETC",
			"BTC/REP":    "BTC_REP",
			"USDT/REP":   "USDT_REP",
			"ETH/REP":    "ETH_REP",
		},
		tradeTypeMap: map[string]string{
			"buy":  constant.TradeTypeBuy,
			"sell": constant.TradeTypeSell,
		},
		recordsPeriodMap: map[string]string{
			"M5":  "300",
			"M15": "900",
			"M30": "1800",
			"H2":  "7200",
			"H4":  "14400",
			"D":   "86400",
		},
		minAmountMap: map[string]float64{
			"BTC/XMR": 0.0,
		},
		records: make(map[string][]Record),
		host:    "https://poloniex.com/",
		logger:  model.Logger{TraderID: opt.TraderID, ExchangeType: opt.Type},
		option:  opt,

		account: make(map[string]float64),

		limit:     10.0,
		lastSleep: time.Now().UnixNano(),
	}
}

// Log : print something to console
func (e *Poloniex) Log(msgs ...interface{}) {
	e.logger.Log(constant.INFO, "", 0.0, 0.0, msgs...)
}

// GetType : get the type of this exchange
func (e *Poloniex) GetType() string {
	return e.option.Type
}

// GetName : get the name of this exchange
func (e *Poloniex) GetName() string {
	return e.option.Name
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
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetAccount() error, ", err)
		return false
	}
	if errMsg := jsoner.Get("error").MustString(); errMsg != "" {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetAccount() error, ", errMsg)
		return false
	}
	resp := map[string]struct {
		Available string
		OnOrders  string
		BtcValue  string
	}{}
	if err = json.Unmarshal(data, &resp); err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetAccount() error, ", err)
		return false
	}
	account := map[string]float64{}
	for k, v := range resp {
		account[k] = conver.Float64Must(v.Available)
		account["Frozen"+k] = conver.Float64Must(v.OnOrders)
	}
	return account
}

// Trade : place an order
func (e *Poloniex) Trade(tradeType string, stockType string, _price, _amount interface{}, msgs ...interface{}) interface{} {
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

func (e *Poloniex) buy(stockType string, price, amount float64, msgs ...interface{}) interface{} {
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
	_, json, err := e.getAuthJSON(e.host+"tradingApi", []string{
		"command=buy",
		"stockType=" + stockType,
		fmt.Sprintf("rate=%f", price),
		fmt.Sprintf("amount=%f", amount),
	})
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "Buy() error, ", err)
		return false
	}
	if errMsg := json.Get("error").MustString(); errMsg != "" {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "Buy() error, ", errMsg)
		return false
	}
	e.logger.Log(constant.BUY, stockType, price, amount, msgs...)
	return fmt.Sprint(json.Get("orderNumber").Interface())
}

func (e *Poloniex) sell(stockType string, price, amount float64, msgs ...interface{}) interface{} {
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
	_, json, err := e.getAuthJSON(e.host+"tradingApi", []string{
		"command=sell",
		"stockType=" + stockType,
		fmt.Sprintf("rate=%f", price),
		fmt.Sprintf("amount=%f", amount),
	})
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "Sell() error, ", err)
		return false
	}
	if errMsg := json.Get("error").MustString(); errMsg != "" {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "Sell() error, ", errMsg)
		return false
	}
	e.logger.Log(constant.BUY, stockType, price, amount, msgs...)
	return fmt.Sprint(json.Get("orderNumber").Interface())
}

// GetOrder : get details of an order
func (e *Poloniex) GetOrder(stockType, id string) interface{} {
	stockType = strings.ToUpper(stockType)
	if _, ok := e.stockTypeMap[stockType]; !ok {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetOrder() error, unrecognized stockType: ", stockType)
		return false
	}
	return Order{ID: id, StockType: stockType}
}

// GetOrders : get all unfilled orders
func (e *Poloniex) GetOrders(stockType string) interface{} {
	stockType = strings.ToUpper(stockType)
	orders := []Order{}
	if _, ok := e.stockTypeMap[stockType]; !ok {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetOrders() error, unrecognized stockType: ", stockType)
		return false
	}
	if e.simulate {
		return orders
	}
	_, json, err := e.getAuthJSON(e.host+"tradingApi", []string{
		"command=returnOpenOrders",
		"stockType=" + stockType,
	})
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetOrders() error, ", err)
		return false
	}
	if errMsg := json.Get("error").MustString(); errMsg != "" {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetOrders() error, ", errMsg)
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
			TradeType:  e.tradeTypeMap[orderJSON.Get("type").MustString()],
			StockType:  stockType,
		})
	}
	return orders
}

// GetTrades : get all filled orders recently
func (e *Poloniex) GetTrades(stockType string) interface{} {
	stockType = strings.ToUpper(stockType)
	orders := []Order{}
	if _, ok := e.stockTypeMap[stockType]; !ok {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetTrades() error, unrecognized stockType: ", stockType)
		return false
	}
	if e.simulate {
		return orders
	}
	_, json, err := e.getAuthJSON(e.host+"tradingApi", []string{
		"command=returnTradeHistory",
		"stockType=" + stockType,
	})
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetTrades() error, ", err)
		return false
	}
	if errMsg := json.Get("error").MustString(); errMsg != "" {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetTrades() error, ", errMsg)
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
			TradeType:  e.tradeTypeMap[orderJSON.Get("type").MustString()],
			StockType:  stockType,
		})
	}
	return orders
}

// CancelOrder : cancel an order
func (e *Poloniex) CancelOrder(order Order) bool {
	if e.simulate {
		e.logger.Log(constant.CANCEL, order.StockType, order.Price, order.Amount-order.DealAmount, order)
		return true
	}
	_, json, err := e.getAuthJSON(e.host+"tradingApi", []string{
		"command=cancelOrder",
		"orderNumber=" + order.ID,
	})
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "CancelOrder() error, ", err)
		return false
	}
	if errMsg := json.Get("error").MustString(); errMsg != "" {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "CancelOrder() error, ", errMsg)
		return false
	}
	e.logger.Log(constant.CANCEL, order.StockType, order.Price, order.Amount-order.DealAmount, order)
	return true
}

// getTicker : get market ticker & depth
func (e *Poloniex) getTicker(stockType string, sizes ...interface{}) (ticker Ticker, err error) {
	e.lastTimes++
	stockType = strings.ToUpper(stockType)
	if _, ok := e.stockTypeMap[stockType]; !ok {
		err = fmt.Errorf("GetTicker() error, unrecognized stockType: %+v", stockType)
		return
	}
	size := 20
	if len(sizes) > 0 && conver.IntMust(sizes[0]) > 0 {
		size = conver.IntMust(sizes[0])
	}
	resp, err := get(fmt.Sprintf("%vpublic?command=returnOrderBook&stockType=%v&depth=%v", e.host, e.stockTypeMap[stockType], size))
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
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, err)
		return false
	}
	return ticker
}

// GetRecords : get candlestick data
func (e *Poloniex) GetRecords(stockType, period string, sizes ...interface{}) interface{} {
	e.lastTimes++
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
	interval := conver.Int64Must(e.recordsPeriodMap[period])
	start := time.Now().Unix() - interval*int64(size)
	if start < 0 {
		start = 0
	}
	resp, err := get(fmt.Sprintf("%vpublic?command=returnChartData&stockType=%v&start=%v&end=9999999999&period=%v", e.host, e.stockTypeMap[stockType], start, e.recordsPeriodMap[period]))
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetRecords() error, ", err)
		return false
	}
	json, err := simplejson.NewJson(resp)
	if err != nil {
		e.logger.Log(constant.ERROR, "", 0.0, 0.0, "GetRecords() error, ", err)
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
			recordsNew = append([]Record{{
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
