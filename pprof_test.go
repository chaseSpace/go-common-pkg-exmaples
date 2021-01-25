package main

import "testing"

const URL = "https://github.com/EDDYCJY"

var datas []string

func Add(str string) string {
	data := []byte(str)
	sData := string(data)
	datas = append(datas, sData)

	return sData
}

func TestAdd(t *testing.T) {
	s := Add(URL)
	if s == "" {
		t.Errorf("Test.Add error!")
	}
}

func BenchmarkAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Add(URL)
	}
}

/*
go test -run=^$ -bench=. -cpuprofile=cpu.prof
-- 生成cpu.prof文件后启动一个http服务查看分析图
go tool pprof -http=:8080 cpu.prof

查看火焰图：http://localhost:8080/ui/flamegraph
*/
