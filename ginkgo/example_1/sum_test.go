package example_1_test

import (
	"github.com/chaseSpace/go-common-pkg-exmaples/ginkgo/example_1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"log"
)

/*
. "github.com/xxx" 这样的import语法是直接导入目标库的命名空间
在使用其对象时不需要输入xxx.method/func, 而是直接输入method/func
这只是个go的import语法，并不是要求必须这样做，但却是常见做法
*/

var _ = Describe("Sum", func() {
	// 声明变量 用来在 Describe和BeforeEach/AfterEach等函数块中共享
	// 这些变量一般作为测试函数的参数和返回值，结合下面It中的内容就懂了

	type testItem struct {
		slice       []int
		expectedSum int
	}
	// 通常会在开始处声明一些变量，用来共享给所有测试用例
	var (
		runningWhichCase int
		item             *testItem
	)

	// BeforeEach在每个测试用例开始前执行，可以用来做初始化工作
	BeforeEach(func() {
		// 每个测试用例都能够得到相同的初始item
		item = &testItem{
			slice:       []int{1, 2, 3},
			expectedSum: 6,
		}
		runningWhichCase++
		log.Printf("running test NO.%d begin\n", runningWhichCase)
	})
	// BeforeEach在每个测试用例完成后执行，可以用来做资源清理工作
	AfterEach(func() {
		item = nil
		log.Printf("running test NO.%d end\n", runningWhichCase)
	})

	// Describe就像个容器，里面包含了一系列具有类似行为的测试用例
	// It()的第一个参数text表示对这个用例的预期描述，格式一般是`should xxx`
	// 这样在观察日志时也很方便
	Describe("test sum", func() {
		// It函数包含了一个最小粒度的测试用例，其中可以包含测试代码和断言
		// 不能再往里面嵌套It
		It("should be equal sum(vals...) with expectedSum1", func() {
			s, _ := example_1.Sum(item.slice)
			// Expect是gomega提供的断言方法
			// 它有一个等效的方法叫 `Ω`,这是个符号（中文输入法输入oumu会出来），用哪个都行
			// To等效于Should, Equal是gomega提供的匹配方法，还有BeNil/BeTrue/BeFalse等等
			Expect(s).To(Equal(item.expectedSum))
		})
		// Specify等效于It
		// 使用It/Specify是为了方便贴合自然语言
		// 你可以将方法名与第一个text参数连起来阅读
		Specify("Sum should return EmptySliceErr", func() {
			item.slice = []int{}
			_, err := example_1.Sum(item.slice)
			Expect(err).To(MatchError(example_1.EmptySliceErr))
		})
	})
})
