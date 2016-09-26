package trader

import (
	"time"

	"github.com/miaolz123/conver"
	"github.com/miaolz123/samaritan/constant"
	"github.com/miaolz123/samaritan/model"
)

type global struct {
	trader model.Trader
}

func (g global) Sleep(intervals ...interface{}) {
	interval := int64(0)
	if len(intervals) > 0 {
		interval = conver.Int64Must(intervals[0])
	}
	if interval > 0 {
		time.Sleep(time.Duration(interval * 1000000))
	}
}

func (g global) Log(msgs ...interface{}) {
	g.trader.Logger.Log(constant.INFO, 0.0, 0.0, msgs...)
}

func (g global) LogProfit(msgs ...interface{}) {
	profit := 0.0
	if len(msgs) > 0 {
		profit = conver.Float64Must(msgs[0])
	}
	g.trader.Logger.Log(constant.PROFIT, 0.0, profit, msgs[1:]...)
}
