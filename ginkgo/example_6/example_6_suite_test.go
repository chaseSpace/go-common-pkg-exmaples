package example_6_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestExample6(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Example6 Suite")
}
