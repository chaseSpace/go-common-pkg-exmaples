package main

import (
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

func startDB() {
	var err error
	//  "user:pass@(host:port)/DBNAME?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open("mysql", "root:123@tcp(192.168.31.11:3306)/adminbg?charset=utf8&parseTime=True&loc=Local&timeout=10000ms")
	if err != nil {
		panic(err)
	}
	_ = db
}

func closeDB() {
	db.Close()
}
