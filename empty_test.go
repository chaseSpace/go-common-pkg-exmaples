package main

import (
	"fmt"
	"net/url"
	"testing"
)

func TestWs(t *testing.T) {
	var durl = url.URL{
		Scheme:   "http",
		Host:     "1.1.1.1:22",
		Path:     "/api/v1/sadasd",
		RawQuery: fmt.Sprintf("fid=%s&filename=%s", "filename", fmt.Sprintf("新增统计数据.csv"))}
	fmt.Println(durl.String())
}
