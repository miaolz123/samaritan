package handler

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris"
	"github.com/miaolz123/samaritan/model"
)

type strategyHandler struct {
	*iris.Context
}

// Get /strategy
func (c strategyHandler) Get() {
	self, err := model.GetUser(jwtmid.Get(c.Context).Claims.(jwt.MapClaims)["sub"])
	if err != nil {
		c.Error(fmt.Sprint(err), iris.StatusServiceUnavailable)
		return
	}
	strategies, err := model.GetStrategies(self)
	if err != nil {
		c.Error(fmt.Sprint(err), iris.StatusServiceUnavailable)
		return
	}
	c.JSON(iris.StatusOK, strategies)
}

// Post /strategy
func (c strategyHandler) Post() {
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
	req := model.Strategy{}
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

// Put /strategy
func (c strategyHandler) Put() {
	db, err := model.NewOrm()
	if err != nil {
		c.Error(fmt.Sprint(err), iris.StatusServiceUnavailable)
		return
	}
	req := model.Strategy{}
	if err := c.ReadJSON(&req); err != nil {
		c.Error(fmt.Sprint(err), iris.StatusBadRequest)
		return
	}
	strategy := model.Strategy{}
	if err := db.First(&strategy, req.ID).Error; err != nil {
		c.Error(fmt.Sprint(err), iris.StatusServiceUnavailable)
		return
	}
	strategy.Name = req.Name
	strategy.Description = req.Description
	strategy.Script = req.Script
	if err := db.Save(&strategy).Error; err != nil {
		c.Error(fmt.Sprint(err), iris.StatusServiceUnavailable)
		return
	}
	c.JSON(iris.StatusOK, strategy)
}
