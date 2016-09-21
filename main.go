package main

import (
	"log"

	"github.com/go-ini/ini"
	"github.com/miaolz123/samaritan/handler"
)

func main() {
	conf, err := ini.Load("config.ini")
	if err != nil {
		log.Fatalln("Load config.ini error:", err)
	}
	handler.Server.Listen(":" + conf.Section("").Key("ServerPort").String())
}
