package main

func main() {
	foo()
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
*/
