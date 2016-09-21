package handler

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris"
	"github.com/miaolz123/samaritan/model"
)

type userHandler struct {
	*iris.Context
}

// Post /login
func userLogin(c *iris.Context) {
	req := struct {
		Name     string
		Password string
	}{}
	resp := struct {
		Token string
	}{}
	err := c.ReadJSON(&req)
	if err != nil {
		c.Error(fmt.Sprint(err), iris.StatusBadRequest)
		return
	}
	if req.Name == "" || req.Password == "" {
		c.Error("Name and Password can not be empty", iris.StatusBadRequest)
		return
	}
	db, err := model.NewOrm()
	if err != nil {
		c.Error(fmt.Sprint(err), iris.StatusServiceUnavailable)
		return
	}
	user := model.User{}
	if err := db.Where(&model.User{Name: req.Name, Password: req.Password}).First(&user).Error; err != nil {
		c.Error("Name or Password wrong", iris.StatusUnauthorized)
		return
	}
	claims := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		Subject:   user.Name,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	if resp.Token, err = token.SignedString(signKey); err != nil {
		c.Error(fmt.Sprint(err), iris.StatusInternalServerError)
		return
	}
	c.JSON(iris.StatusOK, resp)
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
	db, err := model.NewOrm()
	if err != nil {
		resp["msg"] = fmt.Sprint(err)
		c.JSON(iris.StatusOK, resp)
		return
	}
	self, err := model.GetUser(jwtmid.Get(c.Context).Claims.(jwt.MapClaims)["sub"])
	if err != nil {
		resp["msg"] = fmt.Sprint(err)
		c.JSON(iris.StatusOK, resp)
		return
	}
	req := model.User{}
	if err := c.ReadJSON(&req); err != nil {
		resp["msg"] = fmt.Sprint(err)
		c.JSON(iris.StatusOK, resp)
		return
	}
	if req.ID > 0 {
		user := model.User{}
		if err := db.First(&user, req.ID).Error; err != nil {
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
		if err := db.Save(&user).Error; err != nil {
			resp["msg"] = fmt.Sprint(err)
			c.JSON(iris.StatusOK, resp)
			return
		}
		resp["success"] = true
		c.JSON(iris.StatusOK, resp)
		return
	}
	if req.Level >= self.Level {
		req.Level = self.Level - 1
	}
	if err := db.Create(&req).Error; err != nil {
		resp["msg"] = fmt.Sprint(err)
		c.JSON(iris.StatusOK, resp)
		return
	}
	resp["success"] = true
	c.JSON(iris.StatusOK, resp)
}
