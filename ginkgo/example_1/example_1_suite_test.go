package example_1_test

import (
	"log"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestExample1(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Example1 Suite")
}

/*
BeforeSuite和AfterSuite通常写在x_suite_test.go里面
x_test.go只存放具体的测试用例代码，因为你可能拥有x1_test.go, x2_test.go
*/

// BeforeSuite在整个测试套件运行前运行，只运行一次
var _ = BeforeSuite(func() {
	log.Println("do BeforeSuite===")
})

// AfterSuite在整个测试套件运行结束后运行，只运行一次
var _ = AfterSuite(func() {
	log.Println("do AfterSuite===")
})
