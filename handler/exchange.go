package handler

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/hprose/hprose-golang/rpc"
	"github.com/kataras/iris"
	"github.com/miaolz123/samaritan/constant"
	"github.com/miaolz123/samaritan/model"
)

type exchange struct{}

// Types ...
func (exchange) Types(_ string, ctx rpc.Context) (resp response) {
	resp.Data = constant.ExchangeTypes
	resp.Success = true
	return
}

// List ...
func (exchange) List(size, page int64, ctx rpc.Context) (resp response) {
	username := ctx.GetString("username")
	if username == "" {
		resp.Message = "Authorization wrong"
		return
	}
	self, err := model.GetUser(username)
	if err != nil {
		resp.Message = fmt.Sprint(err)
		return
	}
	total, exchanges, err := self.ExchangeList(size, page)
	if err != nil {
		resp.Message = fmt.Sprint(err)
		return
	}
	resp.Data = struct {
		Total int64
		List  []model.Exchange
	}{
		Total: total,
		List:  exchanges,
	}
	resp.Success = true
	return
}

// Put
func (exchange) Put(req model.Exchange, ctx rpc.Context) (resp response) {
	username := ctx.GetString("username")
	if username == "" {
		resp.Message = "Authorization wrong"
		return
	}
	self, err := model.GetUser(username)
	if err != nil {
		resp.Message = fmt.Sprint(err)
		return
	}
	exchange := req
	if req.ID > 0 {
		if err := model.DB.First(&exchange, req.ID).Error; err != nil {
			resp.Message = fmt.Sprint(err)
			return
		}
		exchange.Name = req.Name
		exchange.Type = req.Type
		exchange.AccessKey = req.AccessKey
		exchange.SecretKey = req.SecretKey
		if err := model.DB.Save(&exchange).Error; err != nil {
			resp.Message = fmt.Sprint(err)
			return
		}
		resp.Success = true
		return
	}
	req.UserID = self.ID
	if err := model.DB.Create(&req).Error; err != nil {
		resp.Message = fmt.Sprint(err)
		return
	}
	resp.Success = true
	return
}

// Delete
func (exchange) Delete(ids []int64, ctx rpc.Context) (resp response) {
	username := ctx.GetString("username")
	if username == "" {
		resp.Message = "Authorization wrong"
		return
	}
	self, err := model.GetUser(username)
	if err != nil {
		resp.Message = fmt.Sprint(err)
		return
	}
	userIds := []int64{}
	if _, users, err := self.UserList(-1, 1); err != nil {
		resp.Message = fmt.Sprint(err)
		return
	} else {
		for _, u := range users {
			userIds = append(userIds, u.ID)
		}
	}
	if err := model.DB.Where("id in (?) AND user_id in (?)", ids, userIds).Delete(&model.Exchange{}).Error; err != nil {
		resp.Message = fmt.Sprint(err)
	} else {
		resp.Success = true
	}
	return
}

type exchangeHandler struct {
	*iris.Context
}

// Get /exchange
func (c exchangeHandler) Get() {
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
	id := c.URLParam("id")
	if id != "" && id != "0" {
		td, err := model.GetTrader(self, id)
		if err != nil {
			resp["msg"] = fmt.Sprint(err)
			c.JSON(iris.StatusOK, resp)
			return
		}
		if self, err = model.GetUserByID(td.UserID); err != nil {
			resp["msg"] = fmt.Sprint(err)
			c.JSON(iris.StatusOK, resp)
			return
		}
	}
	exchanges, err := model.GetExchanges(self)
	if err != nil {
		resp["msg"] = fmt.Sprint(err)
		c.JSON(iris.StatusOK, resp)
		return
	}
	resp["success"] = true
	resp["data"] = exchanges
	c.JSON(iris.StatusOK, resp)
}

// Post /exchange
func (c exchangeHandler) Post() {
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
	req := model.Exchange{}
	if err := c.ReadJSON(&req); err != nil {
		resp["msg"] = fmt.Sprint(err)
		c.JSON(iris.StatusOK, resp)
		return
	}
	if req.ID > 0 {
		exchange := model.Exchange{}
		if err := model.DB.First(&exchange, req.ID).Error; err != nil {
			resp["msg"] = fmt.Sprint(err)
			c.JSON(iris.StatusOK, resp)
			return
		}
		exchange.Name = req.Name
		exchange.Type = req.Type
		exchange.AccessKey = req.AccessKey
		exchange.SecretKey = req.SecretKey
		if err := model.DB.Save(&exchange).Error; err != nil {
			resp["msg"] = fmt.Sprint(err)
			c.JSON(iris.StatusOK, resp)
			return
		}
		resp["success"] = true
		c.JSON(iris.StatusOK, resp)
		return
	}
	req.UserID = self.ID
	if err := model.DB.Create(&req).Error; err != nil {
		resp["msg"] = fmt.Sprint(err)
		c.JSON(iris.StatusOK, resp)
		return
	}
	resp["success"] = true
	c.JSON(iris.StatusOK, resp)
}

// Delete /exchange
func (c exchangeHandler) Delete() {
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
	exchange, err := model.GetExchange(self, id)
	if err != nil {
		db.Rollback()
		resp["msg"] = fmt.Sprint(err)
		c.JSON(iris.StatusOK, resp)
		return
	}
	if err := db.Delete(&exchange).Error; err != nil {
		db.Rollback()
		resp["msg"] = fmt.Sprint(err)
		c.JSON(iris.StatusOK, resp)
		return
	}
	if err := db.Where("exchange_id = ?", id).Delete(&model.TraderExchange{}).Error; err != nil {
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
