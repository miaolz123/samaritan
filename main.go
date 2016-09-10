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
	Log(exchanges[1].GetRecords(BTC, M, 10));
	`
	r := api.New(opts, "ceshi001", scr)
	r.Run()
	time.Sleep(time.Second * 2)
	r.Stop()
}
