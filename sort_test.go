package main

import (
	"log"
	"sort"
	"testing"
)

func TestSort(t *testing.T) {
	var S = []struct {
		integer int
	}{
		{integer: 1}, {2}, {6},
	}

	// 降序
	sort.Slice(S, func(i, j int) bool {
		// 小于=升序，大于=降序（小升大降）
		return S[i].integer > S[j].integer
	})

	log.Println(S)
}
