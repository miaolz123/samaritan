package api

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"

	"github.com/robertkrimen/otto"
)

// Option ...
type Option struct {
	AccessKey string
	SecretKey string
}

// API exchange api collection
type API interface {
	NewAPI(opt Option) map[string]func(otto.FunctionCall) otto.Value
}

type exchangeConf struct {
	name   string
	access string
	secret string
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
