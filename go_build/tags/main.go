package main

import (
	"fmt"
)

func main() {
	fmt.Printf("vars %s\n", some_img_url)
	foo()
	bar()
	//woo() 无法调用
	fmt.Println(MONGODB_HOST)
}

/*
# 1. 不加tags
$ go build .
$ ./tags.exe
This is foo without test_tag

# 2. 加tags
$ go build -tags test_tag .
$ ./tags.exe
This is foo


小结：
1. go build通过-tags来编译满足条件的go文件
2. 在同目录下的go文件，若开头指定了+build注释，则运行存在多个go文件之间存在同名函数
：：注意，至少有一个文件使用 //+build !xxx ，这个文件会作为默认对象被使用，如果没有会报错

-tags 还可应用在 go run, go install, go test等

3. 指定多个tag
// +build tag1,tag2 tag3,!tag4
// +build tag5

>提供多个tag
go build -tags "tag1 tag2"

4. 忽略此文件
// +build ignore
添加后，goland的输入提示中不会包含此文件中的任何对象(常量变量函数等)，如果调用无法通过编译

5. 有一些go预定义的tag，有特别的含义
arm, arm64, 386, amd64, s390x : for different processor architectures
windows, darwin, linux, dragonfly, freebsd, netbsd, openbsd, plan9, solaris, android : for operating systems
cgo : when cgo is enabled
gc or gccgo : for gc and gccgo toolchains
go1.11 for Go version 1.11+, go1.12 for Go version 1.12+ etc.

这些tag可以直接在文件名上以后缀形式使用，如 bar_windows.go 和 bar_linux.go, 等效于 // +build windows/linux
编译时会被自动识别

6. 可在test时区分依赖外部资源的不依赖的测试用例：$ go test -tags "integration"，这个integration是自己定义的。
*/
