package encode_pkg

import (
	"bytes"
	"encoding/binary"
	"testing"
)

/*
::encoding/binary
	这个标准库是用来完成常见的数据类型(整型,布尔)与字节之间的转换的
*/

// 二进制就是字节流，这个例子演示如何把一个int类型转换为固定【两个】字节
func Test_uint16ToByte(t *testing.T) {
	initVar := int32(18888)

	// 先转换类型，实际使用时注意确保目标类型能够保存实际值，如果是byte(18888)就会有数据丢失
	twoByteVar := uint16(initVar)
	buf := bytes.NewBuffer([]byte{})
	// 以小端序规则转换为字节，这很重要，反序列的时候需保持一致
	// 什么是大端序/小端序？ http://www.ruanyifeng.com/blog/2016/11/byte-order.html
	err := binary.Write(buf, binary.LittleEndian, twoByteVar)
	if err != nil {
		t.Error(err)
	}
	// 打印二进制 [11001000 1001001]
	t.Logf("Test_uint16ToByte: %b", buf.Bytes())
}
