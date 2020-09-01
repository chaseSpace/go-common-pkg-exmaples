package main

import (
	"log"
	"testing"
)

/*
gorm接收聚合函数结果
*/

func TestSqlFunc(t *testing.T) {
	startDB()
	defer closeDB()
	db.AutoMigrate(Student{})

	var ss []int // 必须是slice才能接收

	// 这种也是可以的
	//err := db.Raw("select sum(id) as sum from student where id > 0").Pluck("sum", &ss).Error

	err := db.Debug().Model(Student{}).Select("sum(id) as sum").Where("id > 0").Pluck("sum", &ss).Error

	log.Println(err, ss) // nil, []120
}
