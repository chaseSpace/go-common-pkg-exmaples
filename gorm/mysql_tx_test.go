package main

import (
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"testing"
)

// 事务

func CommitTxTest(t *testing.T, db *gorm.DB) {
	err := db.Transaction(func(tx *gorm.DB) error {
		var user = User{Name: "Li", MemberNumber: "Li123", Email: "Li@163.com"}

		// do some database operations in the transaction (use 'tx' from this point, not 'db')
		if err := tx.Create(&user).Error; err != nil {
			// return any error will rollback
			return err
		}

		assert.True(t, user.ID > 0)

		if err := tx.Create(&Order{UserId: int64(user.ID), State: "unpaid"}).Error; err != nil {
			return err
		}

		// return nil will commit
		return nil
	})
	assert.Nil(t, err)
}

func RollbackTxTest(t *testing.T, db *gorm.DB) {
	var user = User{Name: "Li-1", MemberNumber: "Li124", Email: "Li-1@163.com"}
	_ = db.Transaction(func(tx *gorm.DB) error {
		// do some database operations in the transaction (use 'tx' from this point, not 'db')
		if err := tx.Create(&user).Error; err != nil {
			// return any error will rollback
			return err
		}

		assert.True(t, user.ID > 0)

		order := &Order{UserId: int64(user.ID), State: "unpaid"}
		order.ID = 1
		tx.Create(order)

		// ID重复会导致出错
		err := tx.Create(order).Error
		assert.Error(t, err)
		if err != nil {
			return err
		}

		// return nil will commit
		return nil
	})

	// 回滚后这条记录应该找不到
	n := db.Debug().Take(&user).RowsAffected
	assert.True(t, n == 0)
}

func ManualExecTx(t *testing.T, db *gorm.DB) {
	// not test
	_ = func() error {
		// begin a transaction
		tx := db.Begin()
		defer func() {
			// 处理事务内panic
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()

		if tx.Error != nil {
			return tx.Error
		}

		// do some database operations in the transaction (use 'tx' from this point, not 'db')
		err := tx.Create(&User{}).Error
		if err != nil {
			tx.Rollback()
			return err
		}

		// Or commit the transaction
		return tx.Commit().Error
	}
}
