package example_2_test

import (
	"log"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestExample2(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Example2 Suite")
}

// BeforeSuite在整个测试套件运行前运行，只运行一次
var _ = BeforeSuite(func() {
	log.Println("do BeforeSuite===")
	go startLocalHttpServer()
})

// AfterSuite在整个测试套件运行结束后运行，只运行一次
var _ = AfterSuite(func() {
	if httpSvr != nil {
		_ = httpSvr.Close()
	}
	log.Println("do AfterSuite===")
})
