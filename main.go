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
	}
	opts := []api.Option{opt}
	scr := "exchange.Log('Net: ',exchange.GetAccount().Net)"
	api.Run(opts, scr)
}
