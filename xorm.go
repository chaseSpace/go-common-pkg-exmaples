package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"log"
	"os"
	"testing"
	"time"
	"xorm.io/core"
)

var engine *xorm.Engine

/*测试*/

type Attr struct {
	IsTeacher  bool
	Gender     string
	LikeFruits []string
}

// 定义User表
type UserBase struct {
	Id   int64  `xorm:"pk autoincr"` // 根据映射规则 转换为mysql字段 id, 并自动视为主键自增(仅命名为Id时)
	Name string `xorm:"varchar(25) notnull unique 'usr_name'"`
	Desc string `xorm:"varchar(200)"` // 默认可以为null

	// 联合唯一索引(普通索引的用户一致)
	ClassNum string `xorm:"unique(location_index)"` // str默认varchar(255)
	SeatNum  string `xorm:"unique(location_index)"`

	// 忽略字段(不映射到DB)
	OmitField string `xorm:"-"`

	// 只写不读的字段
	ReadOnlyField string `xorm:"->"` // 反之，<- 为只写不读

	// create time字段
	CreateTime time.Time `xorm:"created"`

	// updated time字段
	UpdateTime time.Time `xorm:"updated"`

	// 软删除
	DelTime time.Time `xorm:"deleted"`

	// 默认值
	Note string `xorm:"default('中心小学') varchar(200)"`

	Attr `xorm:"json notnull"` // 表示先转json再入库,但其实默认会对struct转json，默认text类型

	TestMap   map[string]int // slice map 自动通过json转mysql的text类型
	TestInt   int
	TestFloat float32 `xorm:"Numeric"` // float字段只能设置为Numeric类型才能使用 = 查询，否则查不到
}

type College struct {
	Id           int64               // 学校id
	Name         string              `xorm:"varchar(50) unique notnull"`
	IconUrl      string              `xorm:"varchar(200) notnull"`
	ProvinceId   int                 `xorm:"tinyint notnull"` // 所在省份id
	Address      string              `xorm:"varchar(100) notnull"`
	Belong       string              `xorm:"varchar(50) notnull"` // 隶属于
	AcademicType int                 `xorm:"tinyint notnull"`     // 办学类型（本科、专科等）
	Category     int                 `xorm:"tinyint notnull"`     // 院校类型（医药、农林等）
	Level        int                 `xorm:"tinyint notnull"`     // 高校层次（985工程、教育部直属等）
	Tag          []string            `xorm:"json"`                // 学校标签，数组类型(mysql-text)，如 [“985”, “211”]
	Rank         int                 `xorm:"smallint notnull"`    // 全国排名
	Email        []string            `xorm:"json"`                // 联系邮箱， 数组
	Phone        []string            `xorm:"json"`                // 联系电话， 数组
	Site         map[string]string   `xorm:"json"`                // json(mysql-text) {"home":xx, "zhaosheng":xx},分别是官网和招生网
	Desc         string              `xorm:"text notnull"`        // 学校介绍
	Speciality   map[string][]string `xorm:"json"`                // 专业信息，json，{}
	CreateTime   time.Time           `xorm:"created"`             // Insert时自动赋值为当前时间
	UpdateTime   time.Time           `xorm:"updated"`             // Insert或Update时自动赋值为当前时间
}

// insert into prefix_students(id,union_pk,usr_name,attr) values(1,1,磊,attr)
type Student struct {
	UserBase `xorm:"extends"` // 将匿名结构体的字段加载过来
}

// 实现结构体自己的表字段转换规则
//func (s *Students) FromDB([]byte) error{
//
//}
//
//func (s *Students) ToDB() ([]byte, error){
//
//}

func BaseFunc() {
	var err error
	engine, err = xorm.NewEngine("mysql", "root:123@/app_web?charset=utf8")
	if err != nil {
		log.Fatal(111, err)
	}
	engine.ShowSQL(true)

	// 日志写入文件
	f, err := os.OpenFile("sql.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(222, err)
	}
	engine.SetLogger(xorm.NewSimpleLogger(f))
	engine.Logger().SetLevel(core.LOG_DEBUG) // 设置日志级别

	engine.SetMaxIdleConns(1)
	engine.SetMaxOpenConns(1)

	/*结构体与表及字段的映射规则*/
	engine.SetMapper(core.GonicMapper{})       // 设置映射规则 GonicMapper不会将id转换为i_d, SnakeMapper会
	engine.SetTableMapper(core.GonicMapper{})  // 默认和上面一致
	engine.SetColumnMapper(core.GonicMapper{}) // 默认和上面一致

	// 带前缀的映射方案 （针对表）
	tbMapper := core.NewPrefixMapper(core.GonicMapper{}, "prefix_")
	engine.SetTableMapper(tbMapper)

	// note:当某个表名不想按映射规则来时，也可以为结构体定义TableName() string 方法来自定义表名
	// 字段名的映射还可通过结构体中字段的tag指定，如xorm:"'column_name'"
}

func syncTable() {
	engine.ShowSQL()
	err := engine.Charset("utf8mb4").Sync2(new(College))
	//err := engine.Sync2(new(Student))
	if err != nil {
		log.Fatal(err)
	}
}

func TestXormBase(t *testing.T) {
	//s,_ := time.Now().MarshalJSON()
	//fmt.Println(string(s))
	//var ti = &time.Time{}
	//_ = ti.UnmarshalJSON(s)

	BaseFunc()
	syncTable()
}

func Test_CURD(t *testing.T) {

	var student = Student{}

	// ****************INSERT
	//student.Id = 1 // 主键自增 无需指定
	student.Name = "lei01"

	// 联合索引字段
	student.ClassNum = "12"
	student.SeatNum = "3"

	student.Desc = "hello lei"
	student.Gender = "male"
	student.IsTeacher = false
	student.LikeFruits = []string{"apple", "pear"}
	student.TestMap = map[string]int{"A": 1, "B": 2}
	student.ReadOnlyField = "readOnlyF"
	student.OmitField = "notRead"
	student.Note = "nothing"
	student.TestInt = 1

	// 插入成功后，主键以及createTime等字段将被回写到student结构内
	affected, err := engine.Insert(&student)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("SINGLE INSERT", affected)

	// BULK INSERT
	stu02 := student
	stu02.Id = 0
	stu02.SeatNum = "4"
	stu02.Name = "lei02"
	stu03 := student
	stu03.Id = 0
	stu03.SeatNum = "5"
	stu03.Name = "lei03"

	// 两种批量插入方式
	// 注意各个数据库对SQL长度限制，一次性批量插入的SQL语句过长会导致插入失败，
	// 一般150条记录以下是安全的，如大于，需要自己切分多个150条分批插入
	//affected, err = engine.Insert(&stu02, &stu03)
	stus := []*Student{&stu02, &stu03}
	affected, err = engine.Insert(stus)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("BULK INSERT", affected)
	// mysql支持一次性批量插入，这会使得id无法回写到结构体内
	// 不支持一次性批量插入的就会自动一条条的插入，这种可以回写id
	log.Printf("BULE INSERT id1:%d id2:%d\n", stu02.Id, stu03.Id)

	// ****************SELECT
	// 声明一个变量接收数据，单个结构体接收1条记录，即使返回了多条。
	var stu_select = Student{}
	ok, err := engine.Alias("a").Where("a.usr_name = ? and a.class_num = ?", "lei01", "12").Get(&stu_select)
	if !ok {
		log.Println("SELECT NULL")
	} else {
		log.Printf("SELECT SUCC,following\n%+v\n", stu_select)
	}

	// 多条查询
	var stus_select = []*Student{}
	err = engine.Alias("a").Where("a.class_num = ?", "12").Find(&stus_select)
	if err != nil {
		log.Println(err)
		return
	} else {
		log.Printf("SELECT SUCC,size:%d\n", len(stus_select))
	}

	// ****************UPDATE
	var stu_update = new(Student)
	stu_update.Name = "lei100"
	stu_update.TestInt = 0 // nil和 0的字段不会被更新
	affected, err = engine.Where("id = ?", 1).Update(stu_update)
	if err != nil {
		log.Println("UPDATE err:", err)
		return
	} else {
		log.Printf("UPDATE SUCC,size:%d\n", affected)
	}
	// 使用·Col·更新 某个字段为0
	stu_update.Name = "lei100"
	stu_update.TestInt = 0
	affected, err = engine.Where("id = ?", 1).Cols("test_int", "name").Update(stu_update)
	if err != nil {
		log.Println("UPDATE zero err:", err)
		return
	} else {
		log.Printf("UPDATE zero SUCC,size:%d\n", affected)
	}

	// ****************DELETE
	// 如果某个字段tag是 `xorm:"deleted"`，删除将不会真正进行，只是记录删除时间，同时会被正常查询过滤
	var stu_del = new(Student)
	affected, err = engine.Id(3).Delete(stu_del)
	if err != nil {
		log.Println("DEL err:", err)
		return
	} else {
		log.Printf("DEL SUCC,size:%d\n", affected)
	}
}

func Test_RawSqlQuery(t *testing.T) {
	sql := "select * from prefix_student limit 2"
	ret, err := engine.Query(sql)
	if err != nil {
		log.Fatal(err)
	}
	if len(ret) >= 1 {
		for _, item := range ret {
			for k, v := range item {
				log.Printf("%s -- %s\n", k, string(v))
			}

		}
	} else {
		log.Println("No data")
	}
}
