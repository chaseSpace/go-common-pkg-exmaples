package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	type S struct {
		A int
	}
	s := S{128}
	b, _ := json.Marshal(&s) // {"A":128}, 128占3个字节
	println(len(b))          // 9bytes

	println('\x80', len([]byte{'\x80'})) // 128 1   -- 80是十进制数128的十六进制表示
	fmt.Print(s)
}
