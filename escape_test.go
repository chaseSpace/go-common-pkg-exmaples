package main

import (
	"testing"
)

type S struct {
	Name string
}

func marshal2(s3 interface{}) {
	s4 := s3.(S)
	s4.Name = ""
}

//go:noinline
func marshal(s1 interface{}) {
	s2 := s1.(S)
	s2.Name = ""
	marshal2(s2)
}

func escapeType() ([]byte, []S, map[int]int, chan int, func()) {
	var b = []byte("xxx")
	var s = []S{}
	var m = map[int]int{}
	var c = make(chan int)
	var f = func() {}
	return b, s, m, c, f
}

// go test -c -gcflags "-m -m" escape_test.go
func TestEscape(t *testing.T) {
	escapeType()
	//s1 := S{}
	//marshal(s1)
	//b, _ := json.Marshal(s1)
	//log.Println(b)
}
