package main

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"math"
	"testing"
)

/*
base62编码方案

-	简介
	顾名思义，它使用62个可打印字符（A-Z，a-z，0-9）来编码二进制(8bit数据)到字符

-	映射规则
	十进制的0-9 =》62进制 0-9
	十进制的10-35 =》62进制 a-z
	十进制的36-61 =》62进制 A-Z

-	适用场景
	因为base64使用的字符包含"+/"，这两个字符在URL中属于专用字符需要转义后使用（转换为%xx的形式）较为麻烦，还要其他场景也是同样的原因；
	这些场景使用base62显然更合适且简单
*/

const CODE62 = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const CODE_LENTH = 62

var EDOC = map[string]int{"0": 0, "1": 1, "2": 2, "3": 3, "4": 4, "5": 5, "6": 6, "7": 7, "8": 8, "9": 9, "a": 10,
	"b": 11, "c": 12, "d": 13, "e": 14, "f": 15, "g": 16, "h": 17, "i": 18, "j": 19, "k": 20, "l": 21, "m": 22,
	"n": 23, "o": 24, "p": 25, "q": 26, "r": 27, "s": 28, "t": 29, "u": 30, "v": 31, "w": 32, "x": 33, "y": 34,
	"z": 35, "A": 36, "B": 37, "C": 38, "D": 39, "E": 40, "F": 41, "G": 42, "H": 43, "I": 44, "J": 45, "K": 46,
	"L": 47, "M": 48, "N": 49, "O": 50, "P": 51, "Q": 52, "R": 53, "S": 54, "T": 55, "U": 56, "V": 57, "W": 58,
	"X": 59, "Y": 60, "Z": 61}

// 编码：output中的每个byte换为十进制不会超过61
func b62encode(number int) (output []byte) {
	if number == 0 {
		output = []byte{'0'}
		return
	}
	// 使用常见的进制换算方法
	for number > 0 {
		round := number / CODE_LENTH
		remain := number % CODE_LENTH
		output = append(output, CODE62[remain]) // 余数对应映射表中的一个字符
		number = round                          // 商继续下次运算
	}
	return
}

func b62decode(bytes []byte) (result int, err error) {
	for index, char := range bytes {
		v, ok := EDOC[string(char)]
		if !ok {
			return 0, fmt.Errorf("%v isn't a valid char in Base62 rule", string(char))
		}
		result += v * int(math.Pow(CODE_LENTH, float64(index)))
	}
	return
}

func TestBase62(t *testing.T) {
	number := 12893019283
	output := b62encode(number)

	number2, err := b62decode(output)
	require.Nil(t, err)
	require.Equal(t, number2, number)

	// error case
	number3, err := b62decode([]byte{'+'})
	require.EqualError(t, err, fmt.Sprintf("+ isn't a valid char in Base62 rule"))
	require.Equal(t, number3, 0)
}
