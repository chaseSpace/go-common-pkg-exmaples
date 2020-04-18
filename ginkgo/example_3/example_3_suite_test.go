package example_3_test

import (
	"github.com/fatih/color"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestExample3(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Example3 Suite")
}

var _ = AfterSuite(func() {
	// color包方便打印颜色文字
	color.New(color.FgRed).Println("AfterSuite -- PanicTest done!")
})
