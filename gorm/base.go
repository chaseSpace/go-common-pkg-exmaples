package main

import (
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

func startDB() {
	var err error
	//  "user:pass@(host:port)/DBNAME?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open("mysql", GetMysqlUri())
	if err != nil {
		panic(err)
	}
	_ = db
}

func closeDB() {
	db.Close()
}
