package main

func main() {
	// 16进制表示
	//println(0x11)  // 0001 0001
	println(0x303) // 00000001(高8位) 00000010(低8位)
	println(makeWord(3, 3))
	println(LoByte(0x102), HiByte(0x102))
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
