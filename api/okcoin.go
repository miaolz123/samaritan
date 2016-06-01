package api

import (
	"fmt"

	"github.com/cihub/seelog"
	"github.com/robertkrimen/otto"
)

// OKCoin ...
type OKCoin struct {
	option Option
	Logger seelog.LoggerInterface
}

// NewAPI ...
func (e *OKCoin) NewAPI(opt Option) map[string]func(otto.FunctionCall) otto.Value {
	return map[string]func(otto.FunctionCall) otto.Value{
		"Log": func(call otto.FunctionCall) otto.Value {
			var msgs []interface{}
			for _, msg := range call.ArgumentList {
				m, _ := msg.Export()
				msgs = append(msgs, m)
			}
			e.log(msgs...)
			return otto.TrueValue()
		},
	}
}

func (e *OKCoin) log(a ...interface{}) {
	fmt.Println(a...)
}
