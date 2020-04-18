package example_4_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestExample4(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Example4 Suite")
}
