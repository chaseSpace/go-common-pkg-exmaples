package main

import (
	"fmt"
	"github.com/jmoiron/sqlx/reflectx"
	"log"
	"strings"
	"time"

	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/jmoiron/sqlx"
)

var ckDB *sqlx.DB

func initCK() {
	var err error
	ckDB, err = sqlx.Open("clickhouse", "tcp://192.168.1.168:9000?username=default&password=&database=xman&debug=true")
	if err != nil {
		log.Fatal(err)
	}
	ckDB.Mapper = reflectx.NewMapperFunc("json", strings.ToLower)
}

func main() {
	initCK()
	defer ckDB.Close()
	Select()
	Exec()
}

func Select() {
	var items []struct {
		EventDate time.Time `json:"event_date"`
		GiftTime  time.Time `json:"gift_time"`
		Sender    uint32    `json:"sender"`
		Receiver  uint32    `json:"receiver"`
	}

	if err := ckDB.Select(&items, "SELECT event_date,gift_time,sender,receiver FROM gift_data limit 5"); err != nil {
		log.Fatal(err)
	}
	for _, item := range items {
		log.Print(item)
	}
}

func Exec() {
	// ck 的insert  update delete必须用这种事务模式，麻烦
	tx, _ := ckDB.Begin()
	result, err := tx.Exec(`insert into gift_data (event_date,gift_time,sender,receiver) values (?,?,?,?)`, "2020-10-11", "2020-10-11 00:00:00", 100, 100)
	if err != nil {
		log.Fatal("err: ", err)
	}
	err = tx.Commit()
	fmt.Println(111, err)
	fmt.Println(result.RowsAffected()) // not supported
	fmt.Println(result.LastInsertId()) // not supported
	// 不如直接使用gorm 反正都不支持RowsAffected
}
