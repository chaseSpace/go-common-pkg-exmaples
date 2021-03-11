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

	var ss []*int // 必须是slice才能接收, sum的话得是指针类型，因为没有记录的时候为null==nil
	// 这种也是可以的
	//err := db.Raw("select sum(id) as sum from student where id > 0").Pluck("sum", &ss).Error

	db.Debug()
	err := db.Debug().Model(Student{}).Select("sum(id) as sum").Where("id > 0 and vc=b'1'").Pluck("sum", &ss).Error

	log.Println(err, ss) // nil, []120
}
