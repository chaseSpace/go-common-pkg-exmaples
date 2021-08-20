package main

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"
)

/*
本文件内的代码演示了unsafe.Pointer的基本使用方法；
如要了解更高级的使用指针做算数运算，需要先*熟练*掌握go的内存对齐原则，即使这样还是尽量少用指针做算数运算，
因为go版本升级可能会导致结构体的内存布局发生微小变化，然后导致你的算数运算出现偏差，导致程序crash~
# 比如reflect.SliceHeader结构体的注释就表名了，结构体的实现在将来可能会发生转变

进阶查看同目录下文件：unsafe2_test.go
*/

func TestUnsafe(t *testing.T) {
	// =====示例-1   了解*int（泛指所有指针变量）、Pointer、uintptr的转换用法
	v1 := int(1)
	// 基本法则
	// *int => Pointer => uintptr => Pointer => *int
	v2 := (*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&v1))))
	println(v1 == *v2) // true

	// =====示例-2   了解uintptr如何做运算 (先得了解struct的内存分配规则)
	var x = struct {
		a int
		b byte
		c []int
	}{a: 1, b: 'x', c: []int{2}}

	// OffsetOf获取成员a相对内存中结构体起始位置的偏移量，单位bytes
	// 所谓的指针运算就是将 针对于 指针面量 的算术运算，这里是充分利用了对象在内存中的非常规则有序的分配原则
	// 每个变量在内存中的地址都可用一个指针面量来表示，go中用表现为uintptr类型，打印是一个无符号整数
	pt := unsafe.Pointer(uintptr(unsafe.Pointer(&x)) + unsafe.Offsetof(x.a))
	a1 := (*int)(pt)
	println(x.a == *a1) // true

	// =====示例-3   不同类型间的强转换
	// go中有些类型在内存中的结构是一致的，可以通过指针运算做0拷贝转换
	// 比如str 和 bytes, 通常我们使用[]byte(x) 来转str，这会造成内存拷贝
	s := "x"
	b := []byte{'x'}
	// str的内存结构是reflect.StringHeader, bytes是reflect.SliceHeader
	//
	strHead := (*reflect.StringHeader)(unsafe.Pointer(&s))
	sliceHead := reflect.SliceHeader{
		Data: strHead.Data,
		Len:  strHead.Len, // strHead的len也是byte单位
		Cap:  strHead.Len,
	}
	newB := (*[]byte)(unsafe.Pointer(&sliceHead))
	println((*newB)[0] == b[0]) // true
	// 其实str和bytes间的转换可以一步到位
	newB = (*[]byte)(unsafe.Pointer(&s))
	newS := (*string)(unsafe.Pointer(&b))
	println(string(*newB) == *newS)

	// ======示例-4    不能使用uintptr临时变量
	z := "x"
	zptr := uintptr(unsafe.Pointer(&z))
	fmt.Printf("zptr: %#x", zptr) // 0xc00011df08, 这表示此刻 &z的指针面量
	// &z 和 其指针面量是两个概念，前者在程序运行期间，栈内存迁移时也可以指向正确数据
	// 后者是一个固定地址，在内存布局变动后，该地址指向的位置可能没有数据也可能指向其他数据
	// 	这时再将其强转为源类型，会导致不可捕获的panic
	// 安全的做法是不要声明uintptr临时变量，而是在指针运算时直接使用

	// 这里模拟将一个错误的指针面量强转为指针 (IDE会警告用法不对)
	_err_pointer_not_panic_now := (*string)(unsafe.Pointer(zptr + 1))
	defer func() { _ = recover() }()          // 无法recover
	_panic_now := *_err_pointer_not_panic_now // fatal error: unexpected signal during runtime execution
	println(_panic_now)

	/*
		总结：unsafe包下的几个函数 都是用来获取对象的内存排列位置信息的；
		unsafe包的用途&场景：
			-	用于实现零拷贝转换类型
			-	偏底层、高频调用的模块内使用
		用法不对会产生无法recovery的panic！！！
	*/
}
