package suite_test

import (
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"
)

/*
套件测试演示代码
完整的官方示例：https://github.com/stretchr/testify/blob/master/suite/suite_test.go
*/

type MyHttpSrv struct {
	addr string
	s    *http.Server
}

func NewMyHttpSrv(addr string) *MyHttpSrv {
	return &MyHttpSrv{
		addr: addr,
		s: &http.Server{
			Addr:         addr,
			ReadTimeout:  time.Second,
			WriteTimeout: time.Second,
		}}
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("srv recv request, path:%s", r.URL.Path)
	w.Header().Set("content-type", "application/json")
	_, _ = w.Write([]byte("this is a response"))
	w.WriteHeader(200)
}

func (m *MyHttpSrv) Start() {
	http.HandleFunc("/test", testHandler)
	go func() { _ = m.s.ListenAndServe() }()
	log.Printf("http srv [started], with addr:%s -----\n", m.addr)
}

func (m *MyHttpSrv) Stop() {
	_ = m.s.Close()
	log.Printf("http srv [stopped], with addr:%s -----\n", m.addr)
}

var httpSrv *MyHttpSrv

// 定义最基本也是必须的一个对象
type SuiteObj struct {
	suite.Suite
}

// SetupSuite：在在套件测试开始前执行
func (s *SuiteObj) SetupSuite() {
	log.Println("SetupSuite---")
	httpSrv = NewMyHttpSrv("localhost:999")
	httpSrv.Start()
}

// 待测内容
func (s *SuiteObj) TestCore() {
	res, err := http.Get("http://localhost:999/test")
	if err != nil {
		log.Panic(err)
	}
	greeting, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Panic(err)
	}
	s.Equal("this is a response", string(greeting))
	s.Equal(200, res.StatusCode)
	// 其实http handler测试会更简单，可以直接使用s.HTTPSuccess()
	// 而不需要启动server，此处是为了演示套件测试
}

// TearDownSuite 在套件测试结束时执行
func (s *SuiteObj) TearDownSuite() {
	httpSrv.Stop()
	httpSrv = nil
	log.Println("TearDownSuite---")
}

func TestSuiteObj(t *testing.T) {
	suite.Run(t, new(SuiteObj))
}
