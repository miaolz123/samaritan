package trader

import (
	"reflect"
	"time"

	"github.com/miaolz123/conver"
	"github.com/miaolz123/samaritan/api"
	"github.com/miaolz123/samaritan/constant"
)

type task struct {
	name string
	fn   reflect.Value
	args []reflect.Value
}

// Sleep ...
func (g *Global) Sleep(intervals ...interface{}) {
	interval := int64(0)
	if len(intervals) > 0 {
		interval = conver.Int64Must(intervals[0])
	}
	if interval > 0 {
		time.Sleep(time.Duration(interval * 1000000))
	}
}

// Log ...
func (g *Global) Log(msgs ...interface{}) {
	g.Logger.Log(constant.INFO, 0.0, 0.0, msgs...)
}

// LogProfit ...
func (g *Global) LogProfit(msgs ...interface{}) {
	profit := 0.0
	if len(msgs) > 0 {
		profit = conver.Float64Must(msgs[0])
	}
	g.Logger.Log(constant.PROFIT, 0.0, profit, msgs[1:]...)
}

// AddTask ...
func (g *Global) AddTask(e api.Exchange, name string, args ...interface{}) bool {
	t := task{}
	switch name {
	case "Log":
		t.fn = reflect.ValueOf(e.Log)
	case constant.GetAccount:
		t.fn = reflect.ValueOf(e.GetAccount)
	case constant.Buy:
		t.fn = reflect.ValueOf(e.Buy)
	case constant.Sell:
		t.fn = reflect.ValueOf(e.Sell)
	default:
		g.Logger.Log(constant.ERROR, 0.0, 0.0, "Invalid task name")
		return false
	}
	t.name = name
	for _, arg := range args {
		t.args = append(t.args, reflect.ValueOf(arg))
	}
	g.tasks = append(g.tasks, t)
	return true
}

// Do ...
func (g *Global) Do() (results []interface{}) {
	for _, t := range g.tasks {
		var result interface{}
		rs := t.fn.Call(t.args)
		if len(rs) > 0 {
			result = rs[0].Interface()
		}
		results = append(results, result)
	}
	return
}
