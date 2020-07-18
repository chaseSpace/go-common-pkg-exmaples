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
