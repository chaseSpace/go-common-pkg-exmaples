package example_5_test

import (
	. "github.com/chaseSpace/go-common-pkg-exmaples/ginkgo/example_5"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"
)

/*
输出效果：
SUCCESS! -- 0 Passed | 0 Failed | 1 Pending | 1 Skipped
*/
var _ = Describe("Skip", func() {
	// PDescribe和Describe的区别是前者将块内的测试用例标记为pending
	// pending的测试用例不会执行
	// PDescribe == PContext == PIt == PMeasure = XIt == XMeasure == XDescribe == XContext
	// `==` 表示它们都会将各自语句块内的用例标记为pending
	// 默认都会打印skip信息，可以 --noisySkippings=false 来取消打印
	PDescribe("Skip block", func() {
		It("skip", func() {
			Expect(SkipFunc).ShouldNot(Panic())
		})
	})

	var skip = time.Now().Unix()%2 == 0
	Describe("Skip dynamically", func() {
		It("should skip", func() {
			if skip {
				// 第二种方法是动态判断你设置的条件选择是否跳过（使用Skip()）
				// Skip后不需要显式return（注意确保你的代码是在It块内）
				Skip("condition is true, so this case will skip")
			}
			Expect(SkipFunc).ShouldNot(Panic())
		})
	})
})
