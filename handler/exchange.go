package handler

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris"
	"github.com/miaolz123/samaritan/model"
)

type exchange struct{}

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
