package main

import (
	"log"
	"time"

	"github.com/go-ini/ini"
	"github.com/miaolz123/samaritan/api"
)

func main() {
	opts := []api.Option{}
	cfg, _ := ini.Load("appMy.ini")
	for _, s := range cfg.Sections() {
		opt := api.Option{}
		if err := s.MapTo(&opt); err != nil {
			log.Panicln(err)
		}
		if opt.Type != "" {
			opts = append(opts, opt)
		}
	}
	scr := `
	var orders = exchanges[1].GetOrders(BTC);
	Log(exchanges[1].CancelOrder(orders[0]));
	Log(exchanges[1].Sell(BTC, 0, 0.02));
	`
	r := api.New(opts, "ceshi001", scr)
	r.Run()
	time.Sleep(time.Second * 2)
	r.Stop()
}
