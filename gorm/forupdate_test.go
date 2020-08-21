package main

import (
	"github.com/jinzhu/gorm"
	"log"
	"testing"
)

/*
gorm的事务
-------------
1. 如果tx对象在提交或回滚前被传递到函数之外，gorm会使用一个新的db连接，所以这个时候再去操作之前锁住的记录会死锁！
所以只能在一个函数内完成事务操作
*/

type Student struct {
	Id    int
	Name  string
	Score int
}

func (Student) TableName() string {
	return "student"
}

type conn struct {
	id int
}

func getTx(d *gorm.DB) *gorm.DB {
	return d.Begin()
}

// 测试锁住id=1的记录时，并去修改它, 结果是可以（同一个事务中）！
func ForUpdateTest(t *testing.T) {
	tx := getTx(db)
	//var dbConn = new(conn)
	//db.Exec("select connection_id()").Scan(dbConn)
	//log.Printf("connection_id-1:%d", dbConn.id)

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()
	if tx.Error != nil {
		panic(tx.Error)
	}

	var stu = &Student{}
	action := tx.Debug().Set("gorm:query_option", "FOR UPDATE").First(stu, "name=?", "liming")
	if action.Error != nil && action.Error != gorm.ErrRecordNotFound {
		tx.Rollback()
		log.Panic(1, action.Error)
		return
	}
	if action.RowsAffected == 0 {
		tx.Rollback()
		log.Println(2, "no data")
		return
	}

	err := UpdateScore(stu.Id, 99, tx)
	if err != nil {
		panic(err)
	}

	if err := tx.Rollback().Error; err != nil {
		log.Panic("rollback", err)
		return
	}

	return
}

// 这里会阻塞，说明db对象开启了新的会话
func UpdateScore(id, score int, db *gorm.DB) error {
	//var dbConn = new(conn)
	//db.Exec("select connection_id() from student;").Scan(dbConn)
	//log.Printf("connection_id-1:%d", dbConn.id)

	err := db.Debug().Exec("update `student` set `score` = ? where `id` = ?", score, id).Error
	if err != nil {
		log.Panic(err)
	}
	return nil
}

func TestForUpdate(t *testing.T) {
	startDB()
	defer closeDB()
	db.AutoMigrate(Student{})

	err := db.Create(&Student{
		Name:  "liming",
		Score: 100,
	}).Error
	if err != nil {
		panic(err)
	}

	ForUpdateTest(t)
}
