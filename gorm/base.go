package main

import "github.com/jinzhu/gorm"

var db *gorm.DB

func startDB() {
	var err error
	db, err = gorm.Open("mysql", "test_u:1918ddk@(114.115.216.44:33061)/test?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	_ = db
}

func closeDB() {
	db.Close()
}
