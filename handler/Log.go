package handler

import (
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris"
	"github.com/miaolz123/samaritan/model"
)

type pagination struct {
	Current  int
	PageSize int
	Total    int
}

type filters struct {
	Type         []string
	ExchangeType []string
}

type logsReq struct {
	Trader     model.Trader
	Pagination pagination
	Filters    filters
}

// Post /logs
func logs(c *iris.Context) {
	resp := iris.Map{
		"success": false,
		"msg":     "",
	}
	self, err := model.GetUser(jwtmid.Get(c).Claims.(jwt.MapClaims)["sub"])
	if err != nil {
		resp["msg"] = fmt.Sprint(err)
		c.JSON(iris.StatusOK, resp)
		return
	}
	req := logsReq{}
	if err := c.ReadJSON(&req); err != nil {
		resp["msg"] = fmt.Sprint(err)
		c.JSON(iris.StatusOK, resp)
		return
	}
	if req.Pagination.PageSize <= 0 || req.Pagination.PageSize > 100 {
		req.Pagination.PageSize = 20
	}
	if req.Trader, err = model.GetTrader(self, req.Trader.ID); err != nil {
		resp["msg"] = fmt.Sprint(err)
		c.JSON(iris.StatusOK, resp)
		return
	}
	raw := fmt.Sprintf("SELECT COUNT(*) total FROM logs WHERE trader_id = '%v'", req.Trader.ID)
	if len(req.Filters.ExchangeType) > 0 {
		raw += fmt.Sprintf(" AND exchange_type IN (%v)", strings.Join(req.Filters.ExchangeType, ","))
	}
	if len(req.Filters.Type) > 0 {
		raw += fmt.Sprintf(" AND type IN (%v)", strings.Join(req.Filters.Type, ","))
	}
	raw += " ORDER BY timestamp"
	total := struct {
		Total int64
	}{}
	if err := model.DB.Raw(raw).Scan(&total).Error; err != nil {
		resp["msg"] = fmt.Sprint(err)
		c.JSON(iris.StatusOK, resp)
		return
	}
	raw = strings.Replace(raw, "COUNT(*) total", "*", 1)
	raw += fmt.Sprintf(" DESC LIMIT %v OFFSET %v", req.Pagination.PageSize, req.Pagination.PageSize*(req.Pagination.Current-1))
	logs := []model.Log{}
	if err := model.DB.Raw(raw).Scan(&logs).Error; err != nil {
		resp["msg"] = fmt.Sprint(err)
		c.JSON(iris.StatusOK, resp)
		return
	}
	for i, l := range logs {
		logs[i].Time = time.Unix(l.Timestamp, 0).Format("01/02 15:04:05")
	}
	resp["success"] = true
	resp["total"] = total.Total
	resp["data"] = logs
	c.JSON(iris.StatusOK, resp)
}
