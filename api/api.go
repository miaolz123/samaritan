package api

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/miaolz123/samaritan/candyjs"
	"github.com/miaolz123/samaritan/log"
)

// Robot ...
type Robot struct {
	ID         string
	Name       string
	CreateTime time.Time
	UpdateTime time.Time
	script     string
	ctx        *candyjs.Context
	waitGroup  sync.WaitGroup
}

// Option : exchange option
type Option struct {
	Type      string // one of ["okcoin.cn", "huobi"]
	AccessKey string
	SecretKey string
	MainStock string
}

// New : get a robot from opts(options) & scr(javascript code)
func New(opts []Option, name, scr string) *Robot {
	robot := &Robot{
		Name:       name,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
		script:     scr,
		ctx:        candyjs.NewContext(),
		waitGroup:  sync.WaitGroup{},
	}
	exchanges := []interface{}{}
	logger := log.New("global")
	for _, opt := range opts {
		switch opt.Type {
		case "okcoin.cn":
			exchanges = append(exchanges, NewOKCoinCn(opt))
		}
	}
	defer func() {
		if err := recover(); err != nil {
			logger.Do("error", 0.0, 0.0, err)
		}
	}()
	if len(exchanges) < 1 {
		panic("Please add at least one exchange")
	}
	robot.ctx.PushGlobalGoFunction("Log", func(msgs ...interface{}) {
		logger.Do("info", 0.0, 0.0, msgs...)
	})
	robot.ctx.PushGlobalGoFunction("Sleep", func(t float64) {
		time.Sleep(time.Duration(t * 1000000))
	})
	robot.ctx.PushGlobalInterface("exchange", exchanges[0])
	robot.ctx.PushGlobalInterface("exchanges", exchanges)
	return robot
}

// Run ...
func (robot *Robot) Run() {
	robot.waitGroup.Add(1)
	go robot.ctx.EvalString(robot.script)
	go robot.waitGroup.Wait()
}

// Stop ...
func (robot *Robot) Stop() {
	robot.waitGroup.Done()
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
