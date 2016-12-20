package trader

import (
	"fmt"
	"time"

	"github.com/miaolz123/samaritan/api"
	"github.com/miaolz123/samaritan/constant"
	"github.com/miaolz123/samaritan/model"
	"github.com/robertkrimen/otto"
)

// Trader Variable
var (
	Executor      = make(map[int64]*Global)
	errHalt       = fmt.Errorf("HALT")
	exchangeMaker = map[string]func(api.Option) api.Exchange{
		constant.OkCoinCn:     api.NewOKCoinCn,
		constant.Huobi:        api.NewHuobi,
		constant.Poloniex:     api.NewPoloniex,
		constant.Btcc:         api.NewBtcc,
		constant.Chbtc:        api.NewChbtc,
		constant.OkcoinFuture: api.NewOKCoinFuture,
		constant.OandaV20:     api.NewOandaV20,
	}
)

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

// init ...
func (trader *Global) init() (err error) {
	if t := Executor[trader.ID]; t != nil && t.Status > 0 {
		return
	}
	if err := model.DB.First(&trader.Trader, trader.ID).Error; err != nil {
		return err
	}
	self, err := model.GetUserByID(trader.UserID)
	if err != nil {
		return err
	}
	if trader.AlgorithmID <= 0 {
		return fmt.Errorf("Please select a algorithm")
	}
	if err := model.DB.First(&trader.Algorithm, trader.AlgorithmID).Error; err != nil {
		return err
	}
	es, err := self.GetTraderExchanges(trader.ID)
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
	for _, c := range constant.Consts {
		trader.Ctx.Set(c, c)
	}
	for _, e := range es {
		if maker, ok := exchangeMaker[e.Type]; ok {
			opt := api.Option{
				TraderID:  trader.ID,
				Type:      e.Type,
				Name:      e.Name,
				AccessKey: e.AccessKey,
				SecretKey: e.SecretKey,
				// Ctx:       trader.Ctx,
			}
			trader.es = append(trader.es, maker(opt))
		}
	}
	if len(trader.es) == 0 {
		return fmt.Errorf("Please add at least one exchange")
	}
	trader.Ctx.Set("Global", &trader)
	trader.Ctx.Set("G", &trader)
	trader.Ctx.Set("Exchange", trader.es[0])
	trader.Ctx.Set("E", trader.es[0])
	trader.Ctx.Set("Exchanges", trader.es)
	trader.Ctx.Set("Es", trader.es)
	return
}

// Run ...
func Run(trader Global) (err error) {
	if err := trader.init(); err != nil {
		return err
	}
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
		trader.LastRunAt = time.Now()
		trader.Status = 1
		trader.Logger.Log(constant.INFO, "", 0.0, 0.0, "The Trader is running")
		if _, err := trader.Ctx.Run(trader.Algorithm.Script); err != nil {
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
func GetStatus(id int64) (status string) {
	if t := Executor[id]; t != nil {
		status = t.statusLog
	}
	return
}

// Stop ...
func Stop(id int64) (err error) {
	if t, ok := Executor[id]; !ok || t == nil {
		return fmt.Errorf("Can not found the Trader")
	}
	Executor[id].Ctx.Interrupt <- func() { panic(errHalt) }
	return
}

// Clean ...
func Clean(userID int64) {
	for _, t := range Executor {
		if t != nil && t.UserID == userID {
			Stop(t.ID)
		}
	}
}
