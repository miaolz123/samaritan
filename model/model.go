package model

import (
	"github.com/jinzhu/gorm"
	// for Sqlite
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// TraderMap ...
var TraderMap = make(map[uint]*Trader)

func init() {
	var err error
	db, err := gorm.Open("sqlite3", "data.db")
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&User{}, &Exchange{}, &Strategy{}, &TraderExchange{}, &Trader{}, &Log{})
	users := []User{}
	db.Find(&users)
	if len(users) == 0 {
		admin := User{
			Name:     "admin",
			Password: "admin",
			Level:    99,
		}
		if err := db.Create(&admin).Error; err != nil {
			panic(err)
		}
	}
}

// NewOrm ...
func NewOrm() (*gorm.DB, error) {
	return gorm.Open("sqlite3", "data.db")
}
