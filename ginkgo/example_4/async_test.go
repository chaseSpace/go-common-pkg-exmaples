package example_4_test

import (
	. "github.com/chaseSpace/go-common-pkg-exmaples/ginkgo/example_2"
	. "github.com/chaseSpace/go-common-pkg-exmaples/ginkgo/example_4"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("Desc-Async", func() {
	// When等效于Describe，只是语义不同
	When("test async", func() {
		It("test async", func(done Done) {
			var httpAddr = ":8081"
			var domain = "http://localhost" + httpAddr
			go func() {
				// here will block if running normally
				StartHttpServer(httpAddr)
				close(done)
			}()
			rsp, _ := Get(domain+"/hi", time.Second)
			Expect(ReadMsgFromIOReader(rsp.Body)).Should(ContainSubstring("HI"))
			rsp, _ = Get(domain+"/close_me", time.Second)
			Expect(ReadMsgFromIOReader(rsp.Body)).Should(ContainSubstring("OK"))
		})
	})
})
