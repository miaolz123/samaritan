package api

import (
	"fmt"
	"strconv"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/cihub/seelog"
	"github.com/robertkrimen/otto"
)

type OKCoin struct {
	conf        exchangeConf
	Logger      seelog.LoggerInterface
	Records     map[string][]Record
	RecordsJser map[string][]map[string]interface{}
	TimeLast    float64
	Sleeper     float64
}

func NewOKCoin(a API, logger seelog.LoggerInterface) *OKCoin {
	okcoin := OKCoin{Detail: a, Logger: logger, TimeLast: float64(time.Now().UnixNano())}
	okcoin.Records = make(map[string][]Record)
	okcoin.RecordsJser = make(map[string][]map[string]interface{})
	return &okcoin
}

func (e *OKCoin) LogTrade(do string, msgs ...interface{}) {
	e.Logger.Criticalf("%-3s| %-6s| %s", do, e.Detail.Name, fmt.Sprint(msgs...))
	e.Logger.Flush()
}

func (e *OKCoin) LogError(msgs ...interface{}) {
	e.Logger.Errorf("错误 | %-6s| %s", e.Detail.Name, fmt.Sprint(msgs...))
	e.Logger.Flush()
}

func (e *OKCoin) Log(msgs ...interface{}) {
	e.Logger.Infof("信息 | %-6s| %s", e.Detail.Name, fmt.Sprint(msgs...))
	e.Logger.Flush()
}

func (e *OKCoin) SetSleeper() {
	e.Sleeper += e.TimeLast + 1000000000/e.Detail.Limiter - float64(time.Now().UnixNano())
}

func (e *OKCoin) ResetSleeper() {
	e.TimeLast = float64(time.Now().UnixNano())
	e.Sleeper = 0.0
}

func (e *OKCoin) GetSleeper() float64 {
	return e.Sleeper
}

func (e *OKCoin) GetJser() map[string]func(otto.FunctionCall) otto.Value {
	return map[string]func(otto.FunctionCall) otto.Value{
		"Log": func(call otto.FunctionCall) otto.Value {
			var msgs []interface{}
			for _, msg := range call.ArgumentList {
				m, _ := msg.Export()
				msgs = append(msgs, m)
			}
			e.Log(msgs...)
			return otto.TrueValue()
		},
		"GetAccount": func(call otto.FunctionCall) otto.Value {
			a, err := e.GetAccount()
			if err != nil {
				return otto.UndefinedValue()
			}
			acc, err := otto.New().ToValue(S2M(a))
			if err != nil {
				return otto.UndefinedValue()
			}
			return acc
		},
		"Buy": func(call otto.FunctionCall) otto.Value {
			stockType, err := call.Argument(0).ToString()
			if stockType == "undefined" {
				e.LogTrade("买入", "币种错误")
			}
			price, err := call.Argument(1).ToFloat()
			if err != nil {
				e.LogTrade("买入", "价钱错误")
			}
			amount, err := call.Argument(2).ToFloat()
			if amount <= 0 {
				e.LogTrade("买入", "数量错误")
			}
			var msgs []interface{}
			for i := 3; i < len(call.ArgumentList); i++ {
				m, _ := call.ArgumentList[i].Export()
				msgs = append(msgs, m)
			}
			ider, err := e.Buy(stockType, price, amount, msgs...)
			if err != nil {
				return otto.UndefinedValue()
			}
			id, err := otto.ToValue(ider)
			if err != nil {
				return otto.UndefinedValue()
			}
			return id
		},
		"Sell": func(call otto.FunctionCall) otto.Value {
			stockType, err := call.Argument(0).ToString()
			if stockType == "undefined" {
				e.LogTrade("卖出", "币种错误")
			}
			price, err := call.Argument(1).ToFloat()
			if err != nil {
				e.LogTrade("卖出", "价钱错误")
			}
			amount, err := call.Argument(2).ToFloat()
			if amount <= 0 {
				e.LogTrade("卖出", "数量错误")
			}
			var msgs []interface{}
			for i := 3; i < len(call.ArgumentList); i++ {
				m, _ := call.ArgumentList[i].Export()
				msgs = append(msgs, m)
			}
			ider, err := e.Sell(stockType, price, amount, msgs...)
			if err != nil {
				return otto.UndefinedValue()
			}
			id, err := otto.ToValue(ider)
			if err != nil {
				return otto.UndefinedValue()
			}
			return id
		},
		"CancelOrder": func(call otto.FunctionCall) otto.Value {
			orderObj := call.Argument(0).Object()
			keys := orderObj.Keys()
			order := make(map[string]interface{})
			for _, k := range keys {
				value, _ := orderObj.Get(k)
				order[k], _ = value.Export()
			}
			ret, err := otto.ToValue(e.CancelOrder(order))
			if err != nil {
				return otto.UndefinedValue()
			}
			return ret
		},
		"WithDraw": func(call otto.FunctionCall) otto.Value {
			stockType, _ := call.Argument(0).ToString()
			if stockType == "undefined" {
				e.LogTrade("提币", "币种错误")
			}
			address, _ := call.Argument(1).ToString()
			if address == "undefined" {
				e.LogTrade("提币", "地址错误")
			}
			amount, _ := call.Argument(2).ToFloat()
			if amount <= 0 {
				e.LogTrade("提币", "数量错误")
			}
			ider, err := e.WithDraw(stockType, address, amount)
			if err != nil {
				return otto.UndefinedValue()
			}
			id, err := otto.ToValue(ider)
			if err != nil {
				return otto.UndefinedValue()
			}
			return id
		},
		"GetOrder": func(call otto.FunctionCall) otto.Value {
			orderObj := call.Argument(0).Object()
			keys := orderObj.Keys()
			order := make(map[string]interface{})
			for _, k := range keys {
				value, _ := orderObj.Get(k)
				order[k], _ = value.Export()
			}
			retStr, err := e.GetOrder(order)
			if err != nil {
				return otto.UndefinedValue()
			}
			ret, err := otto.New().ToValue(S2M(retStr))
			if err != nil {
				return otto.UndefinedValue()
			}
			return ret
		},
		"GetOrders": func(call otto.FunctionCall) otto.Value {
			stockType, _ := call.Argument(0).ToString()
			if stockType == "undefined" {
				e.LogError(`"GetOrders()"的第一个参数错误`)
				return otto.UndefinedValue()
			}
			orders, err := e.GetOrders(stockType)
			if err != nil {
				return otto.UndefinedValue()
			}
			var ordersMap []map[string]interface{}
			for _, order := range orders {
				ordersMap = append(ordersMap, S2M(order))
			}
			rets, err := otto.New().ToValue(ordersMap)
			if err != nil {
				return otto.UndefinedValue()
			}
			return rets
		},
		"GetTrades": func(call otto.FunctionCall) otto.Value {
			stockType, _ := call.Argument(0).ToString()
			if stockType == "undefined" {
				e.LogError(`"GetTrades()"的第一个参数错误`)
				return otto.UndefinedValue()
			}
			orders, err := e.GetTrades(stockType)
			if err != nil {
				return otto.UndefinedValue()
			}
			var ordersMap []map[string]interface{}
			for _, order := range orders {
				ordersMap = append(ordersMap, S2M(order))
			}
			rets, err := otto.New().ToValue(ordersMap)
			if err != nil {
				return otto.UndefinedValue()
			}
			return rets
		},
		"GetTicker": func(call otto.FunctionCall) otto.Value {
			stockType, _ := call.Argument(0).ToString()
			if stockType == "undefined" {
				e.LogError(`"GetTicker()"的第一个参数错误`)
				return otto.UndefinedValue()
			}
			t, err := e.GetTickerJser(stockType)
			if err != nil {
				return otto.UndefinedValue()
			}
			ticker, err := otto.New().ToValue(t)
			if err != nil {
				return otto.UndefinedValue()
			}
			return ticker
		},
		"GetRecords": func(call otto.FunctionCall) otto.Value {
			stockType, _ := call.Argument(0).ToString()
			if stockType == "undefined" {
				e.LogError(`"GetRecords()"的第一个参数错误`)
				return otto.UndefinedValue()
			}
			period, _ := call.Argument(1).ToString()
			rs, err := e.GetRecordsJser(stockType, period)
			if err != nil {
				return otto.UndefinedValue()
			}
			records, err := otto.New().ToValue(rs)
			if err != nil {
				return otto.UndefinedValue()
			}
			return records
		},
	}
}

func (e *OKCoin) GetAccount() (acc Account, err error) {
	params := []string{
		"api_key=" + e.Detail.Access,
	}
	sign := e.sign(params)
	e.SetSleeper()
	get, err := httpPost(okcoinHost+"userinfo.do", append(params, "sign="+sign))
	if err != nil {
		e.LogError(err)
		return
	}
	js, err := simplejson.NewJson(get)
	if err != nil {
		e.LogError(err)
		return
	}
	result := js.Get("result").MustBool()
	if result {
		acc.Balance, _ = strconv.ParseFloat(js.GetPath("info", "funds", "asset", "total").MustString(), 64)
		acc.Net, _ = strconv.ParseFloat(js.GetPath("info", "funds", "asset", "net").MustString(), 64)
		acc.Money, _ = strconv.ParseFloat(js.GetPath("info", "funds", "free", "cny").MustString(), 64)
		acc.Btc, _ = strconv.ParseFloat(js.GetPath("info", "funds", "free", "btc").MustString(), 64)
		acc.Ltc, _ = strconv.ParseFloat(js.GetPath("info", "funds", "free", "ltc").MustString(), 64)
		acc.Stock = 0.0
	} else {
		err = fmt.Errorf("%s", fmt.Sprint("获取用户信息出错，错误代码：", js.Get("error_code").MustInt()))
		e.LogError(err)
	}
	return
}

func (e *OKCoin) Buy(stockType string, price, amount float64, msg ...interface{}) (id int, err error) {
	params := []string{
		"api_key=" + e.Detail.Access,
		"symbol=" + okcoinStockType[stockType],
	}
	typeParam := "type=buy_market"
	amountParam := fmt.Sprint("price=", amount)
	if price > 0 {
		typeParam = "type=buy"
		amountParam = fmt.Sprint("amount=", amount)
		params = append(params, fmt.Sprint("price=", price))
	}
	params = append(params, typeParam, amountParam)
	sign := e.sign(params)
	e.SetSleeper()
	get, err := httpPost(okcoinHost+"trade.do", append(params, "sign="+sign))
	if err != nil {
		e.LogError(err)
		return
	}
	js, err := simplejson.NewJson(get)
	if err != nil {
		e.LogError(err)
		return
	}
	result := js.Get("result").MustBool()
	if result {
		logMsg := ""
		if len(msg) > 0 {
			logMsg = "，" + fmt.Sprint(msg...)
		}
		e.LogTrade("买入", stockType, "价钱：", price, "，数量：", amount, logMsg)
		id = js.Get("order_id").MustInt()
	} else {
		err = fmt.Errorf("%s", fmt.Sprint(stockType, "买入下单失败，错误代码：", js.Get("error_code").MustInt()))
		e.LogError(err)
	}
	return
}

func (e *OKCoin) Sell(stockType string, price, amount float64, msg ...interface{}) (id int, err error) {
	params := []string{
		"api_key=" + e.Detail.Access,
		"symbol=" + okcoinStockType[stockType],
		fmt.Sprint("amount=", amount),
	}
	typeParam := "type=sell_market"
	if price > 0 {
		typeParam = "type=sell"
		params = append(params, fmt.Sprint("price=", price))
	}
	params = append(params, typeParam)
	sign := e.sign(params)
	e.SetSleeper()
	get, err := httpPost(okcoinHost+"trade.do", append(params, "sign="+sign))
	if err != nil {
		e.LogError(err)
		return
	}
	js, err := simplejson.NewJson(get)
	if err != nil {
		e.LogError(err)
		return
	}
	result := js.Get("result").MustBool()
	if result {
		logMsg := ""
		if len(msg) > 0 {
			logMsg = "，" + fmt.Sprint(msg...)
		}
		e.LogTrade("卖出", stockType, "价钱：", price, "，数量：", amount, "，", logMsg)
		id = js.Get("order_id").MustInt()
	} else {
		err = fmt.Errorf("%s", fmt.Sprint(stockType, "卖出下单失败，错误代码：", js.Get("error_code").MustInt()))
		e.LogError(err)
	}
	return
}

func (e *OKCoin) CancelOrder(order map[string]interface{}) (result bool) {
	params := []string{
		"api_key=" + e.Detail.Access,
		"symbol=" + okcoinStockType[fmt.Sprint(order["StockType"])],
		fmt.Sprint("order_id=", order["Id"]),
	}
	sign := e.sign(params)
	e.SetSleeper()
	get, err := httpPost(okcoinHost+"cancel_order.do", append(params, "sign="+sign))
	if err != nil {
		e.LogError(err)
		return
	}
	js, err := simplejson.NewJson(get)
	if err != nil {
		e.LogError(err)
		return
	}
	result = js.Get("result").MustBool()
	if result {
		e.LogTrade("撤销", e.Detail.Name, fmt.Sprintf("%+v", order))
	} else {
		e.LogError(order["StockType"], "撤销委托单失败，错误代码：", js.Get("error_code").MustInt())
	}
	return
}

func (e *OKCoin) GetOrder(order map[string]interface{}) (ret Order, err error) {
	params := []string{
		"api_key=" + e.Detail.Access,
		"symbol=" + okcoinStockType[fmt.Sprint(order["StockType"])],
		fmt.Sprint("order_id=", order["Id"]),
	}
	sign := e.sign(params)
	e.SetSleeper()
	get, err := httpPost(okcoinHost+"order_info.do", append(params, "sign="+sign))
	if err != nil {
		e.LogError(err)
		return
	}
	js, err := simplejson.NewJson(get)
	if err != nil {
		e.LogError(err)
		return
	}
	result := js.Get("result").MustBool()
	if result {
		orderJson := js.Get("orders").GetIndex(0)
		ret = Order{
			Id:         orderJson.Get("order_id").MustInt(),
			Price:      orderJson.Get("price").MustFloat64(),
			Amount:     orderJson.Get("amount").MustFloat64(),
			DealAmount: orderJson.Get("deal_amount").MustFloat64(),
			OrderType:  okcoinOrderType[orderJson.Get("type").MustString()],
			StockType:  fmt.Sprint((order["StockType"])),
			Status:     orderJson.Get("status").MustInt(),
		}
	} else {
		err = fmt.Errorf("%s", fmt.Sprint(order["StockType"], "获取订单详情出错，错误代码：", js.Get("error_code").MustInt()))
		e.LogError(err)
	}
	return
}

func (e *OKCoin) WithDraw(stockType, address string, amount float64) (id int, err error) {
	params := []string{
		"api_key=" + e.Detail.Access,
		"symbol=" + okcoinStockType[stockType],
		"chargefee=" + okcoinWithFee[stockType],
		"trade_pwd=" + e.Detail.Tradepw,
		"withdraw_address=" + address,
		fmt.Sprint("withdraw_amount=", amount),
	}
	sign := e.sign(params)
	e.SetSleeper()
	get, err := httpPost(okcoinHost+"withdraw.do", append(params, "sign="+sign))
	if err != nil {
		e.LogError(err)
		return
	}
	js, err := simplejson.NewJson(get)
	if err != nil {
		e.LogError(err)
		return
	}
	result := js.Get("result").MustBool()
	if result {
		e.LogTrade("提币", fmt.Sprint(amount, stockType))
		id = js.Get("withdraw_id").MustInt()
	} else {
		err = fmt.Errorf("%s", fmt.Sprint(stockType, "提币失败，错误代码：", js.Get("error_code").MustInt()))
		e.LogError(err)
	}
	return
}

func (e *OKCoin) GetOrders(stockType string) (orders []Order, err error) {
	params := []string{
		"api_key=" + e.Detail.Access,
		"symbol=" + okcoinStockType[stockType],
		"status=0",
		"current_page=1",
		"page_length=200",
	}
	sign := e.sign(params)
	e.SetSleeper()
	get, err := httpPost(okcoinHost+"order_history.do", append(params, "sign="+sign))
	if err != nil {
		e.LogError(err)
		return
	}
	js, err := simplejson.NewJson(get)
	if err != nil {
		e.LogError(err)
		return
	}
	result := js.Get("result").MustBool()
	if result {
		ordersJson := js.Get("orders")
		count := len(ordersJson.MustArray())
		for i := 0; i < count; i++ {
			orderJson := ordersJson.GetIndex(i)
			order := Order{
				Id:         orderJson.Get("order_id").MustInt(),
				Price:      orderJson.Get("price").MustFloat64(),
				Amount:     orderJson.Get("amount").MustFloat64(),
				DealAmount: orderJson.Get("deal_amount").MustFloat64(),
				OrderType:  okcoinOrderType[orderJson.Get("type").MustString()],
				StockType:  stockType,
			}
			orders = append(orders, order)
		}
	} else {
		err = fmt.Errorf("%s", fmt.Sprint(stockType, "获取未完成订单出错，错误代码：", js.Get("error_code").MustInt()))
		e.LogError(err)
	}
	return
}

func (e *OKCoin) GetTrades(stockType string) (orders []Order, err error) {
	params := []string{
		"api_key=" + e.Detail.Access,
		"symbol=" + okcoinStockType[stockType],
		"status=1",
		"current_page=1",
		"page_length=200",
	}
	sign := e.sign(params)
	e.SetSleeper()
	get, err := httpPost(okcoinHost+"order_history.do", append(params, "sign="+sign))
	if err != nil {
		e.LogError(err)
		return
	}
	js, err := simplejson.NewJson(get)
	if err != nil {
		e.LogError(err)
		return
	}
	result := js.Get("result").MustBool()
	if result {
		ordersJson := js.Get("orders")
		count := len(ordersJson.MustArray())
		for i := 0; i < count; i++ {
			orderJson := ordersJson.GetIndex(i)
			order := Order{
				Id:         orderJson.Get("order_id").MustInt(),
				Price:      orderJson.Get("price").MustFloat64(),
				Amount:     orderJson.Get("amount").MustFloat64(),
				DealAmount: orderJson.Get("deal_amount").MustFloat64(),
				OrderType:  okcoinOrderType[orderJson.Get("type").MustString()],
				StockType:  stockType,
			}
			orders = append(orders, order)
		}
	} else {
		err = fmt.Errorf("%s", fmt.Sprint(stockType, "获取交易历史出错，错误代码：", js.Get("error_code").MustInt()))
		e.LogError(err)
	}
	return
}

func (e *OKCoin) GetTicker(stockType string) (t Ticker, err error) {
	get, err := httpGet("https://www.okcoin.cn/api/v1/depth.do?symbol=" + okcoinStockType[stockType] + "&size=5")
	if err != nil {
		e.LogError(err)
		return
	}
	js, err := simplejson.NewJson(get)
	if err != nil {
		e.LogError(err)
		return
	}
	marketOrdersJson := js.Get("bids")
	count := len(marketOrdersJson.MustArray())
	for i := 0; i < count; i++ {
		marketOrder := marketOrdersJson.GetIndex(i)
		bid := MarketOrder{Price: marketOrder.GetIndex(0).MustFloat64(), Amount: marketOrder.GetIndex(1).MustFloat64()}
		t.Bids = append(t.Bids, bid)
	}
	marketOrdersJson = js.Get("asks")
	count = len(marketOrdersJson.MustArray())
	for i := count; i > 0; i-- {
		marketOrder := marketOrdersJson.GetIndex(i - 1)
		ask := MarketOrder{Price: marketOrder.GetIndex(0).MustFloat64(), Amount: marketOrder.GetIndex(1).MustFloat64()}
		t.Asks = append(t.Asks, ask)
	}
	t.Buy = t.Bids[0].Price
	t.Sell = t.Asks[0].Price
	t.Mid = (t.Buy + t.Sell) / 2
	return
}

func (e *OKCoin) GetTickerJser(stockType string) (t map[string]interface{}, err error) {
	get, err := httpGet("https://www.okcoin.cn/api/v1/depth.do?symbol=" + okcoinStockType[stockType] + "&size=5")
	if err != nil {
		e.LogError(err)
		return
	}
	js, err := simplejson.NewJson(get)
	if err != nil {
		e.LogError(err)
		return
	}
	marketOrdersJson := js.Get("bids")
	count := len(marketOrdersJson.MustArray())
	var Bids []map[string]interface{}
	Buy := marketOrdersJson.GetIndex(0).GetIndex(0).MustFloat64()
	for i := 0; i < count; i++ {
		marketOrder := marketOrdersJson.GetIndex(i)
		Bids = append(Bids, map[string]interface{}{"Price": marketOrder.GetIndex(0).MustFloat64(), "Amount": marketOrder.GetIndex(1).MustFloat64()})
	}
	marketOrdersJson = js.Get("asks")
	count = len(marketOrdersJson.MustArray())
	var Asks []map[string]interface{}
	Sell := marketOrdersJson.GetIndex(count - 1).GetIndex(0).MustFloat64()
	for i := count; i > 0; i-- {
		marketOrder := marketOrdersJson.GetIndex(i - 1)
		Asks = append(Asks, map[string]interface{}{"Price": marketOrder.GetIndex(0).MustFloat64(), "Amount": marketOrder.GetIndex(1).MustFloat64()})
	}
	Mid := (Buy + Sell) / 2
	t = map[string]interface{}{
		"Mid":  Mid,
		"Buy":  Buy,
		"Sell": Sell,
		"Bids": Bids,
		"Asks": Asks,
	}
	return
}

func (e *OKCoin) GetRecords(stockType, period string) (rs []Record, err error) {
	get, err := httpGet("https://www.okcoin.cn/api/v1/kline.do?symbol=" + okcoinStockType[stockType] + "&type=" + okcoinPeriod[period] + "&size=1000")
	if err != nil {
		e.LogError(err)
		return
	}
	js, err := simplejson.NewJson(get)
	if err != nil {
		e.LogError(err)
		return
	}
	timeLast := 0
	rs = e.Records[period]
	if len(rs) > 0 {
		timeLast = rs[len(rs)-1].Time
	}
	rsArrEnd := len(js.MustArray())
	count := 0
	var rsArr [1000]Record
	for i := rsArrEnd; i > 0; i-- {
		rJson := js.GetIndex(i - 1)
		rTimeInt64 := rJson.GetIndex(0).MustInt64() / 1000
		rTimeFmt := time.Unix(rTimeInt64, 0)
		rTimeInt, _ := strconv.Atoi(rTimeFmt.Format("20060102030405"))
		if rTimeInt < timeLast {
			break
		} else {
			rsArr[i-1] = Record{
				Time:   rTimeInt,
				Open:   rJson.GetIndex(1).MustFloat64(),
				High:   rJson.GetIndex(2).MustFloat64(),
				Low:    rJson.GetIndex(3).MustFloat64(),
				Close:  rJson.GetIndex(4).MustFloat64(),
				Volume: rJson.GetIndex(5).MustFloat64(),
			}
			count += 1
			if rTimeInt == timeLast {
				e.Records[period] = e.Records[period][:len(rs)-1]
			}
		}
	}
	e.Records[period] = append(e.Records[period], rsArr[rsArrEnd-count:rsArrEnd]...)
	if len(e.Records[period]) > 1000 {
		e.Records[period] = e.Records[period][len(e.Records[period])-1000:]
	}
	rs = e.Records[period]
	return
}

func (e *OKCoin) GetRecordsJser(stockType, period string) (rs []map[string]interface{}, err error) {
	get, err := httpGet("https://www.okcoin.cn/api/v1/kline.do?symbol=" + okcoinStockType[stockType] + "&type=" + okcoinPeriod[period] + "&size=1000")
	if err != nil {
		e.LogError(err)
		return
	}
	js, err := simplejson.NewJson(get)
	if err != nil {
		e.LogError(err)
		return
	}
	timeLast := 0
	rs = e.RecordsJser[period]
	if len(rs) > 0 {
		timeLast = rs[len(rs)-1]["Time"].(int)
	}
	rsArrEnd := len(js.MustArray())
	count := 0
	var rsArr [1000]map[string]interface{}
	for i := rsArrEnd; i > 0; i-- {
		rJson := js.GetIndex(i - 1)
		rTimeInt64 := rJson.GetIndex(0).MustInt64() / 1000
		rTimeFmt := time.Unix(rTimeInt64, 0)
		rTimeInt, _ := strconv.Atoi(rTimeFmt.Format("20060102150405"))
		if rTimeInt < timeLast {
			break
		} else {
			rsArr[i-1] = map[string]interface{}{
				"Time":   rTimeInt,
				"Open":   rJson.GetIndex(1).MustFloat64(),
				"High":   rJson.GetIndex(2).MustFloat64(),
				"Low":    rJson.GetIndex(3).MustFloat64(),
				"Close":  rJson.GetIndex(4).MustFloat64(),
				"Volume": rJson.GetIndex(5).MustFloat64(),
			}
			count += 1
			if rTimeInt == timeLast {
				e.RecordsJser[period] = e.RecordsJser[period][:len(rs)-1]
			}
		}
	}
	e.RecordsJser[period] = append(e.RecordsJser[period], rsArr[rsArrEnd-count:rsArrEnd]...)
	if len(e.RecordsJser[period]) > 1000 {
		e.RecordsJser[period] = e.RecordsJser[period][len(e.RecordsJser[period])-1000:]
	}
	rs = e.RecordsJser[period]
	return
}
