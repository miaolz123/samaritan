package handler

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris"
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
	req := model.Trader{}
	if err := c.ReadJSON(&req); err != nil {
		c.Error(fmt.Sprint(err), iris.StatusBadRequest)
		return
	}
	if err := model.RunTrader(req); err != nil {
		c.Error(fmt.Sprint(err), iris.StatusBadRequest)
		return
	}
	c.JSON(iris.StatusOK, req)
}

// Post /stop
func traderStop(c *iris.Context) {
	req := model.Trader{}
	if err := c.ReadJSON(&req); err != nil {
		c.Error(fmt.Sprint(err), iris.StatusBadRequest)
		return
	}
	if model.StopTrader(req) {
		c.JSON(iris.StatusOK, req)
	}
	c.Error("", iris.StatusServiceUnavailable)
}
