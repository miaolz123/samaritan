package handler

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris"
	"github.com/miaolz123/samaritan/model"
	"github.com/miaolz123/samaritan/trader"
)

type traderHandler struct {
	*iris.Context
}

// Get /trader
func (c traderHandler) Get() {
	resp := iris.Map{
		"success": false,
		"msg":     "",
	}
	self, err := model.GetUser(jwtmid.Get(c.Context).Claims.(jwt.MapClaims)["sub"])
	if err != nil {
		resp["msg"] = fmt.Sprint(err)
		c.JSON(iris.StatusOK, resp)
		return
	}
	traders, err := model.GetTraders(self)
	if err != nil {
		resp["msg"] = fmt.Sprint(err)
		c.JSON(iris.StatusOK, resp)
		return
	}
	for i := range traders {
		if t := trader.Executor[traders[i].ID]; t != nil {
			traders[i].Status = t.Status
		}
	}
	resp["success"] = true
	resp["data"] = traders
	c.JSON(iris.StatusOK, resp)
}

// Get /trader/:id
func (c traderHandler) GetBy(id string) {
	resp := iris.Map{
		"success": false,
		"msg":     "",
	}
	self, err := model.GetUser(jwtmid.Get(c.Context).Claims.(jwt.MapClaims)["sub"])
	if err != nil {
		resp["msg"] = fmt.Sprint(err)
		c.JSON(iris.StatusOK, resp)
		return
	}
	rs, err := model.GetTraderExchanges(self, id)
	if err != nil {
		resp["msg"] = fmt.Sprint(err)
		c.JSON(iris.StatusOK, resp)
		return
	}
	resp["success"] = true
	resp["data"] = rs
	c.JSON(iris.StatusOK, resp)
}

// Post /trader
func (c traderHandler) Post() {
	resp := iris.Map{
		"success": false,
		"msg":     "",
	}
	db, err := model.NewOrm()
	if err != nil {
		resp["msg"] = fmt.Sprint(err)
		c.JSON(iris.StatusOK, resp)
		return
	}
	defer db.Close()
	db = db.Begin()
	self, err := model.GetUser(jwtmid.Get(c.Context).Claims.(jwt.MapClaims)["sub"])
	if err != nil {
		db.Rollback()
		resp["msg"] = fmt.Sprint(err)
		c.JSON(iris.StatusOK, resp)
		return
	}
	req := model.Trader{}
	if err := c.ReadJSON(&req); err != nil {
		db.Rollback()
		resp["msg"] = fmt.Sprint(err)
		c.JSON(iris.StatusOK, resp)
		return
	}
	if req.ID > 0 {
		trader := model.Trader{}
		if err := db.First(&trader, req.ID).Error; err != nil {
			db.Rollback()
			resp["msg"] = fmt.Sprint(err)
			c.JSON(iris.StatusOK, resp)
			return
		}
		trader.Name = req.Name
		trader.Status = req.Status
		trader.Exchanges = req.Exchanges
		if err := db.Save(&trader).Error; err != nil {
			db.Rollback()
			resp["msg"] = fmt.Sprint(err)
			c.JSON(iris.StatusOK, resp)
			return
		}
		db.Commit()
		resp["success"] = true
		c.JSON(iris.StatusOK, resp)
		return
	}
	req.UserID = self.ID
	if err := db.Create(&req).Error; err != nil {
		db.Rollback()
		resp["msg"] = fmt.Sprint(err)
		c.JSON(iris.StatusOK, resp)
		return
	}
	for _, e := range req.Exchanges {
		traderExchange := model.TraderExchange{
			TraderID:   req.ID,
			ExchangeID: e.ID,
		}
		if err := db.Create(&traderExchange).Error; err != nil {
			db.Rollback()
			resp["msg"] = fmt.Sprint(err)
			c.JSON(iris.StatusOK, resp)
			return
		}
	}
	db.Commit()
	resp["success"] = true
	c.JSON(iris.StatusOK, resp)
}

// PostBy /trader/:id
func (c traderHandler) PostBy(id string) {
	resp := iris.Map{
		"success": false,
		"msg":     "",
	}
	self, err := model.GetUser(jwtmid.Get(c.Context).Claims.(jwt.MapClaims)["sub"])
	if err != nil {
		resp["msg"] = fmt.Sprint(err)
		c.JSON(iris.StatusOK, resp)
		return
	}
	td, err := model.GetTrader(self, id)
	if err != nil {
		resp["msg"] = fmt.Sprint(err)
		c.JSON(iris.StatusOK, resp)
		return
	}
	req := model.Exchange{}
	if err := c.ReadJSON(&req); err != nil {
		resp["msg"] = fmt.Sprint(err)
		c.JSON(iris.StatusOK, resp)
		return
	}
	if req.ID > 0 {
		traderExchange := model.TraderExchange{
			TraderID:   td.ID,
			ExchangeID: req.ID,
		}
		if err := model.DB.Create(&traderExchange).Error; err != nil {
			resp["msg"] = fmt.Sprint(err)
			c.JSON(iris.StatusOK, resp)
			return
		}
	}
	resp["success"] = true
	c.JSON(iris.StatusOK, resp)
}

// DeleteBy /trader/:id
func (c traderHandler) DeleteBy(id string) {
	resp := iris.Map{
		"success": false,
		"msg":     "",
	}
	self, err := model.GetUser(jwtmid.Get(c.Context).Claims.(jwt.MapClaims)["sub"])
	if err != nil {
		resp["msg"] = fmt.Sprint(err)
		c.JSON(iris.StatusOK, resp)
		return
	}
	if _, err := model.GetTrader(self, id); err != nil {
		resp["msg"] = fmt.Sprint(err)
		c.JSON(iris.StatusOK, resp)
		return
	}
	traderExchange := model.TraderExchange{}
	traderExchange.ID, _ = c.URLParamInt64("id")
	if err := model.DB.Unscoped().Delete(&traderExchange).Error; err != nil {
		resp["msg"] = fmt.Sprint(err)
		c.JSON(iris.StatusOK, resp)
		return
	}
	resp["success"] = true
	c.JSON(iris.StatusOK, resp)
}

// Post /run
func traderRun(c *iris.Context) {
	resp := iris.Map{
		"success": false,
		"msg":     "",
	}
	req := model.Trader{}
	if err := c.ReadJSON(&req); err != nil {
		resp["msg"] = fmt.Sprint(err)
		c.JSON(iris.StatusOK, resp)
		return
	}
	if err := trader.Run(req); err != nil {
		resp["msg"] = fmt.Sprint(err)
		c.JSON(iris.StatusOK, resp)
		return
	}
	resp["success"] = true
	c.JSON(iris.StatusOK, resp)
}

// Post /stop
func traderStop(c *iris.Context) {
	resp := iris.Map{
		"success": false,
		"msg":     "",
	}
	req := model.Trader{}
	if err := c.ReadJSON(&req); err != nil {
		resp["msg"] = fmt.Sprint(err)
		c.JSON(iris.StatusOK, resp)
		return
	}
	if err := trader.Stop(req); err != nil {
		resp["msg"] = fmt.Sprint(err)
		c.JSON(iris.StatusOK, resp)
		return
	}
	resp["success"] = true
	c.JSON(iris.StatusOK, resp)
}
