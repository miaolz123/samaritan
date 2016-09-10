package api

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Account struct
type Account struct {
	Total         float64
	Net           float64
	Balance       float64
	FrozenBalance float64
	BTC           float64
	FrozenBTC     float64
	LTC           float64
	FrozenLTC     float64
	Stock         float64
	FrozenStock   float64
}

// Order struct
type Order struct {
	ID         string
	Price      float64
	Amount     float64
	DealAmount float64
	OrderType  int
	StockType  string
}

// Ticker struct
type Ticker struct {
	Bids []MarketOrder
	Buy  float64
	Mid  float64
	Sell float64
	Asks []MarketOrder
}

// MarketOrder struct
type MarketOrder struct {
	Price  float64
	Amount float64
}

// Record struct
type Record struct {
	Time   int64
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume float64
}

func signMd5(params []string) string {
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
