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
		ExchangeType: "global",
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
		}
	}
	if len(exchanges) == 0 {
		err = fmt.Errorf("Please add at least one exchange")
		return
	}
	g := global{trader: trader}
	trader.Ctx.Set("Global", g)
	trader.Ctx.Set("G", g)
	trader.Ctx.Set("Exchange", exchanges[0])
	trader.Ctx.Set("E", exchanges[0])
	trader.Ctx.Set("Exchanges", exchanges)
	trader.Ctx.Set("Es", exchanges)
	go func() {
		defer func() {
			if err := recover(); err != nil && err != errHalt {
				Executor[trader.ID].Logger.Log(constant.ERROR, 0.0, 0.0, err)
			}
			trader.Status = 0
			Executor[trader.ID].Logger.Log(constant.INFO, 0.0, 0.0, "The Trader stop running")
		}()
		trader.LastRunAt = time.Now().Unix()
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
