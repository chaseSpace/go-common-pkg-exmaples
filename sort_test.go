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
		return S[i].integer > S[j].integer
	})

	log.Println(S)
}
