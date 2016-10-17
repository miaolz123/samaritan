package api

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var client = http.DefaultClient

// // Account struct
// type Account struct {
// 	Total         float64
// 	Net           float64
// 	Balance       float64
// 	FrozenBalance float64
// 	BTC           float64
// 	FrozenBTC     float64
// 	LTC           float64
// 	FrozenLTC     float64
// 	Stock         float64
// 	FrozenStock   float64
// }

// Position struct
type Position struct {
	Price         float64
	Leverage      int
	Amount        float64
	ConfirmAmount float64
	FrozenAmount  float64
	Profit        float64
	ContractType  string
	TradeType     string
	StockType     string
}

// Order struct
type Order struct {
	ID         string
	Price      float64
	Amount     float64
	DealAmount float64
	TradeType  string
	StockType  string
}

// Ticker struct
type Ticker struct {
	Bids []OrderBook
	Buy  float64
	Mid  float64
	Sell float64
	Asks []OrderBook
}

// OrderBook struct
type OrderBook struct {
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

func base64Encode(data string) string {
	return base64.StdEncoding.EncodeToString([]byte(data))
}

func signMd5(params []string) string {
	m := md5.New()
	m.Write([]byte(strings.Join(params, "&")))
	return hex.EncodeToString(m.Sum(nil))
}

func signSha512(params []string, key string) string {
	h := hmac.New(sha512.New, []byte(key))
	h.Write([]byte(strings.Join(params, "&")))
	return hex.EncodeToString(h.Sum(nil))
}

func signSha1(params []string, key string) string {
	h := hmac.New(sha1.New, []byte(key))
	h.Write([]byte(strings.Join(params, "&")))
	return hex.EncodeToString(h.Sum(nil))
}

func signChbtc(params []string, key string) string {
	sha := sha1.New()
	sha.Write([]byte(key))
	secret := hex.EncodeToString(sha.Sum(nil))
	h := hmac.New(md5.New, []byte(secret))
	h.Write([]byte(strings.Join(params, "&")))
	return hex.EncodeToString(h.Sum(nil))
}

func post(url string, data []string) (ret []byte, err error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(strings.Join(data, "&")))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
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

func get(url string) (ret []byte, err error) {
	req, err := http.NewRequest("GET", url, strings.NewReader(""))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
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
