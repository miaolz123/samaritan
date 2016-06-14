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
    function ceshi() { Sleep(3000); Log("sleep ..."); }
	var acc = exchange.Sell("BTC",4950,50,"test");
    while (true) { ceshi(); }
	if (acc) exchange.Log(212121, acc);
	Log(111);`
	r := api.New(opts, "ceshi001", scr)
	r.Run()
	time.Sleep(time.Second * 10)
	r.Stop()
}
