package main

import (
	"bytes"
)

func main() {
	b := bytes.Buffer{}
	b.WriteString("1231111111111111111111111111231111111111111111111111111211111111")
	println(b.Len())
	s, _ := b.ReadString('x')
	println(b.Len(), b.Cap())
	b.WriteString(s)
	println(b.Len())
}
