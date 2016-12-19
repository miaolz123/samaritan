package handler

import (
	"fmt"

	"github.com/hprose/hprose-golang/rpc"
	"github.com/miaolz123/samaritan/constant"
	"github.com/miaolz123/samaritan/model"
)

type algorithm struct{}

// List ...
func (algorithm) List(size, page int64, order string, ctx rpc.Context) (resp response) {
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
	total, algorithms, err := self.AlgorithmList(size, page, order)
	if err != nil {
		resp.Message = fmt.Sprint(err)
		return
	}
	resp.Data = struct {
		Total int64
		List  []model.Algorithm
	}{
		Total: total,
		List:  algorithms,
	}
	resp.Success = true
	return
}

// Put
func (algorithm) Put(req model.Algorithm, ctx rpc.Context) (resp response) {
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
	algorithm := req
	if req.ID > 0 {
		if err := model.DB.First(&algorithm, req.ID).Error; err != nil {
			resp.Message = fmt.Sprint(err)
			return
		}
		algorithm.Name = req.Name
		algorithm.Description = req.Description
		algorithm.Script = req.Script
		algorithm.EvnDefault = req.EvnDefault
		if err := model.DB.Save(&algorithm).Error; err != nil {
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
func (algorithm) Delete(ids []int64, ctx rpc.Context) (resp response) {
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
	if err := model.DB.Where("id in (?) AND user_id in (?)", ids, userIds).Delete(&model.Algorithm{}).Error; err != nil {
		resp.Message = fmt.Sprint(err)
	} else {
		resp.Success = true
	}
	return
}
