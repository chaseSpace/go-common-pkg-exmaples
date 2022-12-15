package main

import (
	"log"
	"math/rand"
	"testing"
	"time"
)

// 奖品配置列表
var prize = []struct {
	id   int
	prob int // 权重
}{
	{id: 3, prob: 30}, {id: 4, prob: 40}, {id: 1, prob: 10}, {id: 2, prob: 20},
}

// 抽奖函数
func probability() int {
	// sum是上面配置的权重总和
	sum := 100
	// 返回的是奖品的index
	for i := 0; i < len(prize); i++ {
		if rand.Intn(sum) < prize[i].prob {
			return i
		}
		sum -= prize[i].prob
	}
	// 最后返回概率最大的奖品，但永远不会走到这一步
	return -1
}

func TestX(t *testing.T) {
	rand.Seed(time.Now().Unix())

	var statis = []struct {
		id     int
		number int
	}{ // 这里的顺序记得和上面配置一致
		{id: 3, number: 0}, {id: 4, number: 0}, {id: 1, number: 0}, {id: 2, number: 0},
	}
	// 抽N次 统计结果
	for i := 0; i < 1; i++ {
		idx := probability()
		statis[idx].number++
	}
	// [{1 100111} {2 199661} {3 300041} {4 400187}]
	log.Println(statis)

	var s []int
	s = append(s, 1, 2, 3)
}
