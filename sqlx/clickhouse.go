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

/*
ck官方推荐使用sql或sqlx库来操作db，但经过测试，发现sqlx对ck的支持也是一般般，看下面的描述；
所以如果你的项目里面已经用了某个orm库，比如gorm，可直接使用，不需要添加sql或sqlx
*/
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
	// 这里数据必须使用?占位，不能直接写在sql里，有点恶心。。
	result, err := tx.Exec(`insert into gift_data (event_date,gift_time,sender,receiver) values (?,?,?,?)`, "2020-10-11", "2020-10-11 00:00:00", 100, 100)
	if err != nil {
		log.Fatal("err: ", err)
	}
	err = tx.Commit()
	fmt.Println(111, err)
	fmt.Println(result.RowsAffected()) // not supported
	fmt.Println(result.LastInsertId()) // not supported
	// 不如直接使用gorm 反正都不支持RowsAffected
	// 如果是gorm v1.20版本，连接ck的方法如下
	//gormConf := gorm.Config{
	//	NamingStrategy: schema.NamingStrategy{
	//		SingularTable: true,
	//	},
	//}
	//// gorm1.20版本也没有ck的dialector，用postgres代替, gorm提供的mysql驱动不行
	// CkDB, err = gorm.Open(postgres.New(postgres.Config{
	//	DriverName: "clickhouse",
	//	DSN:        "tcp://192.168.1.168:9000?username=default&password=xxx&database=DDDBBBread_timeout=10&write_timeout=20",
	//},
	//), &gormConf)
}
