package model

import (
	"log"
	"strings"
	"time"

	"github.com/hprose/hprose-golang/io"
	"github.com/jinzhu/gorm"
	// for db SQL
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/miaolz123/samaritan/config"
)

var (
	// DB Database
	DB     *gorm.DB
	dbType = config.String("dbtype")
	dbURL  = config.String("dburl")
)

func init() {
	io.Register((*User)(nil), "User", "json")
	io.Register((*Exchange)(nil), "Exchange", "json")
	var err error
	DB, err = gorm.Open(strings.ToLower(dbType), dbURL)
	if err != nil {
		log.Printf("Connect to %v database error: %v\n", dbType, err)
		dbType = "sqlite3"
		dbURL = "custom/data.db"
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
			Username: "admin",
			Password: "admin",
			Level:    99,
		}
		if err := DB.Create(&admin).Error; err != nil {
			log.Fatalln("Create admin error:", err)
		}
	}
	DB.LogMode(false)
	go ping()
}

func ping() {
	for {
		if err := DB.Exec("SELECT 1").Error; err != nil {
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
