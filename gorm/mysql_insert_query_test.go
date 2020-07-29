package main

import (
	"database/sql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

/*
gorm的坑：
	1. db.Model(&table_struct).Find(&other_struct) 会查到已被删除的记录，还是用回Find(&table_struct)

Mysql的注意点：
	1. rows affected这个属性在update时，如果新旧数据一致，它也是0，并不代表记录不存在
*/

/* 数据库需要对应driver
import _ "github.com/jinzhu/gorm/dialects/mysql"
// import _ "github.com/jinzhu/gorm/dialects/postgres"
// import _ "github.com/jinzhu/gorm/dialects/sqlite"
// import _ "github.com/jinzhu/gorm/dialects/mssql"

*/

// table定义
// 表名会默认被创建成复数，即users，可以禁用此特点
type User struct {
	/* gorm.Model建表后的结果， uint32 == uint ==> int(10) unsigned
	`id` int(10) unsigned NOT NULL AUTO_INCREMENT,
	 `created_at` datetime DEFAULT NULL,
	 `updated_at` datetime DEFAULT NULL,
	 `deleted_at` datetime DEFAULT NULL,
	*/
	gorm.Model // 进去看它的代码就知道其作用：可选，主要是嵌入一些基本字段，如 id，updatedAt,createdAt,deletedAt
	// 写了它就不需要再定义id这些基本字段，注意DeletedAt字段是指针，因为在数据未被删除时这个字段应该是nil

	//ID string `gorm:"primary_key"`//  primary_key标签也是可选的，gorm默认把id当主键
	Name         string        `gorm:"column:u_name"` // tag修改字段名（默认字段命名规则是小写+下划线）
	Age          sql.NullInt64 `gorm:"default:'18'"`  // 默认值会在字段为nil或类型零值时被使用
	Birthday     *time.Time
	Email        string `gorm:"type:varchar(100);unique_index"` // 另一种用法是type:text
	Role         string `gorm:"size:255"`                       // size:255等效于varchar(255)
	MemberNumber string `gorm:"unique;not null"`                // create unique_index and not null, so unique_index == unique
	Num          int    `gorm:"AUTO_INCREMENT"`                 // set num to auto incrementable
	Address      string `gorm:"index:addr"`                     // create index with name `addr` for address
	IgnoreMe     int    `gorm:"-"`                              // ignore this field
}

// table
type Order struct {
	gorm.Model
	State  string
	Amount int64
	UserId int64
}

/*
CreatedAt：有此字段的表在插入时此字段无需传入，gorm会设置为当前时间，UpdatedAt和DeletedAt同理
有DeletedAt字段的表数据在删除时不会真的删除，而是给这个字段设置删除时间（除非你用手写SQL去删除）
*/
// 修改表名
func (User) TableName() string {
	return "admin_users"
}

/* 禁用复数
db,err := gorm.Open(...)
db.SingularTable(true)
*/

/*
一些快捷的建表、查询方法（建议项目中不要用，不规范）
// Create `deleted_users` table with struct User's definition
db.Table("deleted_users").CreateTable(&User{})

var deleted_users []User
db.Table("deleted_users").Find(&deleted_users)
// SELECT * FROM deleted_users;

db.Table("deleted_users").Where("name = ?", "jinzhu").Delete()
// DELETE FROM deleted_users WHERE name = 'jinzhu';
*/

func TestMysql(t *testing.T) {
	// 自定义table命名规则, 不要和TableName()同时存在
	//gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
	//	return "prefix_" + defaultTableName
	//}

	// uri方式连接
	// user:password@(localhost)/dbname?charset=utf8&parseTime=True&loc=Local
	db, err := gorm.Open("mysql", "test_u:1918ddkkdd@(114.115.216.44:33061)/test?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db = db.DropTableIfExists(&User{}, &Order{})
	// processes err
	if db.Error != nil {
		panic(err)
	}

	// create table
	db.CreateTable(&User{}, &Order{})

	//if !db.HasTable(&User{}) {
	//}

	InsertTest(t, db)
	CommonQueryTest(t, db)
	//QueryNotTest(t, db)
	//QueryOrTest(t, db)
	//MoreSimpleQueryTest(t, db)
	//FirstOrInitQueryTest(t, db)
	//FirstOrCreateQueryTest(t, db)
	//SubQueryTest(t, db)
	//SelectTest(t, db)
	//LimitTest(t, db)
	//OffsetTest(t, db)
	//CountTest(t, db)
	//JoinTest(t, db)
	//ScanTest(t, db)

	UpdateAllFields(t, db)
	UpdateWantedFields(t, db)
}

func InsertTest(t *testing.T, db *gorm.DB) {
	var record = User{
		Name:         "x",
		Num:          0,
		Email:        "e1",
		MemberNumber: "0",
	}

	ok := db.NewRecord(record) // true, 检测记录的主键是否零值，不插数据，也不会与db交互.
	assert.True(t, ok)

	db.Create(&record) // insert, db生成的id将会写入record

	assert.NotEqual(t, record.ID, 0)
	ok = db.NewRecord(record) // false, because record.id already exists in db.
	assert.False(t, ok)

	record.ID = 0
	record.Age = sql.NullInt64{Valid: true, Int64: 0}
	record.Email = "e2"
	record.MemberNumber = "1"
	db.Create(&record)

	record.ID = 0
	record.Age = sql.NullInt64{Valid: true, Int64: 17}
	record.Email = "e3"
	record.MemberNumber = "2"
	db.Create(&record)
}

func CommonQueryTest(t *testing.T, db *gorm.DB) {
	var user User
	// 不要将条件放在结构体内，不会读取的，只有主键会被作为条件, 后面的Take方法也是
	db.Debug().First(&user, "u_name=?", "x") // SELECT * FROM users WHERE u_name=x ORDER BY id LIMIT 1;
	assert.NotEqual(t, user.ID, 0)

	u := new(User)
	// Get one record, no specified order (只使用主键查询，其他字段不会使用)
	// SELECT * FROM `admin_users`  WHERE `admin_users`.`deleted_at` IS NULL AND `admin_users`.`id` = 1 LIMIT 1
	db.Debug().Take(u, "u_name=?", "x")
	log.Printf("111 %+v", u)

	// Get last record, order by primary key
	db.Last(&user)

	// 获取不存在的记录
	takeErr := db.Take(&User{}, "u_name=?", "NOT_EXIST").Error
	FindErr := db.Find(&[]User{}, "u_name=?", "NOT_EXIST").Error
	// !!! 注意这个err，当接收者是一个结构体时且数据未找到时返回
	assert.Equal(t, takeErr, gorm.ErrRecordNotFound)
	// slice接收，则是nil
	assert.Nil(t, FindErr)

	var users []User
	// Get all records
	db.Find(&users)
	assert.True(t, len(users) > 1)

	// Get record with primary key (后面参数只会传递给整型主键)
	db.First(&user, 10)

	// where可以自定义字符串形式的条件，使用问号占位参数
	// 还可以传入带值的struct，其中的值作为条件查询
	// 还可以传入map类型，其中的k-v作为条件查询
	// 还可以传入slice类型，不过元素只能是整型，作为主键字段参数，执行IN查询（如果主键不是整型，应该会报错）

	var user1 User
	// where
	db.Where("u_name = ?", "x").First(&user1)
	// SELECT * FROM users WHERE name = 'x' ORDER BY id LIMIT 1;
	assert.True(t, user1.Name == "x")

	var users1 []User
	db.Where("u_name = ?", "x").Find(&users1)
	// SELECT * FROM users WHERE name = 'x';
	assert.True(t, len(users) > 1)

	// IN
	db.Where("u_name IN (?)", []string{"x", "jinzhu 2"}).Find(&users)
	// SELECT * FROM users WHERE name in ('x','jinzhu 2');

	// LIKE
	db.Where("u_name LIKE ?", "%x%").Find(&users)
	// SELECT * FROM users WHERE name LIKE '%jin%';

	// AND
	db.Where("u_name = ? AND age >= ?", "x", "18").Find(&users)
	// SELECT * FROM users WHERE name = 'x' AND age >= 22;

	var users2 []User
	// Time
	db.Where("updated_at > ?", time.Now().Add(-time.Hour)).Find(&users2)
	// SELECT * FROM users WHERE updated_at > 'an hour ago';
	assert.True(t, len(users2) > 1)

	var users3 []User
	// BETWEEN
	db.Where("created_at BETWEEN ? AND ?", time.Now().Add(-time.Hour), time.Now()).Find(&users3)
	// SELECT * FROM users WHERE created_at BETWEEN 'an hour ago' AND 'now time';
	assert.True(t, len(users3) > 1)

	var user2 User
	// struct作为条件查询
	// 注意：struct作为条件时，其中字段的零值将会被gorm忽略，比如0,'',false
	// 如果不想被忽略，只能在定义结构体时，将字段类型定义为指针或scanner/valuer, 如sql.NullInt64/NullString...
	db.Where(&User{Name: "x", Age: sql.NullInt64{Int64: 18, Valid: true}}).First(&user2)
	// SELECT * FROM users WHERE name = "x" AND age = 18 ORDER BY id LIMIT 1;
	assert.True(t, user2.Age.Int64 == 18)

	var user3 User
	db.Where(&User{Name: "x", Age: sql.NullInt64{Int64: 0, Valid: true}}).First(&user3)
	// SELECT * FROM users WHERE name = "x" AND age = 0 ORDER BY id LIMIT 1;
	assert.True(t, user3.ID > 0 && user3.Age.Int64 == 0)

	// Map作为条件查询
	db.Where(map[string]interface{}{"u_name": "x", "age": 18}).Find(&users)
	// SELECT * FROM users WHERE name = "x" AND age = 18;

	// Slice作为条件查询
	db.Where([]int64{1, 21, 22}).Find(&users)
	// SELECT * FROM users WHERE id IN (1, 21, 22);
}

func QueryNotTest(t *testing.T, db *gorm.DB) {
	var user User
	db.Not("u_name", "jinzhu").First(&user)
	// SELECT * FROM users WHERE name <> "jinzhu" ORDER BY id LIMIT 1;
	assert.True(t, user.ID > 0)

	var users []User
	// Not In
	db.Not("u_name", []string{"jinzhu", "jinzhu 2"}).Find(&users)
	// SELECT * FROM users WHERE name NOT IN ("jinzhu", "jinzhu 2");
	assert.True(t, len(users) > 1)

	// Not In slice of primary keys
	db.Not([]int64{1, 2, 3}).First(&user)
	// SELECT * FROM users WHERE id NOT IN (1,2,3) ORDER BY id LIMIT 1;

	// Special case
	db.Not([]int64{}).First(&user)
	// SELECT * FROM users ORDER BY id LIMIT 1;

	// Plain SQL
	db.Not("u_name = ?", "jinzhu").First(&user)
	// SELECT * FROM users WHERE NOT(name = "jinzhu") ORDER BY id LIMIT 1;

	// Struct
	db.Not(User{Name: "jinzhu"}).First(&user)
	// SELECT * FROM users WHERE name <> "jinzhu" ORDER BY id LIMIT 1;
}

func QueryOrTest(t *testing.T, db *gorm.DB) {
	var users []User
	db.Where("role = ?", "").Or("role = ?", "super_admin").Find(&users)
	// SELECT * FROM users WHERE role = 'admin' OR role = 'super_admin';
	assert.True(t, len(users) > 1)

	var users1 []User
	// Struct
	db.Where("u_name = 'x'").Or(User{Name: "jinzhu 2"}).Find(&users1)
	// SELECT * FROM users WHERE u_name = 'x' OR name = 'jinzhu 2';
	assert.True(t, len(users) > 1)

	// Map
	db.Where("u_name = 'jinzhu'").Or(map[string]interface{}{"u_name": "jinzhu 2"}).Find(&users)
	// SELECT * FROM users WHERE u_name = 'jinzhu' OR u_name = 'jinzhu 2';
	assert.True(t, len(users) == 0)
}

// gorm称之为inline condition，内联查询，我看来就是更简单的一种查询写法
func MoreSimpleQueryTest(t *testing.T, db *gorm.DB) {
	var users []User
	var user User

	// Get by primary key (only works for integer primary key)
	db.First(&user, 2)
	// SELECT * FROM users WHERE id = 2;
	assert.True(t, user.ID == 2)

	// Get by primary key if it were a non-integer type
	db.First(&user, "id = ?", "string_primary_key")
	// SELECT * FROM users WHERE id = 'string_primary_key';

	// Plain SQL
	db.Find(&user, "u_name = ?", "jinzhu")
	// SELECT * FROM users WHERE u_name = "jinzhu";

	db.Find(&users, "u_name <> ? AND age = ?", "jinzhu", 18)
	// SELECT * FROM users WHERE u_name <> "jinzhu" AND age = 18;
	assert.True(t, len(users) > 0)

	var users1 []User
	// Struct
	db.Find(&users1, User{Age: sql.NullInt64{Int64: 18, Valid: true}})
	// SELECT * FROM users WHERE age = 18;
	assert.True(t, len(users1) > 0)

	var users2 []User
	// Map
	db.Find(&users2, map[string]interface{}{"age": 18})
	// SELECT * FROM users WHERE age = 18;
	assert.True(t, len(users2) > 0)
}

func FirstOrInitQueryTest(t *testing.T, db *gorm.DB) {
	var user User
	// 先介绍for update
	db.Set("gorm:query_option", "FOR UPDATE").First(&user, 10)
	// SELECT * FROM users WHERE id = 10 FOR UPDATE;

	// FirstOrInit, 获取匹配的第一条，如果没有就用给定的条件初始化传入的user（没有往db插入），仅支持struct和map
	// 针对不存在的数据
	db.FirstOrInit(&user, User{Name: "non_existing"})
	// user -> User{ ID: N!=0, Name: "non_existing"}
	assert.True(t, user.ID == 0 && user.Name == "non_existing")

	// 针对存在的数据，另外的2种写法
	db.Where(User{Name: "x"}).FirstOrInit(&user)
	db.FirstOrInit(&user, map[string]interface{}{"u_name": "x"})

	assert.True(t, user.Age.Int64 == 18)

	var user1 User
	// 使用Attrs方法，是FirstOrInit的扩展，它将仅作为初始化的参数与查询参数隔离开

	// 针对不存在的数据
	db.Where(User{Name: "non_existing"}).Attrs(User{Age: sql.NullInt64{Int64: 18, Valid: true}}).FirstOrInit(&user1)
	// SELECT * FROM USERS WHERE name = 'non_existing' ORDER BY id LIMIT 1;
	// user -> User{Name: "non_existing", Age: 18}
	assert.True(t, user1.ID == 0 && user1.Name == "non_existing" && user1.Age.Int64 == 18)

	// 更简便的写法
	db.Where(User{Name: "non_existing"}).Attrs("age", 18).FirstOrInit(&user)
	// SELECT * FROM USERS WHERE name = 'non_existing' ORDER BY id LIMIT 1;
	// user -> User{Name: "non_existing", Age: 18}

	// 针对存在的数据
	db.Where(User{Name: "x`"}).Attrs(User{Age: sql.NullInt64{Int64: 18, Valid: true}}).FirstOrInit(&user)
	// SELECT * FROM USERS WHERE u_name = x' ORDER BY id LIMIT 1;
	// user -> User{Id: n, Name: "x", Age: 18}

	var user2 User
	// 使用Assign方法，是FirstOrInit的扩展，它将仅作为初始化的参数与查询参数隔离开
	// 与Attrs的不同是，<无论是否匹配到结果>，都会将数据设置到传入的结构体
	// 针对不存在的数据
	db.Where(User{Name: "non_existing"}).Assign(User{Age: sql.NullInt64{Int64: 20, Valid: true}}).FirstOrInit(&user2)
	// user -> User{Name: "non_existing", Age: 20}
	assert.True(t, user2.ID == 0 && user2.Name == "non_existing" && user2.Age.Int64 == 20)

	// 针对存在的数据
	db.Where(User{Name: "x"}).Assign(User{Age: sql.NullInt64{Int64: 22, Valid: true}}).FirstOrInit(&user2)
	// SELECT * FROM USERS WHERE name = x' ORDER BY id LIMIT 1;
	// user -> User{Id: n, Name: "x", Age: 22}
	assert.True(t, user2.ID > 0 && user2.Name == "x" && user2.Age.Int64 == 22)
}

func FirstOrCreateQueryTest(t *testing.T, db *gorm.DB) {
	// FirstOrCreate, 获取匹配的第一条，如果没有就用给定的条件往db插入一条数据，仅支持struct和map
	var user User
	// not found
	db.FirstOrCreate(&user, User{Name: "non_existing"})
	// INSERT INTO "users" (name) VALUES ("non_existing");
	// user -> User{Id: N, Name: "non_existing"}
	assert.True(t, user.ID > 0 && user.Name == "non_existing")

	user.ID = 0

	// Found
	db.Where(User{Name: "x"}).FirstOrCreate(&user)
	// user -> User{Id: N, Name: "x"}

	var user1 User
	// Attrs()，是FirstOrCreate的扩展，它将仅作为初始化和插入db的参数与查询参数隔离开
	// not found
	db.Where(User{Name: "non_existing666", Email: "e123", MemberNumber: "m123"}).Attrs(User{Age: sql.NullInt64{Int64: 20, Valid: true}}).FirstOrCreate(&user1)
	// SELECT * FROM users WHERE name = 'non_existing' ORDER BY id LIMIT 1;
	// INSERT INTO "users" (name, age) VALUES ("non_existing666", 20);
	// user -> User{Id: N, Name: "non_existing666", Age: 20}
	assert.True(t, user1.ID > 0 && user1.Name == "non_existing666", user1.Age.Int64 == 20)

	var user2 User
	// Found
	db.Where(User{Name: "x"}).Attrs(User{Age: sql.NullInt64{Int64: 20, Valid: true}}).FirstOrCreate(&user2)
	// SELECT * FROM users WHERE name = 'x' ORDER BY id LIMIT 1;
	// user -> User{Id: N, Name: "x", Age: 18}
	assert.True(t, user2.Age.Int64 == 18)

	var user3 User
	// Assign(), 是FirstOrCreate的扩展，它将仅作为初始化和插入db的参数与查询参数隔离开
	// 与Attrs的不同是，<无论是否匹配到结果>，都会将数据插入/更新到db

	// not found
	db.Where(User{Name: "non_existing667", Email: "e124", MemberNumber: "m124"}).Assign(User{Age: sql.NullInt64{Int64: 20, Valid: true}}).FirstOrCreate(&user3)
	// SELECT * FROM users WHERE name = 'non_existing667' ORDER BY id LIMIT 1;
	// INSERT INTO "users" (name, age) VALUES ("non_existing667", 20);
	// user -> User{Id: 112, Name: "non_existing667", Age: 20}

	var user4 User
	// Found
	db.Where(User{Name: "x"}).Assign(User{Age: sql.NullInt64{Int64: 25, Valid: true}}).FirstOrCreate(&user4)
	// SELECT * FROM users WHERE name = 'x' ORDER BY id LIMIT 1;
	// UPDATE users SET age=30 WHERE id = 111;
	// user -> User{Id: 111, Name: "x", Age: 25}
	assert.True(t, user4.ID > 0 && user4.Age.Int64 == 25)
}

func SubQueryTest(t *testing.T, db *gorm.DB) {
	// 遇到的问题：下面注释的写法，第二条记录无法插入
	//db.FirstOrCreate(&Order{State: "paid", Amount: 10})
	//db.FirstOrCreate(&Order{State: "paid", Amount: 20})

	// 先给orders表插几条测试数据,忽略错误处理
	db.Where(&Order{State: "paid", Amount: 10, UserId: 1}).FirstOrCreate(&Order{})
	db.Where(&Order{State: "paid", Amount: 20, UserId: 1}).FirstOrCreate(&Order{})

	var orders []Order
	db.Where("amount>?", db.Table("orders").Select("AVG(amount)").Where("state=?", "paid").SubQuery()).Find(&orders)
	// SELECT * FROM "orders"  WHERE "orders"."deleted_at" IS NULL AND (amount > (SELECT AVG(amount) FROM "orders"  WHERE (state = 'paid')));
	assert.True(t, len(orders) == 1 && orders[0].Amount == 20)
}

func SelectTest(t *testing.T, db *gorm.DB) {
	var users []User
	db.Select("u_name, age").Find(&users)
	// SELECT name, age FROM users;

	db.Select([]string{"u_name", "age"}).Find(&users)
	// SELECT name, age FROM users;

	assert.True(t, len(users) > 1)

	loc, _ := time.LoadLocation("Asia/Shanghai")
	ti := time.Date(2020, 10, 10, 0, 0, 0, 0, loc)
	rows, err := db.Table("admin_users").Select("COALESCE(birthday,?)", ti).Rows()
	// SELECT COALESCE(age,'42') FROM admin_users;

	assert.Nil(t, err)
	var count int
	for rows.Next() {
		var birth []uint8 // 上面参数是time，但scan时却是[]uint8，其实是因为上面传入时变成了str
		err = rows.Scan(&birth)
		assert.Nil(t, err)
		assert.True(t, string(birth) == ti.String()[:19]) // 2020-10-10 00:00:00
		count++
	}

	assert.True(t, count > 1)
}

func OrderTest(t *testing.T, db *gorm.DB) {
	var users []User
	db.Order("age desc, name").Find(&users)
	// SELECT * FROM users ORDER BY age desc, name;

	// Multiple orders
	db.Order("age desc").Order("name").Find(&users)
	// SELECT * FROM users ORDER BY age desc, name;

	// ReOrder
	db.Order("age desc").Find(&users).Order("age", true).Find(&users)
	// SELECT * FROM users ORDER BY age desc; (users1)
	// SELECT * FROM users ORDER BY age; (users2)
}

func LimitTest(t *testing.T, db *gorm.DB) {
	var users []User
	db.Limit(2).Find(&users)
	// SELECT * FROM users LIMIT 2;
	assert.True(t, len(users) == 2)

	var users1 []User
	// Cancel limit condition with -1
	db.Limit(10).Find(&users).Limit(-1).Find(&users1)
	// SELECT * FROM users LIMIT 10; (users1)
	// SELECT * FROM users; (users2)
	assert.True(t, len(users1) > 2)
}

// offset需要配合limit使用
func OffsetTest(t *testing.T, db *gorm.DB) {
	var users []User
	db.Debug().Offset(1).Find(&users)
	// SELECT * FROM users OFFSET 3; !!! 这条命令在mysql上是无效的，offset前必须加limit
	// 通过debug发现gorm最终执行的是：SELECT * FROM `admin_users`  WHERE `admin_users`.`deleted_at` IS NULL
	// offset并没有生效, 下同

	var users1 []User
	// Cancel offset condition with -1
	db.Debug().Offset(10).Find(&users).Offset(-1).Find(&users1)
	// SELECT * FROM users OFFSET 10; (users1)
	// SELECT * FROM users; (users2)

	var users2 []User
	// 下面使用limit加offset, 翻译后的sql是正常的，即 ... limit 1 offset 1;
	db.Debug().Limit(1).Offset(1).Find(&users2)
	assert.True(t, len(users2) == 1)

	// 分页查询
	var users3 []User
	var page = 2
	var pageSize = 1
	db.Limit(pageSize).Offset((page - 1) * pageSize).Find(&users3)
	assert.True(t, len(users3) == 1)
}

// gorm提示：Count()必须是链式调用的最后一个方法
func CountTest(t *testing.T, db *gorm.DB) {
	var users []User
	var count int
	db.Where("u_name = ?", "x").Or("u_name = ?", "jinzhu 2").Find(&users).Count(&count)
	assert.True(t, count > 0)

	var count1 int
	db.Model(&User{}).Where("u_name = ?", "x").Count(&count1)
	// SELECT count(*) FROM users WHERE u_name = 'x'; (count)
	assert.True(t, count == count1)

	//db.Table("deleted_users").Count(&count)
	// SELECT count(*) FROM deleted_users;

	//db.Table("deleted_users").Select("count(distinct(name))").Count(&count)
	// SELECT count( distinct(name) ) FROM deleted_users; (count)
}

func JoinTest(t *testing.T, db *gorm.DB) {
	// db.Model()可以替换为db.Table("tableName")
	rows, err := db.Model(&User{}).Select("admin_users.u_name, orders.state").Joins("left join orders on orders.user_id = admin_users.id").Rows()
	assert.Nil(t, err)

	var count int
	type Result struct {
		UName string
		State *string // 如果要scan的数据可能包含null，就得声明为指针类型，当读取null时设置为nil
		// 否则不知道字段在db中是null还是空字符串
		// 不过这种方式scan遇到null不会报错
	}
	for rows.Next() {
		// 两种scan方式
		//  	1. 直接调用rows.Scan(dest1,dest2...)，要求传入与接收字段数量相同的变量
		// 		2. 调用db.ScanRows(&struct) // 就是上面的Result
		var r Result
		err = db.ScanRows(rows, &r)
		//log.Printf("uname:%s, state:%v\n", r.UName, r.State)

		//var (
		//	name string
		//	state *string // 如果要scan的数据可能包含null，就得声明为指针类型，否则这种方式scan遇到null要出错
		//)
		//err = rows.Scan(&name, &state)
		//log.Printf("uname:%s, state:%v\n", name, state)

		assert.Nil(t, err)
		count++
	}
	assert.True(t, count > 0)

	// 取部分字段
	//db.Table("users").Select("users.name, emails.email").Joins("left join emails on emails.user_id = users.id").Scan(&results)
	//
	//// multiple joins with parameter
	//db.Joins("JOIN emails ON emails.user_id = users.id AND emails.email = ?", "jinzhu@example.org").Joins("JOIN credit_cards ON credit_cards.user_id = users.id").Where("credit_cards.number = ?", "411111111111").Find(&user)
}

// gorm的Scan
func ScanTest(t *testing.T, db *gorm.DB) {

	// 定义的切片元素类型必须是struct，不能是[]int这种，无法被scan
	type Result struct {
		Name string
		Age  int
	}

	var result []Result
	db.Table("admin_users").Select("u_name, age").Where("u_name = ?", "x").Scan(&result)
	assert.True(t, len(result) > 1)

	// slice元素不是struct，scan出来全是默认值，是无效的
	var result1 []uint
	db.Model(&User{}).Select("id").Scan(&result1)
	assert.True(t, len(result1) > 0)
	for _, id := range result1 {
		assert.True(t, id == 0)
	}

	// slice元素不是struct，scan出来全是默认值，是无效的
	var result11 []string
	db.Model(&User{}).Select("u_name").Scan(&result11)
	assert.True(t, len(result11) > 0)
	for _, name := range result11 {
		assert.True(t, name == "")
	}

	var result2 []Result
	// Raw SQL
	db.Raw("SELECT u_name, age FROM admin_users WHERE u_name = ?", "x").Scan(&result2)
	assert.Equal(t, result, result2)

}
