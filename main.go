package main

import (
	"time"

	"github.com/go-ini/ini"
	"github.com/miaolz123/samaritan/api"
)

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
	scr := `exchange.Log(exchange.SetMainStock("LTC"));
	exchange.Log(exchange.GetAccount().Net);
	exchange.Buy("BTC",4650,1);
	var acc = exchange.GetOrders("BTC");
	exchange.Log(212121, acc);
	exchange.CancelOrder(acc);
	Log(111);`
	r := api.New(opts, "ceshi001", scr)
	r.Run()
	time.Sleep(time.Second * 2)
	r.Stop()
}
