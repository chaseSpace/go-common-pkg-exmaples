package main

import (
	"testing"
	"unsafe"
)

/*
-- 测试目的：证明常用的bytes与str之间的转换写法是有内存分配操作的，而str的底层原本就是[]byte，所以应当存在直接转换，
无内存分配的类型互转方法。
-- 测试步骤：
##### 首先是 str2bytes 的两种转换方法比较
命令：go test -bench=Str2Bytes -run=^$ -benchmem -memprofile memprofile.out
输出：
$ go test -bench=BenchmarkBytes2Str  -run=^$ -benchmem -memprofile memprofile.out
goos: windows
goarch: amd64
pkg: github.com/chaseSpace/go-common-pkg-exmaples
BenchmarkBytes2Str-8            100000000               11.9 ns/op             3 B/op          1 allocs/op
BenchmarkBytes2Str_2-8          2000000000               0.22 ns/op            0 B/op          0 allocs/op
PASS
ok      github.com/chaseSpace/go-common-pkg-exmaples    2.051s
结果说明：[]byte("xxx")的方式有一次内存分配操作，第二种方法则零内存分配开销，每次操作耗时对比 11.9 ns VS  0.22ns
进一步查看有内存分配的代码行，刚才已经导出了基准测试的内存操作记录，现在使用go自带的pprof分析工具查看：
$ go tool pprof memprofile.out
Type: alloc_space
Time: Mar 11, 2021 at 2:16pm (CST)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) list BenchmarkStr2Bytes   // <---------------- 解读：list 使用正则匹配“BenchmarkStr2Byte”的函数名，并展示匹配函数每一行的内存分配开销(有才标识)
Total: 744.03MB
ROUTINE ======================== github.com/chaseSpace/go-common-pkg-exmaples.BenchmarkStr2Bytes in D:\develop\go\go-common-pkg-exmaples\bytesconvstr_test.go
  743.01MB   743.01MB (flat, cum) 99.86% of Total
         .          .     49:
		 .          .     50:func BenchmarkStr2Bytes(b *testing.B) {
		 .          .     51:   var bin []byte
		 .          .     52:   a := "aaa"
		 .          .     53:   for i := 0; i < b.N; i++ {
  743.01MB   743.01MB     54:           bin = []byte(a)  // 注意：这个743MB是累计分配的内存大小，不是单次
		 .          .     55:           // 这里需注意，不能简单写成 _ = []byte(a)，经测试，这种写法在基准测试中看不到内存分配次数增加，
		 .          .     56:           //但操作时间仍大于 BenchmarkStr2Byte_2 中的方法
		 .          .     57:   }
		 .          .     58:   _ = bin
		 .          .     59:}
(pprof)
结果说明：显然，list BenchmarkStr2Byte 这个命令是可以匹配到第二个方法（BenchmarkStr2Byte_2）的，但是结果没有展示，这是因为该函数没有额外的内存分配操作；
虽然说 a := "aaa" 也是有内存分配操作，但是go应该是过滤了这种字面量声明式的行。

##### 其次是 bytes2str 的两种转换方法比较
这里可以使用上面的方法进行测试，不再表述
*/
func BenchmarkStr2Bytes(b *testing.B) {
	var bin []byte
	a := "aaa"
	for i := 0; i < b.N; i++ {
		bin = []byte(a)
		// 这里需注意，不能简单写成 _ = []byte(a)，经测试，这种写法在基准测试中看不到内存分配次数增加，
		//但操作时间仍大于 BenchmarkStr2Byte_2 中的方法
	}
	_ = bin
}

func BenchmarkStr2Bytes_2(b *testing.B) {
	var bin []byte
	a := "aaa"
	for i := 0; i < b.N; i++ {
		bin = *(*[]byte)(unsafe.Pointer(&a))
	}
	_ = bin
}

func BenchmarkBytes2Str(b *testing.B) {
	var a string
	bin := []byte{'a', 'b', 'c'}
	for i := 0; i < b.N; i++ {
		a = string(bin)
	}
	_ = a
}

// go test -bench=Bytes2Str -run=^$ -benchmem
func BenchmarkBytes2Str_2(b *testing.B) {
	var a string
	bin := []byte{'a', 'b', 'c'}
	for i := 0; i < b.N; i++ {
		a = *(*string)(unsafe.Pointer(&bin))
	}
	_ = a
}
