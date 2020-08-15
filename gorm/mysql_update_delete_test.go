package main

import (
	"database/sql"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"testing"
)

func UpdateAllFields(t *testing.T, db *gorm.DB) {
	var user User
	db.First(&user)
	assert.True(t, user.ID > 0)

	user.Name = "x1"
	user.Age = sql.NullInt64{
		Int64: 60,
		Valid: true,
	}
	db.Save(&user)
	// UPDATE users SET u_name='x1', age=100, birthday=old_birthday, updated_at = time.now() WHERE id=user.id;

	var user1 User
	user1.ID = user.ID
	db.Take(&user1)
	assert.True(t, user1.Name == "x1")
}

func UpdateWantedFields(t *testing.T, db *gorm.DB) {
	var user User
	user.ID = 1
	// Update single attribute if it is changed
	db.Model(&user).Update("u_name", "hello") // field, value
	// UPDATE users SET u_name='hello', updated_at=time.Now() WHERE id=111;
	assert.True(t, user.Name == "hello")

	// Update single attribute with combined conditions
	db.Model(&user).Where("age = ?", 18).Update("u_name", "hello")
	// UPDATE users SET u_name='hello', updated_at=time.Now() WHERE id=2 AND age=18;

	// Update multiple attributes with `map`, will only update those changed fields
	db.Model(&user).Updates(map[string]interface{}{"u_name": "hello", "age": 18})
	// UPDATE users SET u_name='hello', age=18 updated_at=time.Now() WHERE id=111;

	// Update multiple attributes with `struct`, will only update those changed & non blank fields
	db.Model(&user).Updates(User{Name: "hello", Age: sql.NullInt64{Int64: 20, Valid: true}})
	// UPDATE users SET u_name='hello', age=18, updated_at = time.Now() WHERE id = 111;

	// WARNING when update with struct, GORM will only update those fields that with non blank value
	// For below Update, nothing will be updated as "", 0, false are blank values of their types
	db.Model(&user).Updates(User{Name: "", Age: sql.NullInt64{Int64: 21, Valid: true}})
	assert.True(t, user.Name == "hello", user.Age.Int64 == 21)
}

// gorm使用主键删除记录，主键为空时会删除所有记录，除非关闭这个规则: db.BlockGlobalUpdate(true)
func Delete(t *testing.T, db *gorm.DB) {
	var user User
	user.ID = 1
	db.BlockGlobalUpdate(true)
	// Delete an existing record
	db.Delete(&user)
	// DELETE from users where id=1;

	// Add extra SQL option for deleting SQL
	db.Set("gorm:delete_option", "OPTION (OPTIMIZE FOR UNKNOWN)").Delete(&user)
	// DELETE from users where id=1 OPTION (OPTIMIZE FOR UNKNOWN);

	// Batch delete
	db.Where("u_name LIKE ?", "%jinzhu%").Delete(User{})
	// DELETE from users where email LIKE "%jinzhu%";

	db.Delete(User{}, "u_name LIKE ?", "%jinzhu%")
	// DELETE from users where email LIKE "%jinzhu%";

	// soft delete, 只要表结构体有 DeleteAt 字段，默认就是软删除
	db.Delete(&user)
	// UPDATE users SET deleted_at="2013-10-29 10:23" WHERE id = 111;

	// Batch Delete
	db.Where("age = ?", 20).Delete(&User{})
	// UPDATE users SET deleted_at="2013-10-29 10:23" WHERE age = 20;

	// Soft deleted records will be ignored when query them
	db.Where("age = 20").Find(&user)
	// SELECT * FROM users WHERE age = 20 AND deleted_at IS NULL;

	// Find soft deleted records with Unscoped
	db.Unscoped().Where("age = 20").Find(&user)
	// SELECT * FROM users WHERE age = 20;

	// Delete record permanently with Unscoped
	db.Unscoped().Delete(&user)
	//// DELETE FROM orders WHERE id=10;
}
