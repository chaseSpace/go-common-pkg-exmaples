package example_2_test

import (
	"github.com/chaseSpace/go-common-pkg-exmaples/ginkgo/example_2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"log"
	"net/http"
	"time"
)

var httpSvr *http.Server
var httpAddr = ":8080"

func startLocalHttpServer() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(404)
		_, _ = writer.Write([]byte("404"))
	})
	http.HandleFunc("/200", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(200)
		writer.Header().Add("location", "http://")
		_, _ = writer.Write([]byte("200"))
	})
	httpSvr = &http.Server{Addr: httpAddr}
	err := httpSvr.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}

var _ = Describe("Request", func() {
	// 声明struct用来包含一个测试用例需要的参数和返回值
	type argWithReturnVal struct {
		url               string
		timeout           time.Duration
		rspBodyContainStr string
		statusCode        int
	}

	// Describe就像个容器，里面包含了一系列具有类似行为的测试用例
	Describe("test sum", func() {
		// Context用来给用例分类, 比如不同的参数导致函数有不同的行为
		// 这个分类针对的是statusCode
		Context("test statusCode", func() {
			// BeforeEach嵌入到Context里面
			// 就表示这个BeforeEach只会在这个Context下的测试用例运行之前执行
			// 同样的，BeforeEach可以嵌入到 Describe/Context/When语句块内部
			// 这几个方法的作用都是一样的，只有方法名不同
			// 为了方便我们组织多个测试用例
			BeforeEach(func() {
				log.Println("BeforeEach within test statusCode")
			})
			It("should be 200 for code and contain substr", func() {
				// 不需要共享的变量可以放到测试用例函数内
				normalTest := argWithReturnVal{
					url:               "https://www.baidu.com",
					rspBodyContainStr: "baidu.com",
					statusCode:        200}
				rsp, err := example_2.Get(normalTest.url)
				var buf = make([]byte, 1000)
				Expect(err).Should(BeNil())
				Expect(rsp.StatusCode).Should(Equal(normalTest.statusCode))
				_, _ = rsp.Body.Read(buf)
				defer rsp.Body.Close()
				Expect(buf).Should(ContainSubstring(normalTest.rspBodyContainStr))
			})

			It("should be 404 for code", func() {
				// 不需要共享的变量可以放到测试用例函数内
				normalTest := argWithReturnVal{
					url:        "http://localhost" + httpAddr,
					statusCode: 404}
				rsp, err := example_2.Get(normalTest.url)
				Expect(err).Should(BeNil())
				Expect(rsp.StatusCode).Should(Equal(normalTest.statusCode))
			})
		})

		// Context用来给用例分类, 这个分类是针对timeout
		Context("test timeout", func() {
			It("should be timeout", func() {
				timeoutTest := argWithReturnVal{
					url:     "https://www.youtube.com",
					timeout: time.Nanosecond}
				_, err := example_2.Get(timeoutTest.url, timeoutTest.timeout)
				Expect(err).Should(Not(BeNil()))
				//Expect(err).ShouldNot(BeNil()) 等效
				Expect(err.Error()).Should(ContainSubstring("Client.Timeout"))
			})
		})

	})
})
