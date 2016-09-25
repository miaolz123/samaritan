package model

import (
	"log"
	"strings"
	"time"

	"github.com/go-ini/ini"
	"github.com/jinzhu/gorm"
	// for data SQL
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var (
	// DB Database
	DB     *gorm.DB
	dbType string
	dbURL  string
)

func init() {
	conf, err := ini.Load("config.ini")
	if err != nil {
		log.Fatalln("Load config.ini error:", err)
	}
	dbType = conf.Section("").Key("DatabaseType").String()
	dbURL = conf.Section("").Key("DatabaseURL").String()
	DB, err = gorm.Open(strings.ToLower(dbType), dbURL)
	if err != nil {
		log.Printf("Connect to %v database error: %v\n", dbType, err)
		dbType = "sqlite3"
		dbURL = "data.db"
		DB, err = gorm.Open(dbType, dbURL)
		if err != nil {
			log.Fatalln("Connect to database error:", err)
		}
	}
	DB.AutoMigrate(&User{}, &Exchange{}, &Strategy{}, &TraderExchange{}, &Trader{}, &Log{})
	users := []User{}
	DB.Find(&users)
	if len(users) == 0 {
		admin := User{
			Name:     "admin",
			Password: "admin",
			Level:    99,
		}
		if err := DB.Create(&admin).Error; err != nil {
			log.Fatalln("Create admin error:", err)
		}
	}
	go ping()
}

func ping() {
	for {
		if err := DB.DB().Ping(); err != nil {
			log.Println("Database ping error:", err)
			if DB, err = gorm.Open(strings.ToLower(dbType), dbURL); err != nil {
				log.Println("Retry connect to database error:", err)
			}
		}
		time.Sleep(time.Minute)
	}
}

// NewOrm ...
func NewOrm() (*gorm.DB, error) {
	return gorm.Open(strings.ToLower(dbType), dbURL)
}
