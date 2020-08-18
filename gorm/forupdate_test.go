package main

import (
	"github.com/jinzhu/gorm"
	"log"
	"testing"
)

// 测试锁住id=3的记录时，能不能去修改它, 结果是可以！
func ForUpdate(t *testing.T, db *gorm.DB) {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()
	if tx.Error != nil {
		panic(tx.Error)
	}

	action := tx.Debug().Set("gorm:query_option", "FOR UPDATE").First(&User{}, "id=?", 31)
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

	err := tx.Debug().Delete(User{}, "id=?", 3).Error
	if err != nil {
		tx.Rollback()
		log.Panic(2, err)
		return
	}
	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		log.Panic(3, err)
		return
	}
	return
}

func TestForUpdate(t *testing.T) {
	db, err := gorm.Open("mysql", "test_u:1918ddkk@(114.115.216.44:33061)/test?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	ForUpdate(t, db)
}
