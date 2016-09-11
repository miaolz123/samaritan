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
	self, err := model.GetUser(jwtmid.Get(c.Context).Claims.(jwt.MapClaims)["sub"])
	if err != nil {
		c.Error(fmt.Sprint(err), iris.StatusServiceUnavailable)
		return
	}
	users, err := model.GetUsers(self)
	if err != nil {
		c.Error(fmt.Sprint(err), iris.StatusServiceUnavailable)
		return
	}
	c.JSON(iris.StatusOK, users)
}

// GetBy /user/:name
func (c userHandler) GetBy(name string) {
	self, err := model.GetUser(jwtmid.Get(c.Context).Claims.(jwt.MapClaims)["sub"])
	if err != nil {
		c.Error(fmt.Sprint(err), iris.StatusServiceUnavailable)
		return
	}
	user, err := model.GetUser(name)
	if err != nil {
		c.Error(fmt.Sprint(err), iris.StatusServiceUnavailable)
		return
	}
	if self.Level > user.Level || self.Name == user.Name {
		c.JSON(iris.StatusOK, user)
		return
	}
	c.Error("", iris.StatusForbidden)
}

// Post /user
func (c userHandler) Post() {
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
	req := struct {
		Name     string
		Password string
		Level    int
	}{}
	if err := c.ReadJSON(&req); err != nil {
		c.Error(fmt.Sprint(err), iris.StatusBadRequest)
		return
	}
	count := int64(0)
	if err := db.Model(&model.User{}).Where("name = ?", req.Name).Count(&count).Error; err != nil {
		c.Error(fmt.Sprint(err), iris.StatusServiceUnavailable)
		return
	}
	if count > 0 {
		c.Error("Name exists", iris.StatusBadRequest)
		return
	}
	user := model.User{
		Name:     req.Name,
		Password: req.Password,
		Level:    req.Level,
	}
	if user.Level >= self.Level {
		user.Level = self.Level - 1
	}
	if err := db.Create(&user).Error; err != nil {
		c.Error(fmt.Sprint(err), iris.StatusServiceUnavailable)
		return
	}
	c.JSON(iris.StatusOK, user)
}
