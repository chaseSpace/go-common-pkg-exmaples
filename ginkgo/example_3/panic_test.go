package example_3_test

import (
	. "github.com/chaseSpace/go-common-pkg-exmaples/ginkgo/example_3"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"log"
)

/*
预期测试结果：FAIL! -- 1 Passed | 2 Failed | 0 Pending | 0 Skipped
*/

var _ = Describe("Describe-PanicTest", func() {

	Describe("Describe-test panic", func() {
		It("should panic", func() {
			// 有时候It的text参数不能简短的介绍复杂的用例内容，就用到By，它可以在It内部多次使用，用来当做日志一下
			// 比如用例包含4个步骤，那么在对应步骤的起始处添加By("description for this step...")
			// 文本会在错误时打印出来
			By("Document your test case-0,  but this case won't display, because it will be passed for test")
			// 测试panic，需要用func封装
			Expect(func() { DoPanic(true) }).Should(Panic())
		})
		// 只要不是函数内新创建的goroutine中的panic不会中断整个测试进程
		// 始终可以观察到AfterSuite的执行
		It("should not panic", func() {
			By("Document your test case-1, It will be convenient for you to get the notes of the use case when it fails.",
				func() {
					log.Println("you can pass to By with a fun(), it will immediately run.")
				})
			Expect(func() { DoPanic(true) }).ShouldNot(Panic())
		})

		// 新创建的goroutine抛出的panic会导致整个测试停止，测试进程立即停止
		// 如果需要创建新的goroutine来运行用例，需要defer GinkgoRecover()
		// func(done Done){ ... close(done)} 是固定写法
		// >>为什么要在新的goroutine中运行用例，因为默认It方法默认同步执行，这样就可以异步
		// done管道是用来与goroutine通信的, 参考example_4
		It("should not panic in a goroutine", func(done Done) {
			By("Document your test case-2, It will be convenient for you to get the notes of the use case when it fails.",
				func() {
					log.Println("you can pass to By with a fun(), it will immediately run.")
				})
			go func() {
				defer GinkgoRecover() // 不加这行就会导致测试进程中断，后面的测试不会继续执行
				Expect(func() { DoPanic(true) }).ShouldNot(Panic())
				close(done)
			}()
		})
	})
})
