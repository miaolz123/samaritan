package api

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"

	"github.com/miaolz123/go-duktape/candyjs"
	"github.com/miaolz123/samaritan/log"
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
	GetAccount() (map[string]interface{}, error)
}

// Run : run a strategy from opts(options) & scr(script)
func Run(opts []Option, scr string) {
	exchanges := []Exchange{}
	logger := log.New("test")
	for _, opt := range opts {
		opt.log = logger
		switch opt.Type {
		case "okcoin.cn":
			exchanges = append(exchanges, NewOKCoinCn(opt))
		}
	}
	ctx := candyjs.NewContext()
	defer func() {
		if err := recover(); err != nil {
			logger.Do("global", "error", 0.0, 0.0, err)
		}
	}()
	if len(exchanges) < 1 {
		panic("Please add at least one Exchange")
	}
	ctx.PushGlobalGoFunction("log", func(a ...interface{}) {
		logger.Do("global", "info", 0.0, 0.0, a...)
	})
	_, err := ctx.PushGlobalGoFunction("add", func(a, b float64) float64 {
		return a + b
	})
	if err != nil {
		logger.Do("global", "error", 0.0, 0.0, err)
	}
	ctx.PushGlobalInterface("exchange", exchanges[0])
	ctx.PushGlobalInterface("exchanges", exchanges)
	ctx.EvalString(scr)
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
	if resp == nil {
		err = fmt.Errorf("[POST %s] HTTP Error Info: %v", url, err)
	} else if resp.StatusCode == 200 {
		ret, _ = ioutil.ReadAll(resp.Body)
		resp.Body.Close()
	} else {
		err = fmt.Errorf("[POST %s] HTTP Status: %d, Info: %v", url, resp.StatusCode, err)
	}
	return ret, err
}

func get(url string) ([]byte, error) {
	var ret []byte
	resp, err := http.Get(url)
	if resp == nil {
		err = fmt.Errorf("[GET %s] HTTP Error Info: %v", url, err)
	} else if resp.StatusCode == 200 {
		ret, _ = ioutil.ReadAll(resp.Body)
		resp.Body.Close()
	} else {
		err = fmt.Errorf("[GET %s] HTTP Status: %d, Info: %v", url, resp.StatusCode, err)
	}
	return ret, err
}
