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
	strategies, err := model.GetStrategies(self)
	if err != nil {
		resp["msg"] = fmt.Sprint(err)
		c.JSON(iris.StatusOK, resp)
		return
	}
	resp["success"] = true
	resp["data"] = strategies
	c.JSON(iris.StatusOK, resp)
}

// Post /strategy
func (c strategyHandler) Post() {
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
	req := model.Strategy{}
	if err := c.ReadJSON(&req); err != nil {
		resp["msg"] = fmt.Sprint(err)
		c.JSON(iris.StatusOK, resp)
		return
	}
	if req.ID > 0 {
		strategy, err := model.GetStrategy(self, req.ID)
		if err != nil {
			resp["msg"] = fmt.Sprint(err)
			c.JSON(iris.StatusOK, resp)
			return
		}
		strategy.Name = req.Name
		strategy.Description = req.Description
		strategy.Script = req.Script
		if err := model.DB.Save(&strategy).Error; err != nil {
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

// Delete /strategy
func (c strategyHandler) Delete() {
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
	strategy, err := model.GetStrategy(self, id)
	if err != nil {
		db.Rollback()
		resp["msg"] = fmt.Sprint(err)
		c.JSON(iris.StatusOK, resp)
		return
	}
	if err := db.Delete(&strategy).Error; err != nil {
		db.Rollback()
		resp["msg"] = fmt.Sprint(err)
		c.JSON(iris.StatusOK, resp)
		return
	}
	if err := db.Exec(`UPDATE traders SET strategy_id = 0 WHERE strategy_id = ?`, id).Error; err != nil {
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
