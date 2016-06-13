package main

import (
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
	scr := `exchange.Log(exchange.SetMainStock());
	var acc = exchange.Buy("BTC",-1,50,"test");
	exchange.Log(202020, typeof(acc));
	if (acc) exchange.Log(212121, acc);
	Log(111);`
	api.Run(opts, scr)
}
