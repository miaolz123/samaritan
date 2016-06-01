package api

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/miaolz123/conver"
	"github.com/miaolz123/samaritan/log"
	"github.com/robertkrimen/otto"
)

// Option : exchange option
type Option struct {
	Type      string // one of ["okcoin.cn", "huobi"]
	AccessKey string
	SecretKey string
	log       log.Logger
}

// Exchange : exchange interface
type Exchange interface {
	GetMethods() map[string]func(otto.FunctionCall) otto.Value
}

// Run : run a strategy from opts(options) & scr(script)
func Run(opts []Option, scr interface{}) {
	exchanges := []Exchange{}
	logger := log.New("test")
	for _, opt := range opts {
		opt.log = logger
		switch opt.Type {
		case "okcoin.cn":
			exchanges = append(exchanges, NewOKCoinCn(opt))
		}
	}
	vm := otto.New()
	defer func() {
		if r := recover(); r != nil {
			logger.Do("global", "error", fmt.Sprint(r), 0.0, 0.0)
		} else {
			endScript, _ := vm.Get("end")
			endScript.Call(endScript)
			logger.Do("global", "info", "End succeed", 0.0, 0.0)
		}
	}()
	if len(exchanges) < 1 {
		panic("Please add at least one Exchange")
	}
	es := []map[string]func(otto.FunctionCall) otto.Value{}
	for _, e := range exchanges {
		es = append(es, e.GetMethods())
	}
	vm.Set("Log", func(call otto.FunctionCall) otto.Value {
		msgs := ""
		for _, msg := range call.ArgumentList {
			m, _ := msg.Export()
			msgs += conver.StringMust(m, "undefined")
		}
		logger.Do("global", "info", msgs, 0.0, 0.0)
		return otto.TrueValue()
	})
	vm.Set("LogProfit", func(call otto.FunctionCall) otto.Value {
		pro, err := call.Argument(0).ToFloat()
		if err != nil {
			logger.Do("global", "error", "The first argument of LogProfit() must be a number", 0.0, 0.0)
			return otto.FalseValue()
		}
		msgs := ""
		for i, msg := range call.ArgumentList {
			if i == 0 {
				continue
			}
			m, _ := msg.Export()
			msgs += conver.StringMust(m, "undefined")
		}
		logger.Do("global", "profit", msgs, pro, 0.0)
		return otto.TrueValue()
	})
	vm.Set("Sleep", func(call otto.FunctionCall) otto.Value {
		sleeper, err := call.Argument(0).ToFloat()
		if err != nil {
			return otto.FalseValue()
		}
		time.Sleep(time.Duration(sleeper * 1000000))
		return otto.TrueValue()
	})
	vm.Set("exchange", es[0])
	vm.Set("exchanges", es)
	if _, err := vm.Run(scr); err != nil {
		panic(err)
	}
	mainScript, err := vm.Get("main")
	if err != nil || mainScript.IsUndefined() {
		panic("Can not find the function: main()")
	}
	_, err = mainScript.Call(mainScript)
	if err != nil {
		panic(err)
	}
}

func sign(params []string) string {
	sort.Strings(params)
	m := md5.New()
	m.Write([]byte(strings.Join(params, "&")))
	return hex.EncodeToString(m.Sum(nil))
}

func post(url string, data []string) ([]byte, error) {
	var ret []byte
	resp, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(strings.Join(data, "&")))
	if resp.StatusCode == 200 {
		ret, _ = ioutil.ReadAll(resp.Body)
		resp.Body.Close()
	} else {
		err = fmt.Errorf("HTTP Status: %d, Info: %v", resp.StatusCode, err)
	}
	return ret, err
}

func get(url string) ([]byte, error) {
	var ret []byte
	resp, err := http.Get(url)
	if resp.StatusCode == 200 {
		ret, _ = ioutil.ReadAll(resp.Body)
		resp.Body.Close()
	} else {
		err = fmt.Errorf("HTTP Status: %d, Info: %v", resp.StatusCode, err)
	}
	return ret, err
}
