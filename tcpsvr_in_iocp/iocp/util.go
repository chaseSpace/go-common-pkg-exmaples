package iocp

// MakeWord 传入一个低位和高位字节数据参数，合并为一个双字节数据，作用同C++同名宏函数
func MakeWord(lowByte, highByte int8) uint16 {
	return uint16(highByte)<<8 | uint16(lowByte)
}

// LoByte 取一个双字节数据最低（最右边）字节的内容，作用同C++同名函数
func LoByte(x uint16) uint16 {
	x &= 0x000F // 高位清零就得到低位
	return x
}

// HiByte 取一个双字节数据最高（最左边）字节的内容，作用同C++同名函数
func HiByte(x uint16) uint16 {
	x &= 0xFF00 // 低位清零就得到高位
	return x >> 8
}
