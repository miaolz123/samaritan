package config

import (
	"log"
	"strings"

	"github.com/go-ini/ini"
)

var confs = make(map[string]string)

func init() {
	conf, err := ini.InsensitiveLoad("custom/config.ini")
	if err != nil {
		conf, err = ini.InsensitiveLoad("config.ini")
		if err != nil {
			log.Fatalln("Load config.ini error:", err)
		}
	}
	keys := conf.Section("").KeyStrings()
	for _, k := range keys {
		confs[k] = conf.Section("").Key(k).String()
	}
	if confs["logstimezone"] == "" {
		confs["logstimezone"] = "Local"
	}
}

// String ...
func String(key string) string {
	return confs[strings.ToLower(key)]
}
