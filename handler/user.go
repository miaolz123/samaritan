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
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
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
