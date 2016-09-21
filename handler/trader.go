package handler

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris"
	"github.com/miaolz123/samaritan/api"
	"github.com/miaolz123/samaritan/candyjs"
	"github.com/miaolz123/samaritan/constant"
	"github.com/miaolz123/samaritan/model"
)

type traderHandler struct {
	*iris.Context
}

// Get /trader
func (c traderHandler) Get() {
	self, err := model.GetUser(jwtmid.Get(c.Context).Claims.(jwt.MapClaims)["sub"])
	if err != nil {
		c.Error(fmt.Sprint(err), iris.StatusServiceUnavailable)
		return
	}
	traders, err := model.GetTraders(self)
	if err != nil {
		c.Error(fmt.Sprint(err), iris.StatusServiceUnavailable)
		return
	}
	c.JSON(iris.StatusOK, traders)
}

// Post /trader
func (c traderHandler) Post() {
	db, err := model.NewOrm()
	if err != nil {
		c.Error(fmt.Sprint(err), iris.StatusServiceUnavailable)
		return
	}
	self, err := model.GetUser(jwtmid.Get(c.Context).Claims.(jwt.MapClaims)["sub"])
	if err != nil {
		c.Error(fmt.Sprint(err), iris.StatusServiceUnavailable)
		return
	}
	req := model.Trader{}
	if err := c.ReadJSON(&req); err != nil {
		c.Error(fmt.Sprint(err), iris.StatusBadRequest)
		return
	}
	req.UserID = self.ID
	if err := db.Create(&req).Error; err != nil {
		c.Error(fmt.Sprint(err), iris.StatusServiceUnavailable)
		return
	}
	c.JSON(iris.StatusOK, req)
}

// Put /trader
func (c traderHandler) Put() {
	db, err := model.NewOrm()
	if err != nil {
		c.Error(fmt.Sprint(err), iris.StatusServiceUnavailable)
		return
	}
	req := model.Trader{}
	if err := c.ReadJSON(&req); err != nil {
		c.Error(fmt.Sprint(err), iris.StatusBadRequest)
		return
	}
	trader := model.Trader{}
	if err := db.First(&trader, req.ID).Error; err != nil {
		c.Error(fmt.Sprint(err), iris.StatusServiceUnavailable)
		return
	}
	trader.Name = req.Name
	trader.Status = req.Status
	trader.Exchanges = req.Exchanges
	if err := db.Save(&trader).Error; err != nil {
		c.Error(fmt.Sprint(err), iris.StatusServiceUnavailable)
		return
	}
	c.JSON(iris.StatusOK, trader)
}

// Post /run
func traderRun(c *iris.Context) {
	fmt.Println(888888)
	trader := model.Trader{}
	if err := c.ReadJSON(&trader); err != nil {
		c.Error(fmt.Sprint(err), iris.StatusBadRequest)
		return
	}
	fmt.Println(949494)
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
	fmt.Println(108108)
	trader.Logger = model.Logger{
		TraderID:     trader.ID,
		ExchangeType: "",
	}
	trader.Ctx = candyjs.NewContext()
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
	fmt.Println(125125)
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
		}
	}
	fmt.Println(142142)
	if len(exchanges) == 0 {
		trader.Logger.Log(constant.ERROR, 0.0, 0.0, "Please add at least one exchange")
	}
	for _, c := range constants {
		trader.Ctx.PushGlobalInterface(c, c)
	}
	trader.Ctx.PushGlobalGoFunction("Log", func(msgs ...interface{}) {
		trader.Logger.Log(constant.INFO, 0.0, 0.0, msgs...)
	})
	trader.Ctx.PushGlobalGoFunction("Sleep", func(t float64) {
		time.Sleep(time.Duration(t * 1000000))
	})
	fmt.Println(155155)
	trader.Ctx.PushGlobalInterface("exchange", exchanges[0])
	trader.Ctx.PushGlobalInterface("exchanges", exchanges)
	model.TraderMap[trader.ID] = &trader
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()
		model.TraderMap[trader.ID].Run()
	}()
	fmt.Println(167167)
	c.JSON(iris.StatusOK, trader)
}

// Post /stop
func traderStop(c *iris.Context) {
	req := model.Trader{}
	if err := c.ReadJSON(&req); err != nil {
		c.Error(fmt.Sprint(err), iris.StatusBadRequest)
		return
	}
	model.StopTrader(req)
	c.JSON(iris.StatusOK, req)
}
