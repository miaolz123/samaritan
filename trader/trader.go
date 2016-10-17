package trader

import (
	"fmt"
	"time"

	"github.com/miaolz123/samaritan/api"
	"github.com/miaolz123/samaritan/constant"
	"github.com/miaolz123/samaritan/model"
	"github.com/robertkrimen/otto"
)

// Executor ...
var Executor = make(map[uint]*Global)
var errHalt = fmt.Errorf("HALT")

// Global ...
type Global struct {
	model.Trader
	Logger    model.Logger
	Ctx       *otto.Otto
	es        []api.Exchange
	tasks     []task
	execed    bool
	statusLog string
}

// Run ...
func Run(trader Global) (err error) {
	if t := Executor[trader.ID]; t != nil && t.Status > 0 {
		return
	}
	if err = model.DB.First(&trader.Trader, trader.ID).Error; err != nil {
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
		ExchangeType: "global",
	}
	trader.tasks = []task{}
	trader.Ctx = otto.New()
	trader.Ctx.Interrupt = make(chan func(), 1)
	for _, c := range constant.CONSTS {
		trader.Ctx.Set(c, c)
	}
	for _, e := range es {
		opt := api.Option{
			TraderID:  trader.ID,
			Type:      e.Type,
			Name:      e.Name,
			AccessKey: e.AccessKey,
			SecretKey: e.SecretKey,
			Ctx:       trader.Ctx,
		}
		switch opt.Type {
		case constant.OkCoinCn:
			trader.es = append(trader.es, api.NewOKCoinCn(opt))
		case constant.Huobi:
			trader.es = append(trader.es, api.NewHuobi(opt))
		case constant.Poloniex:
			trader.es = append(trader.es, api.NewPoloniex(opt))
		case constant.Btcc:
			trader.es = append(trader.es, api.NewBtcc(opt))
		case constant.Chbtc:
			trader.es = append(trader.es, api.NewChbtc(opt))
		case constant.OkcoinFuture:
			trader.es = append(trader.es, api.NewOKCoinFuture(opt))
		}
	}
	if len(trader.es) == 0 {
		err = fmt.Errorf("Please add at least one exchange")
		return
	}
	trader.Ctx.Set("Global", &trader)
	trader.Ctx.Set("G", &trader)
	trader.Ctx.Set("Exchange", trader.es[0])
	trader.Ctx.Set("E", trader.es[0])
	trader.Ctx.Set("Exchanges", trader.es)
	trader.Ctx.Set("Es", trader.es)
	go func() {
		defer func() {
			if err := recover(); err != nil && err != errHalt {
				trader.Logger.Log(constant.ERROR, "", 0.0, 0.0, err)
			}
			if exit, err := trader.Ctx.Get("exit"); err == nil && exit.IsFunction() {
				if _, err := exit.Call(exit); err != nil {
					trader.Logger.Log(constant.ERROR, "", 0.0, 0.0, err)
				}
			}
			trader.Status = 0
			trader.Logger.Log(constant.INFO, "", 0.0, 0.0, "The Trader stop running")
		}()
		trader.LastRunAt = time.Now().Unix()
		trader.Status = 1
		trader.Logger.Log(constant.INFO, "", 0.0, 0.0, "The Trader is running")
		if _, err := trader.Ctx.Run(trader.Strategy.Script); err != nil {
			trader.Logger.Log(constant.ERROR, "", 0.0, 0.0, err)
		}
		if main, err := trader.Ctx.Get("main"); err != nil || !main.IsFunction() {
			trader.Logger.Log(constant.ERROR, "", 0.0, 0.0, "Can not get the main function")
		} else {
			if _, err := main.Call(main); err != nil {
				trader.Logger.Log(constant.ERROR, "", 0.0, 0.0, err)
			}
		}
	}()
	Executor[trader.ID] = &trader
	return
}

// GetStatus ...
func GetStatus(id uint) string {
	t := Executor[id]
	if t != nil {
		return t.statusLog
	}
	return ""
}

// Stop ...
func Stop(id uint) (err error) {
	t := Executor[id]
	if t == nil {
		err = fmt.Errorf("Can not found the Trader")
		return
	}
	Executor[id].Ctx.Interrupt <- func() {
		panic(errHalt)
	}
	return
}

// Clean ...
func Clean(userID uint) {
	for _, t := range Executor {
		if t.UserID == userID {
			Stop(t.ID)
		}
	}
}
