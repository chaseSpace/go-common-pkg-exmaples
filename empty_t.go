package main

import (
	"bytes"
)

func main() {
	b := bytes.Buffer{}
	b.WriteString("x111111111111111")
	println(b.Len())
}
