package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"testing"
)

func TestCSV(t *testing.T) {
	f, err := os.Create("test.csv")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	var data = make([][]string, 4)
	data[0] = []string{"标题", "作者", "时间"}
	data[1] = []string{"羊皮卷", "鲁迅", "2008"}
	data[2] = []string{"易筋经", "唐生", "665"}

	f.WriteString("\xEF\xBB\xBF") // 写入一个UTF-8 BOM

	w := csv.NewWriter(f) //创建一个新的写入文件流
	w.WriteAll(data)
	w.Flush()
}
