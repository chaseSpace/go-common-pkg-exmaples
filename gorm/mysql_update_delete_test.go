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
	// UPDATE admin_users SET u_name='x1', age=100, birthday=old_birthday, updated_at = time.now() WHERE id=user.id;

	var user1 User
	user1.ID = user.ID
	db.Take(&user1)
	assert.True(t, user1.Name == "x1")
}

// 注意：update之后不会将数据写回struct
func UpdateWantedFields(t *testing.T, db *gorm.DB) {
	// 首先最好关掉全局更新, 也就是当where条件为空时，返回err，而不是直接更新/删除全部数据
	// 这是避免人为的失误造成严重后果
	db.BlockGlobalUpdate(true)
	err := db.Delete(&User{}).Error
	assert.Error(t, err) // err: missing WHERE clause while deleting

	var user User
	user.ID = 1
	// Update single attribute if it is changed
	db.Model(&user).Update("u_name", "hello") // field, value
	// UPDATE admin_users SET u_name='hello', updated_at=time.Now() WHERE id=1;
	assert.True(t, user.Name == "hello")

	// Update single attribute with combined conditions
	db.Model(&user).Where("age = ?", 18).Update("u_name", "hello")
	// UPDATE admin_users SET u_name='hello', updated_at=time.Now() WHERE id=2 AND age=18;

	// Update multiple attributes with `map`, will only update those changed fields
	db.Model(&user).Updates(map[string]interface{}{"u_name": "hello", "age": 18})
	// UPDATE admin_users SET u_name='hello', age=18 updated_at=time.Now() WHERE id=1;

	// Update multiple attributes with `struct`, will only update those changed & non blank fields
	db.Model(&user).Updates(User{Name: "hello", Age: sql.NullInt64{Int64: 20, Valid: true}})
	// UPDATE admin_users SET u_name='hello', age=20, updated_at = time.Now() WHERE id = 1;

	// WARNING when update with struct, GORM will only update those fields that with non blank value
	// For below Update, nothing will be updated as "", 0, false are blank values of their types
	db.Model(&user).Updates(User{Name: "", Age: sql.NullInt64{Int64: 21, Valid: true}})
	assert.True(t, user.Name == "hello", user.Age.Int64 == 21)
}

func UpdateSelectFields(t *testing.T, db *gorm.DB) {
	var user User
	user.ID = 1
	db.Model(&user).Select("u_name").Updates(map[string]interface{}{"u_name": "hello1", "age": 99})
	// UPDATE admin_users SET u_name='hello', updated_at='2013-11-17 21:34:10' WHERE id=1;
	assert.True(t, user.Name == "hello1" && user.Age.Int64 != 99)

	var user2 User
	user2.ID = 1
	db.Model(&user2).Omit("u_name").Updates(map[string]interface{}{"u_name": "hello222", "age": 99})
	// UPDATE admin_users SET age=18, actived=false, updated_at='2013-11-17 21:34:10' WHERE id=1;
	assert.True(t, user2.Name != "hello222" && user2.Age.Int64 == 99)
}

// 不更新updateAt字段 (仅限UpdateColumn/UpdateColumns)
func UpdateFieldsOnly(t *testing.T, db *gorm.DB) {
	var user User
	user.ID = 1

	db.Take(&user)
	assert.True(t, user.UpdatedAt.Year() == 2020)

	var user1 User
	user1.ID = 1
	// Update single attribute, similar with `Update`
	db.Model(&user1).UpdateColumn("u_name", "hello")
	// UPDATE admin_users SET u_name='hello' WHERE id = 1;
	db.Take(&user1)
	assert.True(t, user.UpdatedAt == user1.UpdatedAt)

	var user2 User
	user2.ID = 1
	// Update multiple attributes, similar with `Updates`
	db.Model(&user2).UpdateColumns(User{Name: "hello", Age: sql.NullInt64{Int64: 25, Valid: true}})
	// UPDATE admin_users SET u_name='hello', age=25 WHERE id = 1;
	db.Take(&user2)
	assert.True(t, user.UpdatedAt == user2.UpdatedAt && user2.Age.Int64 == 25)
}

func BatchUpdate(t *testing.T, db *gorm.DB) {
	n := db.Table("admin_users").Where("id IN (?)", []int{1, 111}).Updates(map[string]interface{}{"u_name": "BatchUpdate", "age": 18}).RowsAffected
	// UPDATE admin_users SET name='hello', age=18 WHERE id IN (10, 11);
	assert.True(t, n == 1)

	// Update with struct only works with none zero values, or use map[string]interface{}
	db.Model(User{Name: "hello"}).Updates(User{Name: "BatchUpdate", Age: sql.NullInt64{Int64: 25, Valid: true}})
	// UPDATE admin_users SET name='hello', age=18 where name='hello' and deleted_at is null;

	// Get updated records count with `RowsAffected`
	n = db.Model(User{Name: "hello"}).Updates(User{Name: "BatchUpdate123", Age: sql.NullInt64{Int64: 25, Valid: true}}).RowsAffected
	assert.True(t, n == 0)
}

func UpdateWithSqlExpr() {
	/*
		DB.Model(&product).Update("price", gorm.Expr("price * ? + ?", 2, 100))
		// UPDATE "products" SET "price" = price * '2' + '100', "updated_at" = '2013-11-17 21:34:10' WHERE "id" = '2';

		DB.Model(&product).Updates(map[string]interface{}{"price": gorm.Expr("price * ? + ?", 2, 100)})
		// UPDATE "products" SET "price" = price * '2' + '100', "updated_at" = '2013-11-17 21:34:10' WHERE "id" = '2';

		DB.Model(&product).UpdateColumn("quantity", gorm.Expr("quantity - ?", 1))
		// UPDATE "products" SET "quantity" = quantity - 1 WHERE "id" = '2';

		DB.Model(&product).Where("quantity > 1").UpdateColumn("quantity", gorm.Expr("quantity - ?", 1))
		// UPDATE "products" SET "quantity" = quantity - 1 WHERE "id" = '2' AND quantity > 1;
	*/
}

func ChangeValuesInHooks() {
	// 如果想在钩子函数中修改待更新的字段值，可以给表结构体定义方法BeforeUpdate, BeforeSave
	// 使用scope.SetColumn, 例如:

	/*
		func(user *User) BeforeSave(scope * gorm.Scope)(err error) {
			if pw, err := bcrypt.GenerateFromPassword(user.Password, 0); err == nil {
				scope.SetColumn("EncryptedPassword", pw)
			}
		}
	*/
}

func AddExtraUpdateOption() {
	// Add extra SQL option for updating SQL
	// db.Model(&user).Set("gorm:update_option", "OPTION (OPTIMIZE FOR UNKNOWN)").Update("name", "hello")
	// UPDATE users SET name='hello', updated_at = '2013-11-17 21:34:10' WHERE id=111 OPTION (OPTIMIZE FOR UNKNOWN);

}

// 默认所有的删除操作都是软删除（前提是有deleted_at字段）
func DeleteOneRecord(t *testing.T, db *gorm.DB) {
	var user User
	user.ID = 1
	// &struct删除要求结构体[主键]必须有值，否则数据全部删除！！！(可以通过db.BlockGlobalUpdate(true)避免)
	// gorm不会使用主键以外的字段来删除
	db.Delete(&user)

	err := db.Take(&user).Error
	assert.True(t, err.Error() == "record not found")

	// db.BlockGlobalUpdate(true) 时，不能在没有where语句时删除数据, 主键以外的字段不会被gorm用作where语句
	db.BlockGlobalUpdate(true)
	err = db.Debug().Delete(User{Name: "hello"}).Error
	assert.Error(t, err) // err: missing WHERE clause while deletin

	// add extra option
	// db.Set("gorm:delete_option", "OPTION (OPTIMIZE FOR UNKNOWN)").Delete(&email)
}

func BatchDelete(t *testing.T, db *gorm.DB) {
	n := db.Where("id = ?", 2).Delete(User{}).RowsAffected
	// DELETE from admin_users where id = 2 and deleted_at is null;
	assert.True(t, n == 1)

	var user User
	user.ID = 2
	err := db.Take(&user).Error
	assert.True(t, err.Error() == "record not found")

	n = db.Delete(User{}, "u_name LIKE ?", "NOT_EXIST%").RowsAffected
	// DELETE from admin_users where u_name LIKE 'NOT_EXIST%' and deleted_at is null;
	assert.True(t, n == 0)
}

// 永久删除数据
func DeletePermanently(t *testing.T, db *gorm.DB) {
	var user User
	user.ID = 1
	// Delete record permanently with Unscoped
	db.Unscoped().Delete(&user)
	// DELETE FROM orders WHERE id=10;

	err := db.Raw("SELECT * FROM admin_users WHERE id=?", 1).Scan(&user).Error
	assert.Error(t, err) // err: record not found
}
