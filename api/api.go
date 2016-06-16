package api

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/miaolz123/samaritan/candyjs"
	"github.com/miaolz123/samaritan/log"
	"github.com/miaolz123/samaritan/task"
)

// Robot ...
type Robot struct {
	ID         string
	Name       string
	CreateTime time.Time
	UpdateTime time.Time
	script     string
	log        log.Logger
	ctx        *candyjs.Context
	runner     *task.Task
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
	constants := []string{
		"BTC",
		"LTC",
		"M",
		"M5",
		"M15",
		"M30",
		"H",
		"D",
		"W",
	}
	robot := &Robot{
		Name:       name,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
		script:     scr,
		log:        log.New("global"),
		ctx:        candyjs.NewContext(),
		runner:     task.New(),
	}
	exchanges := []interface{}{}
	for _, opt := range opts {
		switch opt.Type {
		case "okcoin.cn":
			exchanges = append(exchanges, NewOKCoinCn(opt))
		}
	}
	if len(exchanges) < 1 {
		robot.log.Do("error", 0.0, 0.0, "Please add at least one exchange")
	}
	for _, cons := range constants {
		robot.ctx.PushGlobalInterface(cons, cons)
	}
	robot.ctx.PushGlobalGoFunction("Log", func(msgs ...interface{}) {
		robot.log.Do("info", 0.0, 0.0, msgs...)
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
	robot.runner.Add(1)
	defer robot.Stop()
	robot.log.Do("info", 0.0, 0.0, "Start Running")
	if err := robot.ctx.PevalString(robot.script); err != nil {
		robot.log.Do("error", 0.0, 0.0, err)
	}
}

// Stop ...
func (robot *Robot) Stop() {
	if robot.runner.AllDone() {
		robot.log.Do("info", 0.0, 0.0, "Stop Running")
	}
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
