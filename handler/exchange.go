package handler

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris"
	"github.com/miaolz123/samaritan/model"
)

type exchangeHandler struct {
	*iris.Context
}

// Get /exchange
func (c exchangeHandler) Get() {
	self, err := model.GetUser(jwtmid.Get(c.Context).Claims.(jwt.MapClaims)["sub"])
	if err != nil {
		c.Error(fmt.Sprint(err), iris.StatusServiceUnavailable)
		return
	}
	exchanges, err := model.GetExchanges(self)
	if err != nil {
		c.Error(fmt.Sprint(err), iris.StatusServiceUnavailable)
		return
	}
	c.JSON(iris.StatusOK, exchanges)
}

// Post /exchange
func (c exchangeHandler) Post() {
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
	req := model.Exchange{}
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

// Put /exchange
func (c exchangeHandler) Put() {
	db, err := model.NewOrm()
	if err != nil {
		c.Error(fmt.Sprint(err), iris.StatusServiceUnavailable)
		return
	}
	req := model.Exchange{}
	if err := c.ReadJSON(&req); err != nil {
		c.Error(fmt.Sprint(err), iris.StatusBadRequest)
		return
	}
	exchange := model.Exchange{}
	if err := db.First(&exchange, req.ID).Error; err != nil {
		c.Error(fmt.Sprint(err), iris.StatusServiceUnavailable)
		return
	}
	exchange.Name = req.Name
	exchange.Type = req.Type
	exchange.AccessKey = req.AccessKey
	exchange.SecretKey = req.SecretKey
	if err := db.Save(&exchange).Error; err != nil {
		c.Error(fmt.Sprint(err), iris.StatusServiceUnavailable)
		return
	}
	c.JSON(iris.StatusOK, exchange)
}
