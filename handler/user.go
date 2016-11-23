package handler

import (
	"fmt"

	"github.com/hprose/hprose-golang/rpc"
	"github.com/miaolz123/samaritan/model"
)

type user struct{}

// Login ...
func (user) Login(username, password string, ctx rpc.Context) (resp response) {
	user := model.User{
		Username: username,
		Password: password,
	}
	if user.Username == "" || user.Password == "" {
		resp.Message = "Username and Password can not be empty"
		return
	}
	if err := model.DB.Where(&user).First(&user).Error; err != nil {
		resp.Message = "Username or Password wrong"
		return
	}
	if resp.Data = makeToken(user.Username); resp.Data != "" {
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
		Username: username,
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
func (user) List(size, page int64, order string, ctx rpc.Context) (resp response) {
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
	total, users, err := self.UserList(size, page, order)
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

// Put ...
func (user) Put(req model.User, password string, ctx rpc.Context) (resp response) {
	username := ctx.GetString("username")
	if username == "" {
		resp.Message = "Authorization wrong"
		return
	}
	if req.Username == "" {
		resp.Message = "Request data wrong"
		return
	}
	self, err := model.GetUser(username)
	if err != nil {
		resp.Message = fmt.Sprint(err)
		return
	}
	user := model.User{
		Username: req.Username,
		Level:    req.Level,
		Password: password,
	}
	if req.ID > 0 {
		if err := model.DB.First(&user, req.ID).Error; err != nil {
			resp.Message = fmt.Sprint(err)
			return
		}
		user.Level = req.Level
		if user.Level >= self.Level {
			if user.ID == self.ID {
				user.Level = self.Level
			} else {
				user.Level = self.Level - 1
			}
		}
		if password != "" {
			user.Password = password
		}
		if err := model.DB.Save(&user).Error; err != nil {
			resp.Message = fmt.Sprint(err)
			return
		}
		resp.Success = true
		return
	}
	if password == "" {
		resp.Message = "Password can't be empty"
		return
	}
	if user.Level >= self.Level {
		user.Level = self.Level - 1
	}
	if err := model.DB.Create(&user).Error; err != nil {
		resp.Message = fmt.Sprint(err)
	} else {
		resp.Success = true
	}
	return
}

// Delete ...
func (user) Delete(ids []int64, ctx rpc.Context) (resp response) {
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
	if err := model.DB.Where("id in (?) AND level < ?", ids, self.Level).Delete(&model.User{}).Error; err != nil {
		resp.Message = fmt.Sprint(err)
	} else {
		resp.Success = true
	}
	return
}
