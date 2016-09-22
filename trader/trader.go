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
	db, err := model.NewOrm()
	if err != nil {
		return
	}
	if err = db.First(&trader, trader.ID).Error; err != nil {
		return
	}
	if err = db.First(&trader.Strategy, trader.StrategyID).Error; err != nil {
		return
	}
	if err = db.Model(&trader).Association("Exchanges").Find(&trader.Exchanges).Error; err != nil {
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
	exchanges := []interface{}{}
	for _, e := range trader.Exchanges {
		opt := api.Option{
			TraderID:  trader.ID,
			Type:      e.Type,
			AccessKey: e.AccessKey,
			SecretKey: e.SecretKey,
			MainStock: "BTC",
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
	trader.Ctx.Set("Sleep", func(call otto.FunctionCall) otto.Value {
		if t := conver.Int64Must(call.Argument(0).String()); t > 0 {
			time.Sleep(time.Duration(t * 1000000))
		}
		return otto.UndefinedValue()
	})
	trader.Ctx.Set("exchange", exchanges[0])
	trader.Ctx.Set("exchanges", exchanges)
	go func() {
		defer func() {
			if err := recover(); err != errHalt {
				Executor[trader.ID].Logger.Log(constant.ERROR, 0.0, 0.0, err)
			}
			trader.Status = 0
			Executor[trader.ID].Logger.Log(constant.INFO, 0.0, 0.0, "The Trader stop running")
		}()
		Executor[trader.ID].Logger.Log(constant.INFO, 0.0, 0.0, "The Trader us running")
		trader.Ctx.Run(trader.Strategy.Script)
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
