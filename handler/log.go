package handler

import (
	"fmt"

	"github.com/hprose/hprose-golang/rpc"
	"github.com/miaolz123/samaritan/constant"
	"github.com/miaolz123/samaritan/model"
)

type logger struct{}

type pagination struct {
	Current  int64
	PageSize int64
	Total    int64
}

type filters struct {
	Type         []string
	ExchangeType []string
}

func (logger) List(trader model.Trader, pagination pagination, filters filters, ctx rpc.Context) (resp response) {
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
	if trader, err = self.GetTrader(trader.ID); err != nil {
		resp.Message = fmt.Sprint(err)
		return
	}
	total, logs, err := self.ListLog(trader.ID, pagination.PageSize, pagination.Current)
	if err != nil {
		resp.Message = fmt.Sprint(err)
		return
	}
	resp.Data = struct {
		Total int64
		List  []model.Log
	}{
		Total: total,
		List:  logs,
	}
	resp.Success = true
	return
}

// // Post /logs
// func logs(c *iris.Context) {
// 	resp := iris.Map{
// 		"success": false,
// 		"msg":     "",
// 	}
// 	self, err := model.GetUser(jwtmid.Get(c).Claims.(jwt.MapClaims)["sub"])
// 	if err != nil {
// 		resp["msg"] = fmt.Sprint(err)
// 		c.JSON(iris.StatusOK, resp)
// 		return
// 	}
// 	req := logsReq{}
// 	if err := c.ReadJSON(&req); err != nil {
// 		resp["msg"] = fmt.Sprint(err)
// 		c.JSON(iris.StatusOK, resp)
// 		return
// 	}
// 	if req.Pagination.PageSize <= 0 || req.Pagination.PageSize > 100 {
// 		req.Pagination.PageSize = 20
// 	}
// 	if req.Trader, err = model.GetTrader(self, req.Trader.ID); err != nil {
// 		resp["msg"] = fmt.Sprint(err)
// 		c.JSON(iris.StatusOK, resp)
// 		return
// 	}
// 	raw := fmt.Sprintf("SELECT COUNT(*) total FROM logs WHERE trader_id = '%v'", req.Trader.ID)
// 	if len(req.Filters.ExchangeType) > 0 {
// 		raw += fmt.Sprintf(" AND exchange_type IN (%v)", strings.Join(req.Filters.ExchangeType, ","))
// 	}
// 	if len(req.Filters.Type) > 0 {
// 		raw += fmt.Sprintf(" AND type IN (%v)", strings.Join(req.Filters.Type, ","))
// 	}
// 	total := struct {
// 		Total int64
// 	}{}
// 	if err := model.DB.Raw(raw).Scan(&total).Error; err != nil {
// 		resp["msg"] = fmt.Sprint(err)
// 		c.JSON(iris.StatusOK, resp)
// 		return
// 	}
// 	raw = strings.Replace(raw, "COUNT(*) total", "*", 1)
// 	raw += fmt.Sprintf(" ORDER BY timestamp DESC, id DESC LIMIT %v OFFSET %v",
// 		req.Pagination.PageSize, req.Pagination.PageSize*(req.Pagination.Current-1))
// 	logs := []model.Log{}
// 	if err := model.DB.Raw(raw).Scan(&logs).Error; err != nil {
// 		resp["msg"] = fmt.Sprint(err)
// 		c.JSON(iris.StatusOK, resp)
// 		return
// 	}
// 	loc, err := time.LoadLocation(config.String("logstimezone"))
// 	if err != nil || loc == nil {
// 		loc = time.Local
// 	}
// 	for i, l := range logs {
// 		logs[i].Time = time.Unix(l.Timestamp, 0).In(loc).Format("01/02 15:04:05")
// 	}
// 	resp["success"] = true
// 	resp["total"] = total.Total
// 	resp["data"] = logs
// 	c.JSON(iris.StatusOK, resp)
// }

// // Delete /logs
// func logsDelete(c *iris.Context) {
// 	resp := iris.Map{
// 		"success": false,
// 		"msg":     "",
// 	}
// 	self, err := model.GetUser(jwtmid.Get(c).Claims.(jwt.MapClaims)["sub"])
// 	if err != nil {
// 		resp["msg"] = fmt.Sprint(err)
// 		c.JSON(iris.StatusOK, resp)
// 		return
// 	}
// 	td := model.Trader{}
// 	if td, err = model.GetTrader(self, c.URLParam("id")); err != nil {
// 		resp["msg"] = fmt.Sprint(err)
// 		c.JSON(iris.StatusOK, resp)
// 		return
// 	}
// 	timeType := c.URLParam("type")
// 	minTimestamp := int64(0)
// 	switch timeType {
// 	case constant.LastTime:
// 		if t := trader.Executor[td.ID]; t != nil {
// 			minTimestamp = t.LastRunAt.Unix()
// 		} else {
// 			resp["msg"] = "Not found running trader"
// 			c.JSON(iris.StatusOK, resp)
// 			return
// 		}
// 	case constant.Day:
// 		minTimestamp = time.Now().AddDate(0, 0, -1).Unix()
// 	case constant.Week:
// 		minTimestamp = time.Now().AddDate(0, 0, -7).Unix()
// 	case constant.Month:
// 		minTimestamp = time.Now().AddDate(0, -1, 0).Unix()
// 	}
// 	if err := model.DB.Where("trader_id = ?", td.ID).Where("timestamp < ?", minTimestamp).Delete(&model.Log{}).Error; err != nil {
// 		resp["msg"] = fmt.Sprint(err)
// 		c.JSON(iris.StatusOK, resp)
// 		return
// 	}
// 	resp["success"] = true
// 	c.JSON(iris.StatusOK, resp)
// }

// // Get /profits
// func profits(c *iris.Context) {
// 	resp := iris.Map{
// 		"success": false,
// 		"msg":     "",
// 	}
// 	self, err := model.GetUser(jwtmid.Get(c).Claims.(jwt.MapClaims)["sub"])
// 	if err != nil {
// 		resp["msg"] = fmt.Sprint(err)
// 		c.JSON(iris.StatusOK, resp)
// 		return
// 	}
// 	td := model.Trader{}
// 	if td, err = model.GetTrader(self, c.URLParam("id")); err != nil {
// 		resp["msg"] = fmt.Sprint(err)
// 		c.JSON(iris.StatusOK, resp)
// 		return
// 	}
// 	logs := []model.Log{}
// 	if err := model.DB.Where("trader_id = ?", td.ID).Where("type = 1").Order("timestamp").Order("id").Find(&logs).Error; err != nil {
// 		resp["msg"] = fmt.Sprint(err)
// 		c.JSON(iris.StatusOK, resp)
// 		return
// 	}
// 	amount, _ := c.URLParamInt("amount")
// 	if amount == 0 {
// 		amount = 100
// 	}
// 	if len(logs) < amount {
// 		amount = len(logs)
// 	}
// 	data := []model.Log{}
// 	loc, err := time.LoadLocation(config.String("logstimezone"))
// 	if err != nil || loc == nil {
// 		loc = time.Local
// 	}
// 	for i := 1; i <= amount; i++ {
// 		index := i*len(logs)/amount - 1
// 		data = append(data, model.Log{
// 			Time:   time.Unix(logs[index].Timestamp, 0).In(loc).Format("01/02 15:04:05"),
// 			Amount: logs[index].Amount,
// 		})
// 	}
// 	resp["data"] = data
// 	resp["success"] = true
// 	c.JSON(iris.StatusOK, resp)
// }

// // Get /status
// func status(c *iris.Context) {
// 	resp := iris.Map{
// 		"success": false,
// 		"msg":     "",
// 	}
// 	self, err := model.GetUser(jwtmid.Get(c).Claims.(jwt.MapClaims)["sub"])
// 	if err != nil {
// 		resp["msg"] = fmt.Sprint(err)
// 		c.JSON(iris.StatusOK, resp)
// 		return
// 	}
// 	td := model.Trader{}
// 	if td, err = model.GetTrader(self, c.URLParam("id")); err != nil {
// 		resp["msg"] = fmt.Sprint(err)
// 		c.JSON(iris.StatusOK, resp)
// 		return
// 	}
// 	resp["data"] = trader.GetStatus(td.ID)
// 	resp["success"] = true
// 	c.JSON(iris.StatusOK, resp)
// }
