package main

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"
)

/*
内存对齐遵循的原则：
1. Struct中每个成员分配的起始地址必须是其大小的倍数
2. Struct大小一定是其最大成员类型大小的倍数
golang的内存对齐特点：
-	除了遵循上述两点外，对于嵌套结构体，不会当做单独成员，而是展开其内部的子成员
	所以嵌套结构体不会导致整个结构体内存占用倍增，下面ss内部的N3  N4就是例子

-------------------------------------
#扩展知识#
# 什么是内存对齐
	-	简单来说，这是由于CPU对内存访问的特点引出的，大部分架构的CPU在通过一个变量首地址访问内存时，都【建议】这个地址是N的倍数(N也叫对齐边界/系数)，
		而这个要求就叫做内存对齐，只要是内存对齐的首地址，CPU都可以高效的通过该地址访问内存数据；若没有对齐，有些架构的CPU直接
		不支持读取这个地址。抛出硬件异常；而另一些架构的CPU做了兼容处理，它们将一次非对齐访问拆分为两次对齐访问，因为每次访问的
		数据大小与寄存器单元大小一致，即32位CPU的单个寄存器单元大小（也称字长word size）为4byte，那么一次内存访问就一定读取的是4byte，假设有个非对齐变量是4byte，
		CPU就得分两次对齐访问，得到8byte，然后取出其中4byte。（这里对CPU来说就多做功了，耗时耗力）

# 为什么需要内存对齐
	-	1. 平台原因(移植原因)：不是所有的硬件平台都能访问任意地址上的任意数据的；某些硬件平台只能在某些地址处取某些特定类型的数据，否则抛出硬件异常。
	-	2. 性能原因：数据结构(尤其是栈)应该尽可能地在自然边界上对齐。原因在于，为了访问未对齐的内存，处理器需要作两次内存访问；而对齐的内存访问仅需要一次访问。

# 内存对齐分2种情况：整体对齐和结构体成员对齐
	# 1.整体对齐
		- 针对的是任意类型（包含结构体、基本类型、指针等等），要求任意类型的变量都需要进行内存对齐
			任意类型的变量的首地址（的十进制）必须是对齐系数N的倍数

	# 2.结构体成员对齐
		- 针对结构体，要求其每个成员内存对齐，具体规则是：每个成员的变量的首地址（的十进制）必须是对齐系数N的倍数

# 如何求得 N ？（以下规则适用绝大部分架构CPU）
	- 这里定义 x=CPU字长，y=具体类型大小，那么N=min(x, y)
		规则一：若是指针类型，y=8
		规则二：若是结构体类型，y=max(成员类型大小)

	go中使用unsafe.AlignOf获取变量的对齐系数N，这个N是常数，编译期确定！

	下面基于64位CPU架构进行举例：
		-	var a int8; unsafe.AlignOf(a)==1
		-	var a byte; unsafe.AlignOf(a)==1
		-	var a int16; unsafe.AlignOf(a)==2
		-	var a int64; unsafe.AlignOf(a)==8

	# 下面是go中的切片示例，因为str和切片的runtime实体是结构体，既然是结构体就应用上面的规则二，具体下面的代码有说明
		-	var a string; unsafe.AlignOf(a)==8
		-	var a []int64; unsafe.AlignOf(a)==8
		-	var a []string; unsafe.AlignOf(a)==8

# 我可以不对齐或者修改对齐系数吗？
	- 当然可以，有些CPU也是支持非对齐访问的，只是慢了点
	- 可以修改，只是go应用层不支持，需修改编译器代码；而足够open的C语言是直接支持的，通过在代码首行添加：#pragma pack(n) // n是对齐系数

# 知道某个类型变量的对齐系数后我可以做什么？
	- 哦天哪，你可以做的太多了；go提供了unsafe包。就允许我们可以绕过类型系统直接操作内存地址
		假如两个不同go类型在内存中的布局大小是一致或部分一致的
		我们都可以将类型a的变量数据（或截取部分）直接映射（构造）为类型b的变量数据，注意是零拷贝！！！
	-	比如我们经常看到有些包内部会将[]byte 通过unsafe调用转换为 string
	-	在go应用层来说，我们可以指鹿为马！

# 深入探索：为什么CPU访问内存的首地址一定是对齐的？
	- 偏向硬件层了，自行查询了解
*/

type ss struct {
	N0 bool // offset 0
	N1 bool // offset 1
	N2 bool // offset 2
	/*
		切片 等于 reflect.SliceHeader（结构体）
	*/
	N3 []int               // offset 8  size 24
	N4 reflect.SliceHeader // offset 32 size 24
} // sizeof ss = 56

type ss1 struct {
	N0 bool // offset 0
	N1 bool // offset 1
	N2 bool // offset 2
	/*
		指针类型的大小是uintptr 8bytes
	*/
	N3 *[]int               // offset 8  size 8
	N4 *reflect.SliceHeader // offset 16 size 8
	N5 *string              // offset 24 size 8
	N6 *int                 // offset 32 size 8
} // sizeof ss1 = 40

func TestUnsafe2(t *testing.T) {
	s := ss{}
	fmt.Println(unsafe.Offsetof(s.N0), // 1
		unsafe.Offsetof(s.N1), // 1
		unsafe.Offsetof(s.N2), //
		unsafe.Offsetof(s.N3),
		unsafe.Offsetof(s.N4))

	fmt.Println("s:", unsafe.Sizeof(s), unsafe.Offsetof(s.N3))

	s1 := ss1{}
	fmt.Println("s1:", unsafe.Sizeof(s1), unsafe.Offsetof(s1.N4), unsafe.Offsetof(s1.N5))

	// 特例：空str的size是16，因为 str的运行时实体是 reflect.StringHeader{}，后者的size就是16bytes
	fmt.Println("sizeof str", unsafe.Sizeof(""))                              // 16
	fmt.Println("sizeof StringHeader", unsafe.Sizeof(reflect.StringHeader{})) // 16
	fmt.Println("sizeof byte", unsafe.Sizeof(byte('x')))                      // 1
	fmt.Println("sizeof int", unsafe.Sizeof(int(1)))                          // 8
	fmt.Println("sizeof int8", unsafe.Sizeof(int8(1)))                        // 1

	/*
		   切片的size返回的切片描述符的大小，所有切片的初始Size都是24；
			其实所有切片的运行时表现都是 reflect.SliceHeader{}，而它的size就是24（因为struct内有3个8字节字段）
	*/
	fmt.Println(unsafe.Sizeof(reflect.SliceHeader{})) // 24
	fmt.Println(unsafe.Sizeof([]byte{}))              // 24
	fmt.Println(unsafe.Sizeof([]string{}))            // 24
	fmt.Println(unsafe.Sizeof([]int{}))               // 24
	fmt.Println(unsafe.Sizeof([]ss{}))                // 24

	// 数组size==元素类型大小*长度
	fmt.Println(unsafe.Sizeof([0]int{})) // 0
	fmt.Println(unsafe.Sizeof([1]int{})) // 8
	fmt.Println(unsafe.Sizeof([2]int{})) // 16

	/*
		AlignOf 查看各种类型变量的对齐边界（编译器提前定义），
		任意类型的切片的对齐边界都是8，这个函数主要用于指针的算数运算
	*/
	println("str", unsafe.Alignof(""))                       // 8
	println("byte", unsafe.Alignof(byte('x')))               // 1
	println("int", unsafe.Alignof(int(1)))                   // 8
	println("int8", unsafe.Alignof(int8(1)))                 // 1
	println("int16", unsafe.Alignof(int16(1)))               // 2
	println("int32", unsafe.Alignof(int32(1)))               // 4
	println("[]str", unsafe.Alignof([]string{}))             // 8
	println("[]byte", unsafe.Alignof([]byte{}))              // 8
	println("[]int16", unsafe.Alignof([]int16{}))            // 8
	println("[]int32", unsafe.Alignof([]int32{}))            // 8
	println("[]int64", unsafe.Alignof([]int64{}))            // 8
	println("[]uint64", unsafe.Alignof([]uint64{}))          // 8
	println("struct.bool", unsafe.Alignof(s.N1))             // 1
	println("struct.[]int", unsafe.Alignof(s.N2))            // 8
	println("struct.bool_after_[]int", unsafe.Alignof(s.N3)) // 8

	/*
		指鹿为马
	*/
	var b = []byte{232, 156, 156, 233, 155, 170, 226, 157, 164, 229, 134, 176, 229, 159, 142}
	_str := *(*string)(unsafe.Pointer(&b))
	fmt.Println("指鹿为马:", _str) // 蜜雪❤冰城

	/*
		偷梁换柱-struct版
	*/
	var a struct {
		n0 int
		n1 string
		n2 bool
		n3 struct {
			_    interface{}
			n3_1 []rune
		}
	}
	a.n3.n3_1 = []rune{'*', '‿', '*'}

	_n2 := (*bool)(unsafe.Pointer(uintptr(unsafe.Pointer(&a)) + unsafe.Offsetof(a.n2)))
	println("偷梁换柱之前的a.n2", a.n2) // false
	*_n2 = true
	println("偷梁换柱之后的a.n2", a.n2) // true

	fmt.Println("偷梁换柱之前的a.n3_1", string(a.n3.n3_1)) // *‿*
	// 再修改n3_1
	_n3_1 := (*[]rune)(unsafe.Pointer(uintptr(unsafe.Pointer(&a.n3)) + unsafe.Offsetof(a.n3.n3_1)))
	(*_n3_1)[1] = '︿'
	fmt.Println("偷梁换柱之后的a.n3_1", string(a.n3.n3_1)) // *︿*

	/*
		偷梁换柱-slice版
	*/
	// OffsetOf是编译器帮我计算出的位置，掌握内存对齐后我自己也可以计算出成员的位置（通过unsafe.Alignof）
	var slice1 = []int{6, 7, 6}
	shDataPtrVal := (*reflect.SliceHeader)(unsafe.Pointer(&slice1)).Data            // step1:先拿到底层数组Data指针面量
	_slice1Idx2Ptr := (*int)(unsafe.Pointer(shDataPtrVal + unsafe.Alignof(int(0)))) // step2:再计算下标为2的元素的指针面量，再转为类型变量
	*_slice1Idx2Ptr = 6                                                             // step3:改值
	fmt.Println("偷梁换柱-slice版: ", slice1)                                            // [6 6 6]
	// ↑↑↑ 备注：通过goland编写【偷梁换柱-slice版】示例时，你会看到上述代码中step2的代码会出现IDE的阴影提示，报告用法可能不对
	// 对的IDE没说错，实际环境下是将step1-2合成一行代码，不会分开写（为了演示），这是因为 shDataPtrVal 是uintptr类型，是一个指针面量的数字；
	// 而不是指向数据的指针本体，运行时由于栈扩容，可能指针指向的数据地址可能会变化，但uintptr是数字不会跟着变化，所以不能定义uintptr类型的临时变量

	// 一行代码
	_slice1Idx2Ptr = (*int)(unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&slice1)).Data + unsafe.Alignof(int(0))))
	*_slice1Idx2Ptr = 66
	fmt.Println("偷梁换柱-slice版-2: ", slice1) // [6 66 6]
}
