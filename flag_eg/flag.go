package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

/*
flag包是Go语言标准库提供用来解析命令行参数的包，使得开发命令行工具更为简单
*/

// 自定义类型
type Address string

func (p *Address) Set(s string) error {
	fixed := "San francisco"
	if s == "San francisco" {
		*p = Address(s)
		return nil
	}

	return errors.New("add只能是:" + fixed)
}

func (p *Address) String() string {
	return fmt.Sprintf("%f", *p)
}

func main() {

	help := flag.Bool("help", true, "need help")
	namePtr := flag.String("name", "username", "姓名")
	agePtr := flag.Int("age", 18, "年龄")
	musclePtr := flag.Bool("muscle", true, "是否有肌肉")

	// 这种定义变量 再传指针给flag的方式是推荐使用方式
	var email string
	flag.StringVar(&email, "email", "chenqionghe@sina.com", "邮箱")

	var hello = new(Address)
	flag.Var(hello, "add", "	hello参数")

	flag.Parse() // 开始解析

	args := flag.Args()
	fmt.Println("name:", *namePtr)
	fmt.Println("age:", *agePtr)
	fmt.Println("muscle:", *musclePtr)
	fmt.Println("email:", email)
	fmt.Println("args:", args)
	fmt.Println("add:", *hello)

	flag.Usage = usage

	if *help {
		flag.Usage() // 替换默认函数
	}
}

func usage() {
	_, _ = fmt.Fprintf(os.Stdout, "That below is what you can do!  \nOptions:\n")
	flag.PrintDefaults() // 打印用法
}

/*
go build
flag_test.exe -age=11 -email=123@qq.com -muscle=false -name="哈哈" -add="Los Angeles" 其他1 其他2
*/
