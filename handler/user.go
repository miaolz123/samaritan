package handler

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/hprose/hprose-golang/rpc"
	"github.com/kataras/iris"
	"github.com/miaolz123/samaritan/model"
	"github.com/miaolz123/samaritan/trader"
)

type user struct{}

// Login ...
func (user) Login(username, password string, ctx rpc.Context) (resp response) {
	user := model.User{
		Name:     username,
		Password: password,
	}
	if user.Name == "" || user.Password == "" {
		resp.Message = "Username and Password can not be empty"
		return
	}
	if err := model.DB.Where(&user).First(&user).Error; err != nil {
		resp.Message = "Username or Password wrong"
		return
	}
	if resp.Data = makeToken(user.Name); resp.Data != "" {
		resp.Success = true
	} else {
		resp.Message = "Make token error"
	}
	return
}

// Get ...
func (user) Get(_ string, ctx rpc.Context) (resp response) {
	username := ctx.GetString("username")
	if username == "" {
		resp.Message = "Authorization wrong"
		return
	}
	user := model.User{
		Name: username,
	}
	if err := model.DB.Where(&user).First(&user).Error; err != nil {
		resp.Message = "Authorization username wrong"
		return
	}
	resp.Data = user
	resp.Success = true
	return
}

// List ...
func (user) List(size, page int64, ctx rpc.Context) (resp response) {
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
	total, users, err := self.UserList(size, page)
	if err != nil {
		resp.Message = fmt.Sprint(err)
		return
	}
	resp.Data = struct {
		Total int64
		List  []model.User
	}{
		Total: total,
		List:  users,
	}
	resp.Success = true
	return
}

type userHandler struct {
	*iris.Context
}

// Post /login
func userLogin(c *iris.Context) {
	resp := iris.Map{
		"success": false,
		"msg":     "",
	}
	req := struct {
		Name     string
		Password string
	}{}
	if err := c.ReadJSON(&req); err != nil {
		resp["msg"] = fmt.Sprint(err)
		c.JSON(iris.StatusOK, resp)
		return
	}
	user := model.User{
		Name:     req.Name,
		Password: req.Password,
	}
	if user.Name == "" || user.Password == "" {
		resp["msg"] = "Username and Password can not be empty"
		c.JSON(iris.StatusOK, resp)
		return
	}
	if err := model.DB.Where(&user).First(&user).Error; err != nil {
		resp["msg"] = "Username or Password wrong"
		c.JSON(iris.StatusOK, resp)
		return
	}
	claims := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(),
		Subject:   user.Name,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	if t, err := token.SignedString(signKey); err != nil {
		resp["msg"] = "Username or Password wrong"
		c.JSON(iris.StatusOK, resp)
		return
	} else if t != "" {
		resp["success"] = true
		resp["data"] = t
	}
	c.JSON(iris.StatusOK, resp)
}

// Post /token
func token(c *iris.Context) {
	c.JSON(iris.StatusOK, iris.Map{})
}

// Get /user
func (c userHandler) Get() {
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
	order := c.URLParam("order")
	users, err := model.GetUsers(self, order)
	if err != nil {
		resp["msg"] = fmt.Sprint(err)
		c.JSON(iris.StatusOK, resp)
		return
	}
	resp["success"] = true
	resp["data"] = users
	c.JSON(iris.StatusOK, resp)
}

// Post /user
func (c userHandler) Post() {
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
	req := struct {
		ID       uint
		Name     string
		Password string
		Level    int
	}{}
	if err := c.ReadJSON(&req); err != nil {
		resp["msg"] = fmt.Sprint(err)
		c.JSON(iris.StatusOK, resp)
		return
	}
	user := model.User{
		Name:     req.Name,
		Password: req.Password,
		Level:    req.Level,
	}
	if req.ID > 0 {
		if err := model.DB.First(&user, req.ID).Error; err != nil {
			resp["msg"] = fmt.Sprint(err)
			c.JSON(iris.StatusOK, resp)
			return
		}
		user.Name = req.Name
		user.Level = req.Level
		if user.Level >= self.Level {
			if req.ID == self.ID {
				req.Level = self.Level
			} else {
				req.Level = self.Level - 1
			}
		}
		if req.Password != "" {
			user.Password = req.Password
		}
		if err := model.DB.Save(&user).Error; err != nil {
			resp["msg"] = fmt.Sprint(err)
			c.JSON(iris.StatusOK, resp)
			return
		}
		resp["success"] = true
		c.JSON(iris.StatusOK, resp)
		return
	}
	if user.Level >= self.Level {
		user.Level = self.Level - 1
	}
	if err := model.DB.Create(&user).Error; err != nil {
		resp["msg"] = fmt.Sprint(err)
		c.JSON(iris.StatusOK, resp)
		return
	}
	resp["success"] = true
	c.JSON(iris.StatusOK, resp)
}

// Delete /user
func (c userHandler) Delete() {
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
	user, err := model.GetUserByID(id)
	if err != nil {
		db.Rollback()
		resp["msg"] = fmt.Sprint(err)
		c.JSON(iris.StatusOK, resp)
		return
	}
	if user.Level >= self.Level {
		db.Rollback()
		resp["msg"] = "Insufficient permissions"
		c.JSON(iris.StatusOK, resp)
		return
	}
	if err := db.Unscoped().Delete(&user).Error; err != nil {
		db.Rollback()
		resp["msg"] = fmt.Sprint(err)
		c.JSON(iris.StatusOK, resp)
		return
	}
	exchanges := []model.Exchange{}
	if err := db.Find(&exchanges, "user_id = ?", id).Error; err != nil {
		db.Rollback()
		resp["msg"] = fmt.Sprint(err)
		c.JSON(iris.StatusOK, resp)
		return
	}
	if err := db.Unscoped().Where("user_id = ?", id).Delete(&model.Exchange{}).Error; err != nil {
		db.Rollback()
		resp["msg"] = fmt.Sprint(err)
		c.JSON(iris.StatusOK, resp)
		return
	}
	if len(exchanges) > 0 {
		exchangeIDs := []uint{}
		for _, e := range exchanges {
			exchangeIDs = append(exchangeIDs, e.ID)
		}
		if err := db.Unscoped().Where("exchange_id in (?)", exchangeIDs).Delete(&model.TraderExchange{}).Error; err != nil {
			db.Rollback()
			resp["msg"] = fmt.Sprint(err)
			c.JSON(iris.StatusOK, resp)
			return
		}
	}
	strategies := []model.Strategy{}
	if err := db.Find(&strategies, "user_id = ?", id).Error; err != nil {
		db.Rollback()
		resp["msg"] = fmt.Sprint(err)
		c.JSON(iris.StatusOK, resp)
		return
	}
	if err := db.Unscoped().Where("user_id = ?", id).Delete(&model.Strategy{}).Error; err != nil {
		db.Rollback()
		resp["msg"] = fmt.Sprint(err)
		c.JSON(iris.StatusOK, resp)
		return
	}
	if len(strategies) > 0 {
		strategyIDs := []uint{}
		for _, s := range strategies {
			strategyIDs = append(strategyIDs, s.ID)
		}
		if err := db.Exec(`UPDATE traders SET strategy_id = 0 WHERE strategy_id in (?)`, strategyIDs).Error; err != nil {
			db.Rollback()
			resp["msg"] = fmt.Sprint(err)
			c.JSON(iris.StatusOK, resp)
			return
		}
	}
	traders := []model.Trader{}
	if err := db.Find(&traders, "user_id = ?", id).Error; err != nil {
		db.Rollback()
		resp["msg"] = fmt.Sprint(err)
		c.JSON(iris.StatusOK, resp)
		return
	}
	if err := db.Unscoped().Where("user_id = ?", id).Delete(&model.Trader{}).Error; err != nil {
		db.Rollback()
		resp["msg"] = fmt.Sprint(err)
		c.JSON(iris.StatusOK, resp)
		return
	}
	if len(traders) > 0 {
		traderIDs := []uint{}
		for _, t := range traders {
			traderIDs = append(traderIDs, t.ID)
		}
		if err := db.Unscoped().Where("trader_id in (?)", traderIDs).Delete(&model.TraderExchange{}).Error; err != nil {
			db.Rollback()
			resp["msg"] = fmt.Sprint(err)
			c.JSON(iris.StatusOK, resp)
			return
		}
		if err := db.Unscoped().Where("trader_id in (?)", traderIDs).Delete(&model.Log{}).Error; err != nil {
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
	trader.Clean(user.ID)
	resp["success"] = true
	c.JSON(iris.StatusOK, resp)
}
