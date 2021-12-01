package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
	"testing"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func TestWithSqlite(t *testing.T) {
	db, err := gorm.Open("sqlite3", "sqlite.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer func() {
		if err = db.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	// Migrate the schema
	db.AutoMigrate(&Product{})

	// Create
	db.Create(&Product{Code: "L1212", Price: 1000})

	// Read
	var product Product
	db.First(&product, 1)                   // find product with id 1
	db.First(&product, "code = ?", "L1212") // find product with code l1212

	log.Println(1, product)

	// Update - update product's price to 2000
	db.Model(&product).Update("Price", 2000)
	db.First(&product, 1)
	log.Println(2, product)
	// Delete - delete product
	db.Delete(&product)
}
