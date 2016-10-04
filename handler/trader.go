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
		td, err := model.GetTrader(self, req.ID)
		if err != nil {
			db.Rollback()
			resp["msg"] = fmt.Sprint(err)
			c.JSON(iris.StatusOK, resp)
			return
		}
		td.Name = req.Name
		td.StrategyID = req.StrategyID
		rs, err := model.GetTraderExchanges(self, td.ID)
		if err != nil {
			db.Rollback()
			resp["msg"] = fmt.Sprint(err)
			c.JSON(iris.StatusOK, resp)
			return
		}
		for i, r := range rs {
			if i >= len(req.Exchanges) {
				if err := db.Delete(&r).Error; err != nil {
					db.Rollback()
					resp["msg"] = fmt.Sprint(err)
					c.JSON(iris.StatusOK, resp)
					return
				}
				continue
			}
			if r.Exchange.ID == req.Exchanges[i].ID {
				continue
			}
			r.ExchangeID = req.Exchanges[i].ID
			if err := db.Save(&r).Error; err != nil {
				db.Rollback()
				resp["msg"] = fmt.Sprint(err)
				c.JSON(iris.StatusOK, resp)
				return
			}
		}
		for i, e := range req.Exchanges {
			if i < len(rs) {
				continue
			}
			r := model.TraderExchange{
				TraderID:   td.ID,
				ExchangeID: e.ID,
			}
			if err := db.Create(&r).Error; err != nil {
				db.Rollback()
				resp["msg"] = fmt.Sprint(err)
				c.JSON(iris.StatusOK, resp)
				return
			}
		}
		if err := db.Save(&td).Error; err != nil {
			db.Rollback()
			resp["msg"] = fmt.Sprint(err)
			c.JSON(iris.StatusOK, resp)
			return
		}
		if err := db.Commit().Error; err != nil {
			db.Rollback()
			resp["msg"] = fmt.Sprint(err)
			c.JSON(iris.StatusOK, resp)
			return
		}
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
	if err := db.Commit().Error; err != nil {
		db.Rollback()
		resp["msg"] = fmt.Sprint(err)
		c.JSON(iris.StatusOK, resp)
		return
	}
	resp["success"] = true
	c.JSON(iris.StatusOK, resp)
}

// Delete /trader
func (c traderHandler) Delete() {
	resp := iris.Map{
		"success": false,
		"msg":     "",
	}
	id := c.URLParam("id")
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
	td, err := model.GetTrader(self, id)
	if err != nil {
		db.Rollback()
		resp["msg"] = fmt.Sprint(err)
		c.JSON(iris.StatusOK, resp)
		return
	}
	if err := db.Delete(&td).Error; err != nil {
		db.Rollback()
		resp["msg"] = fmt.Sprint(err)
		c.JSON(iris.StatusOK, resp)
		return
	}
	if err := db.Where("trader_id = ?", id).Delete(&model.TraderExchange{}).Error; err != nil {
		db.Rollback()
		resp["msg"] = fmt.Sprint(err)
		c.JSON(iris.StatusOK, resp)
		return
	}
	if err := db.Commit().Error; err != nil {
		db.Rollback()
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
	if err := trader.Run(trader.Global{Trader: req}); err != nil {
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
	if err := trader.Stop(req.ID); err != nil {
		resp["msg"] = fmt.Sprint(err)
		c.JSON(iris.StatusOK, resp)
		return
	}
	resp["success"] = true
	c.JSON(iris.StatusOK, resp)
}
