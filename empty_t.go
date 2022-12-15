package main

import (
	"fmt"
)

func main() {
	// 16进制表示
	//println(0x11)  // 0001 0001
	println(0x303) // 00000001(高8位) 00000010(低8位)
	println(makeWord(3, 3))
	println(LoByte(0x102), HiByte(0x102))

	println(reverseBits(0xf0000000)) // => 0Xf = 15
	println(^(int8(-12) - 1))

	fmt.Println(GetIntimateLvByExp(20778))

	fmt.Printf("%p %p", get(), get())
}

type x struct {
}

var s = &x{}

func get() *x {
	return s
}

// 实现Cpp中的同名宏函数
func makeWord(lowByte, highByte int8) uint16 {
	// 低8位，高8位
	low := uint16(lowByte)
	high := uint16(highByte) * 1 << 8
	return low + high
}

// LoByte 取一个双字节数据最低（最右边）字节的内容，作用同c++同名函数
func LoByte(x uint16) uint16 {
	x &= 0x000F // 高位清零就得到低位
	return x
}

// HiByte 取一个双字节数据最高（最左边）字节的内容，作用同c++同名函数
func HiByte(x uint16) uint16 {
	x &= 0xFF00 // 低位清零就得到高位
	return x
}

func reverseBits(n uint32) uint32 {
	n2 := uint32(0)
	for i := 0; i < 32; i++ {
		n2 = n2<<1 | (n & 1)
		n >>= 1
	}
	return n2
}

func GetIntimateLvByExp(exp int32) (lv, currLvExp, nextExp int32) {
	lv2Exp := int32(520)
	lv3Exp := int32(520 + 1314)
	lv4Exp := int32(520 + 1314 + 3344)
	lv5Exp := int32(520 + 1314 + 3344 + 5200)

	switch _ = 0; {
	case exp < lv2Exp:
		return 1, 0, lv2Exp
	case exp < lv3Exp:
		return 2, lv2Exp, lv3Exp
	case exp < lv4Exp:
		return 3, lv3Exp, lv4Exp
	case exp < lv5Exp:
		return 4, lv4Exp, lv5Exp
	default:
		lv = (exp-lv5Exp)/5200 + 5
		currLvExp = lv5Exp + (lv-5)*5200
		nextExp = (lv-5+1)*5200 + lv5Exp
		return
	}
}
