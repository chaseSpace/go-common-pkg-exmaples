package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
	"testing"
)

func TestRune(t *testing.T) {
	// 这几个unicode字符肉眼查看都是空格，但实际是不同形式的空格，这点自行查阅
	a := "\u0000\u0020\u3000\u00A0"                                  // len=4
	fmt.Println(fmt.Sprintf("[%v]", a))                              // [    ]
	fmt.Println(" " == "\u0020")                                     // true，常用的空格实际是 \u0020 这个unicode字符
	fmt.Printf("rune-len:%d \n", len([]rune(a)))                     // 多少个rune类型字符, 4
	fmt.Printf("contains space: %v\n", strings.ContainsRune(a, ' ')) // true

	// rune底层是int32即四字节类型, 在go中通常用来表示1个unicode字符，以utf8编码方案存储和传输，也可用单引号创建,
	// 但是，rune类型通过buf.WriteRune()的形式写入buf会忽略低位，所以通过buf.Bytes()得到的也是实际数据类型占的长度，而不是内存中数据占的字节长度
	// 例子：'\u0000'用utf8编码是单字节，通过buf.WriteRune()也只会写1个字节到buf，但是rune类型在内存中占的字节长度是int32,即4个字节（下面通过binary库可一探究竟）

	var buf bytes.Buffer
	buf.WriteRune(' ')      // 单字节
	buf.WriteRune('\u0000') // 在unicode中表示空字符，单字节
	buf.WriteRune('\u0020') // 表示ASCII的空格，编号在unicode定义的0x00-0x7F范围内，以utf8编码仅需单字节即可表示
	buf.WriteRune('\u4F60') // 表示汉字的`你`，在unicode定义的0x0800-0xFFFF之间，以utf8编码需三字节可表示
	buf.WriteRune(0x20C30)  // 0x20C30是unicode字符的另一种16进制书写形式，在unicode定义的0x010000-0x10FFFF之间，以utf8编码需四字节可表示
	fmt.Printf("rune-bin:%b rune equal: %v, byte len of a rune:%d "+
		"toStr:[%s]\n", ' ', ' ' == rune(' '), len(buf.Bytes()), buf.Bytes()) // rune equal: true, byte len of a rune:10 toStr:[   你𠰰]

	// 看看一个rune类型的【空字符】在内存中到底占了多少字节？
	var b2 = &bytes.Buffer{}
	err := binary.Write(b2, binary.LittleEndian, ' ')
	if err != nil {
		panic(err)
	}
	fmt.Printf("binary rune: %b actual-byte-len:%d \n", b2.Bytes(), len(b2.Bytes())) // [100000 0 0 0] actual-byte-len:4

	// 单引号 就是 rune 类型的另一种写法， 而byte是uint8
	fmt.Printf("%T %T %T\n", ' ', rune(' '), byte(' ')) // int32 int32 uint8

}
