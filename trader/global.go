package trader

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/miaolz123/conver"
	"github.com/miaolz123/samaritan/constant"
	"github.com/robertkrimen/otto"
)

type task struct {
	fn   otto.Value
	args []interface{}
}

// Sleep ...
func (g *Global) Sleep(intervals ...interface{}) {
	interval := int64(0)
	if len(intervals) > 0 {
		interval = conver.Int64Must(intervals[0])
	}
	if interval > 0 {
		time.Sleep(time.Duration(interval * 1000000))
	} else {
		for _, e := range g.es {
			e.AutoSleep()
		}
	}
}

// Log ...
func (g *Global) Log(msgs ...interface{}) {
	g.Logger.Log(constant.INFO, "", 0.0, 0.0, msgs...)
}

// LogProfit ...
func (g *Global) LogProfit(msgs ...interface{}) {
	profit := 0.0
	if len(msgs) > 0 {
		profit = conver.Float64Must(msgs[0])
	}
	g.Logger.Log(constant.PROFIT, "", 0.0, profit, msgs[1:]...)
}

// LogStatus ...
func (g *Global) LogStatus(msgs ...interface{}) {
	go func() {
		msg := ""
		for _, m := range msgs {
			v := reflect.ValueOf(m)
			switch v.Kind() {
			case reflect.Struct, reflect.Map, reflect.Slice:
				if bs, err := json.Marshal(m); err == nil {
					msg += string(bs)
					continue
				}
			}
			msg += fmt.Sprintf("%+v", m)
		}
		g.statusLog = msg
	}()
}

// AddTask ...
func (g *Global) AddTask(fn otto.Value, args ...interface{}) bool {
	if g.execed {
		g.tasks = []task{}
		g.execed = false
	}
	if fn.Class() != "Function" {
		g.Logger.Log(constant.ERROR, "", 0.0, 0.0, "AddTask(), Invalid function")
	}
	g.tasks = append(g.tasks, task{fn: fn, args: args})
	return true
}

// ExecTasks ...
func (g *Global) ExecTasks() (results []interface{}) {
	g.execed = true
	for range g.tasks {
		results = append(results, false)
	}
	wg := sync.WaitGroup{}
	for i, t := range g.tasks {
		wg.Add(1)
		go func(i int, t task) {
			result, err := t.fn.Call(t.fn, t.args...)
			if err != nil || result.IsUndefined() || result.IsNull() || result.IsNaN() {
				results[i] = false
			} else {
				results[i] = result
			}
			wg.Done()
		}(i, t)
	}
	wg.Wait()
	return
}
