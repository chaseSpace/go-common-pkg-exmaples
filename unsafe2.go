package main

/*
内存对齐遵循的原则：
1. Struct中每个成员分配的起始地址必须是其大小的倍数
2. Struct大小一定是其最大成员类型大小的倍数
golang的内存对齐特点：
-	除了遵循上述两点外，对于嵌套结构体，不会当做单独成员，而是展开其内部的子成员
	所以嵌套结构体不会导致整个结构体内存占用倍增，下面ss内部的N3  N4就是例子
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

func main() {
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

	/* 任意类型的切片的对齐边界都是8
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

	println("str", unsafe.Alignof(""))                       // 8
	println("byte", unsafe.Alignof(byte('x')))               // 1
	println("int", unsafe.Alignof(int(1)))                   // 8
	println("int8", unsafe.Alignof(int8(1)))                 // 1
	println("[]str", unsafe.Alignof([]string{}))             // 8
	println("[]byte", unsafe.Alignof([]byte{}))              // 8
	println("[]int16", unsafe.Alignof([]int16{}))            // 8
	println("[]int32", unsafe.Alignof([]int32{}))            // 8
	println("[]int64", unsafe.Alignof([]int64{}))            // 8
	println("[]uint64", unsafe.Alignof([]uint64{}))          // 8
	println("struct.bool", unsafe.Alignof(s.N1))             // 1
	println("struct.[]int", unsafe.Alignof(s.N2))            // 8
	println("struct.bool_after_[]int", unsafe.Alignof(s.N3)) // 8
}
