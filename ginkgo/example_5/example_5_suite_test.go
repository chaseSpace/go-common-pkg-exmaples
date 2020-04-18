package example_5_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestExample5(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Example5 Suite")
}
