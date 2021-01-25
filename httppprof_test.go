package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"testing"
	"time"
)

func Test_HttpPprof(t *testing.T) {
	go func() {
		for {
			log.Println(Add("https://github.com/EDDYCJY"))
			time.Sleep(time.Millisecond * 100)
		}
	}()

	http.ListenAndServe("0.0.0.0:6060", nil)
}

/*
命令行执行
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
等待采样30s（URL设置）
可用命令
- 	help
-	top10   flat(当前函数【采样过程中】平均耗时)  flat%（当前函数平均耗时占比）   sum%（当前函数累计耗时占比）
			cum（cumulative,当前函数及上层调用平均耗时）   cum%（当前函数及上层调用平均耗时占比）
      flat  flat%   sum%        cum   cum%
      20ms 22.22% 22.22%       20ms 22.22%  runtime.stdcall1
      10ms 11.11% 33.33%       10ms 11.11%  fmt.(*pp).printArg
      10ms 11.11% 44.44%       10ms 11.11%  runtime.cgocall
      10ms 11.11% 55.56%       10ms 11.11%  runtime.execute
      10ms 11.11% 66.67%       10ms 11.11%  runtime.exitsyscall
      10ms 11.11% 77.78%       20ms 22.22%  runtime.goready.func1
      10ms 11.11% 88.89%       10ms 11.11%  runtime.heapBitsSetType
      10ms 11.11%   100%       10ms 11.11%  time.Now
         0     0%   100%       10ms 11.11%  fmt.(*pp).doPrintln
         0     0%   100%       10ms 11.11%  fmt.Sprintln

*/
