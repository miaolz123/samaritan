package handler

import (
	"fmt"

	"github.com/hprose/hprose-golang/rpc"
	"github.com/miaolz123/samaritan/constant"
	"github.com/miaolz123/samaritan/model"
	"github.com/miaolz123/samaritan/trader"
)

type runner struct{}

// List
func (runner) List(algorithmID int64, ctx rpc.Context) (resp response) {
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
	traders, err := self.TraderList(algorithmID)
	if err != nil {
		resp.Message = fmt.Sprint(err)
		return
	}
	resp.Data = traders
	resp.Success = true
	return
}

// Put
func (runner) Put(req model.Trader, ctx rpc.Context) (resp response) {
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
	db, err := model.NewOrm()
	if err != nil {
		resp.Message = fmt.Sprint(err)
		return
	}
	defer db.Close()
	db = db.Begin()
	runner := req
	if req.ID > 0 {
		if err := db.First(&runner, req.ID).Error; err != nil {
			db.Rollback()
			resp.Message = fmt.Sprint(err)
			return
		}
		runner.Name = req.Name
		runner.Environment = req.Environment
		rs, err := self.GetTraderExchanges(runner.ID)
		if err != nil {
			db.Rollback()
			resp.Message = fmt.Sprint(err)
			return
		}
		for i, r := range rs {
			if i >= len(req.Exchanges) {
				if err := db.Delete(&r).Error; err != nil {
					db.Rollback()
					resp.Message = fmt.Sprint(err)
					return
				}
				continue
			}
			if r.Exchange.ID == req.Exchanges[i].ID {
				continue
			}
			r.ExchangeID = req.Exchanges[i].ID
			if err := db.Save(&r).Error; err != nil {
				db.Rollback()
				resp.Message = fmt.Sprint(err)
				return
			}
		}
		for i, e := range req.Exchanges {
			if i < len(rs) {
				continue
			}
			r := model.TraderExchange{
				TraderID:   runner.ID,
				ExchangeID: e.ID,
			}
			if err := db.Create(&r).Error; err != nil {
				db.Rollback()
				resp.Message = fmt.Sprint(err)
				return
			}
		}
		if err := db.Save(&runner).Error; err != nil {
			db.Rollback()
			resp.Message = fmt.Sprint(err)
			return
		}
		if err := db.Commit().Error; err != nil {
			db.Rollback()
			resp.Message = fmt.Sprint(err)
			return
		}
		resp.Success = true
		return
	}
	req.UserID = self.ID
	if err := db.Create(&req).Error; err != nil {
		db.Rollback()
		resp.Message = fmt.Sprint(err)
		return
	}
	for _, e := range req.Exchanges {
		traderExchange := model.TraderExchange{
			TraderID:   req.ID,
			ExchangeID: e.ID,
		}
		if err := db.Create(&traderExchange).Error; err != nil {
			db.Rollback()
			resp.Message = fmt.Sprint(err)
			return
		}
	}
	if err := db.Commit().Error; err != nil {
		db.Rollback()
		resp.Message = fmt.Sprint(err)
		return
	}
	resp.Success = true
	return
}

// Run
func (runner) Run(req model.Trader, ctx rpc.Context) (resp response) {
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
	user, err := model.GetUserByID(req.UserID)
	if err != nil {
		resp.Message = fmt.Sprint(err)
		return
	}
	if user.Level > self.Level || user.ID != self.ID {
		resp.Message = constant.ErrInsufficientPermissions
		return
	}
	if err := trader.Run(trader.Global{Trader: req}); err != nil {
		resp.Message = fmt.Sprint(err)
		return
	}
	resp.Success = true
	return
}

// Stop
func (runner) Stop(req model.Trader, ctx rpc.Context) (resp response) {
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
	user, err := model.GetUserByID(req.UserID)
	if err != nil {
		resp.Message = fmt.Sprint(err)
		return
	}
	if user.Level > self.Level || user.ID != self.ID {
		resp.Message = constant.ErrInsufficientPermissions
		return
	}
	if err := trader.Stop(req.ID); err != nil {
		resp.Message = fmt.Sprint(err)
		return
	}
	resp.Success = true
	return
}
