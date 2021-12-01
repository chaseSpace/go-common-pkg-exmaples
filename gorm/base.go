package main

import (
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

func GetMysqlUri() string {
	//return "user:pass@(host:port)/DBNAME?charset=utf8mb4&parseTime=True&loc=Local"
	return "root:123@(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
}

func startDB() {
	var err error
	//url := "user:pass@(host:port)/DBNAME?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open("mysql", GetMysqlUri())
	if err != nil {
		panic(err)
	}
	_ = db
}

func closeDB() {
	db.Close()
}
