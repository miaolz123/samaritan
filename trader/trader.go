package trader

import (
	"fmt"
	"strings"
	"time"

	"github.com/miaolz123/conver"
	"github.com/miaolz123/samaritan/api"
	"github.com/miaolz123/samaritan/constant"
	"github.com/miaolz123/samaritan/model"
	"github.com/robertkrimen/otto"
)

// Executor ...
var Executor = make(map[uint]*model.Trader)
var errHalt = fmt.Errorf("HALT")

// Run ...
func Run(trader model.Trader) (err error) {
	if t := Executor[trader.ID]; t != nil && t.Status > 0 {
		return
	}
	if err = model.DB.First(&trader, trader.ID).Error; err != nil {
		return
	}
	self, err := model.GetUserByID(trader.UserID)
	if err != nil {
		return
	}
	if trader.StrategyID <= 0 {
		err = fmt.Errorf("Please select a strategy")
		return
	}
	if err = model.DB.First(&trader.Strategy, trader.StrategyID).Error; err != nil {
		return
	}
	es, err := model.GetTraderExchanges(self, trader.ID)
	if err != nil {
		return
	}
	trader.Logger = model.Logger{
		TraderID:     trader.ID,
		ExchangeType: "",
	}
	trader.Ctx = otto.New()
	trader.Ctx.Interrupt = make(chan func(), 1)
	for _, c := range constant.CONSTS {
		trader.Ctx.Set(c, c)
	}
	exchanges := []api.Exchange{}
	for _, e := range es {
		opt := api.Option{
			TraderID:  trader.ID,
			Type:      e.Type,
			AccessKey: e.AccessKey,
			SecretKey: e.SecretKey,
			MainStock: "BTC",
			Ctx:       trader.Ctx,
		}
		switch opt.Type {
		case "okcoin.cn":
			exchanges = append(exchanges, api.NewOKCoinCn(opt))
		case "huobi":
			exchanges = append(exchanges, api.NewHuobi(opt))
			// default:
			// 	exchanges = append(exchanges, api.NewHuobi(opt))
		}
	}
	if len(exchanges) == 0 {
		err = fmt.Errorf("Please add at least one exchange")
		return
	}
	trader.Ctx.Set("Log", func(call otto.FunctionCall) otto.Value {
		message := ""
		for _, a := range call.ArgumentList {
			message += fmt.Sprintf("%+v, ", a)
		}
		trader.Logger.Log(constant.INFO, 0.0, 0.0, strings.TrimSuffix(message, ", "))
		return otto.UndefinedValue()
	})
	trader.Ctx.Set("LogProfit", func(call otto.FunctionCall) otto.Value {
		profit := 0.0
		message := ""
		for i, a := range call.ArgumentList {
			if i == 0 {
				profit = conver.Float64Must(a)
				continue
			}
			message += fmt.Sprintf("%+v, ", a)
		}
		trader.Logger.Log(constant.PROFIT, 0.0, profit, strings.TrimSuffix(message, ", "))
		return otto.UndefinedValue()
	})
	trader.Ctx.Set("Sleep", func(call otto.FunctionCall) otto.Value {
		if t := conver.Int64Must(call.Argument(0).String()); t > 0 {
			time.Sleep(time.Duration(t * 1000000))
		}
		return otto.UndefinedValue()
	})
	trader.Ctx.Set("Exchange", exchanges[0])
	trader.Ctx.Set("Exchanges", exchanges)
	go func() {
		defer func() {
			if err := recover(); err != nil && err != errHalt {
				Executor[trader.ID].Logger.Log(constant.ERROR, 0.0, 0.0, err)
			}
			trader.Status = 0
			Executor[trader.ID].Logger.Log(constant.INFO, 0.0, 0.0, "The Trader stop running")
		}()
		trader.Status = 1
		Executor[trader.ID].Logger.Log(constant.INFO, 0.0, 0.0, "The Trader is running")
		if _, err := trader.Ctx.Run(trader.Strategy.Script); err != nil {
			Executor[trader.ID].Logger.Log(constant.ERROR, 0.0, 0.0, err)
		}
	}()
	Executor[trader.ID] = &trader
	return
}

// Stop ...
func Stop(trader model.Trader) (err error) {
	t := Executor[trader.ID]
	if t == nil {
		err = fmt.Errorf("Can not found the Trader")
		return
	}
	Executor[trader.ID].Ctx.Interrupt <- func() {
		panic(errHalt)
	}
	Executor[trader.ID].Status = 0
	return
}

// Clean ...
func Clean(userID uint) {
	for _, t := range Executor {
		if t.UserID == userID {
			Stop(*t)
		}
	}
}
