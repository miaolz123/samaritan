package handler

import (
	"fmt"

	"github.com/hprose/hprose-golang/rpc"
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
func (exchange) List(size, page int64, order string, ctx rpc.Context) (resp response) {
	username := ctx.GetString("username")
	if username == "" {
		resp.Message = constant.ErrAuthorizationError
		return
	}
	self, err := model.GetUser(username)
	if err != nil {
		resp.Message = fmt.Sprint(err)
		return
	}
	total, exchanges, err := self.ExchangeList(size, page, order)
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
		resp.Message = constant.ErrAuthorizationError
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
		resp.Message = constant.ErrAuthorizationError
		return
	}
	self, err := model.GetUser(username)
	if err != nil {
		resp.Message = fmt.Sprint(err)
		return
	}
	userIds := []int64{}
	_, users, err := self.UserList(-1, 1, "id")
	if err != nil {
		resp.Message = fmt.Sprint(err)
		return
	}
	for _, u := range users {
		userIds = append(userIds, u.ID)
	}
	if err := model.DB.Where("id in (?) AND user_id in (?)", ids, userIds).Delete(&model.Exchange{}).Error; err != nil {
		resp.Message = fmt.Sprint(err)
	} else {
		resp.Success = true
	}
	return
}
