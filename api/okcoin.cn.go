package api

import (
	"fmt"
	"strings"

	"github.com/bitly/go-simplejson"
	"github.com/miaolz123/conver"
)

// OKCoinCn : the exchange struct of okcoin.cn
type OKCoinCn struct {
	stockMap  map[string]string
	host      string
	option    Option
	mainStock string
}

// NewOKCoinCn : create an exchange struct of okcoin.cn
func NewOKCoinCn(opt Option) *OKCoinCn {
	e := OKCoinCn{
		stockMap:  map[string]string{"BTC": "btc", "LTC": "ltc"},
		host:      "https://www.okcoin.cn/api/v1/",
		option:    opt,
		mainStock: "BTC",
	}
	return &e
}

// Log : Log
func (e *OKCoinCn) Log(msgs ...interface{}) {
	e.option.log.Do(e.option.Type, "info", 0.0, 0.0, msgs...)
}

// GetAccount : GetAccount
func (e *OKCoinCn) GetAccount() (map[string]interface{}, error) {
	account := make(map[string]interface{})
	params := []string{
		"api_key=" + e.option.AccessKey,
		"secret_key=" + e.option.SecretKey,
	}
	params = append(params, "sign="+strings.ToUpper(sign(params)))
	resp, err := post(e.host+"userinfo.do", params)
	if err != nil {
		e.option.log.Do(e.option.Type, "error", 0.0, 0.0, err)
		return account, err
	}
	json, err := simplejson.NewJson(resp)
	if err != nil {
		e.option.log.Do(e.option.Type, "error", 0.0, 0.0, err)
		return account, err
	}
	result := json.Get("result").MustBool()
	if result {
		account["Total"] = conver.Float64Must(json.GetPath("info", "funds", "asset", "total").Interface())
		account["Net"] = conver.Float64Must(json.GetPath("info", "funds", "asset", "net").Interface())
		account["Balance"] = conver.Float64Must(json.GetPath("info", "funds", "free", "cny").Interface())
		account["FrozenBalance"] = conver.Float64Must(json.GetPath("info", "funds", "freezed", "cny").Interface())
		account["BTC"] = conver.Float64Must(json.GetPath("info", "funds", "free", "btc").Interface())
		account["FrozenBTC"] = conver.Float64Must(json.GetPath("info", "funds", "freezed", "btc").Interface())
		account["LTC"] = conver.Float64Must(json.GetPath("info", "funds", "free", "ltc").Interface())
		account["FrozenLTC"] = conver.Float64Must(json.GetPath("info", "funds", "freezed", "ltc").Interface())
		account["Stocks"] = account[e.mainStock]
		account["FrozenStocks"] = account["Frozen"+e.mainStock]
	} else {
		err = fmt.Errorf("%s", fmt.Sprint("GetAccount() error, the error number is ", json.Get("error_code").MustInt()))
		e.option.log.Do(e.option.Type, "error", 0.0, 0.0, err)
	}
	return account, err
}
