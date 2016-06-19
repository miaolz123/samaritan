package main

import (
	"runtime"
	"time"

	"github.com/go-ini/ini"
	"github.com/miaolz123/samaritan/api"
)

func init() {
	runtime.LockOSThread()
}

func main() {
	cfg, _ := ini.Load("app.ini")
	section := cfg.Section("test")
	opt := api.Option{
		Type:      section.Key("Type").MustString(""),
		AccessKey: section.Key("AccessKey").MustString(""),
		SecretKey: section.Key("SecretKey").MustString(""),
		MainStock: section.Key("MainStock").MustString(""),
	}
	opts := []api.Option{opt}
	scr := `exchange.Log(exchange.SetMainStock(LTC));
	if (exchange.GetAccount()) exchange.Log(exchange.GetAccount().Net);
	exchange.Buy(BTC,5050,0.4);
	exchange.Sell(BTC,5100,0.4);
	var acc = exchange.GetOrders(BTC);
	exchange.Log(212121, acc);
	if (acc) {
		Log(exchange.GetOrder(acc[0]));
		exchange.CancelOrder(acc[0]);
	} else Log("all done");
	Log(exchange.GetTicker(BTC).Mid);
	Log(exchange.GetRecords(BTC, M, 9999).length);
	while (true) {
		var rs = exchange.GetRecords(BTC, M, 9999);
		Log(rs.length);
		Log(rs[rs.length-3].Time)
		Log(rs[rs.length-2].Time)
		Log(rs[rs.length-1].Time)
		Log(rs[rs.length-1])
		break;
		Sleep(30000);
	}`
	r := api.New(opts, "ceshi001", scr)
	r.Run()
	time.Sleep(time.Second * 2)
	r.Stop()
}
