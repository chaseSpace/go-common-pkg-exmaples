package example_6_test

import (
	. "github.com/chaseSpace/go-common-pkg-exmaples/ginkgo/example_6"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("IntToStr", func() {
	/*
		0. Measure和It是同级，不能相互嵌套, 但前者内部还可以包含多个b.Time()，每个表示一个用例（参考下面）
		1. 通过--skipMeasurements可以跳过性能测试用例（只运行普通的It）
		2. PMeasure和XMeasure可以标记该用例为pending，不会执行
		3. Measure的func参数拥有固定的签名，它的最后一个参数是运行次数
		4. 结果会输出b.Time()传入func运行的性能统计
	*/
	var willRunTimes int64 = 2000000 // 运行次数
	Measure("tests str to int", func(b Benchmarker) {
		// 第一个测试用例
		// 注意一个Measure内的多个b.Time用例的name不同冲突，负责结果只会打印一个
		costTime := b.Time("WithStrConvItoA", func() {
			Expect(WithStrConvItoA(10000)).To(Equal("10000"))
		}, "this is a test for WithStrConvItoA")
		// 默认打印结果的时间单位是s，但是一般我们需要的是ns，所以重新记录这个值
		b.RecordValueWithPrecision("WithStrConvItoA-ns", float64(costTime.Nanoseconds()), "ns", 2)
		// 对用例运行的总时间进行断言(可选)
		//Expect(costTime.Seconds()).To(BeNumerically("<", 3), "WithStrConvItoA should not take too long")

		// 第二个测试用例
		costTime = b.Time("WithFmtSprintf", func() {
			Expect(WithFmtSprintf(10000)).To(Equal("10000"))
		}, "this is a test for WithFmtSprintf")
		b.RecordValueWithPrecision("WithFmtSprintf-ns", float64(costTime.Nanoseconds()), "ns", 2)

		// 第三个测试用例
		costTime = b.Time("WithFormatInt", func() {
			Expect(WithFormatInt(10000)).To(Equal("10000"))
		}, "this is a test for WithFormatInt")
		b.RecordValueWithPrecision("WithFormatInt-ns", float64(costTime.Nanoseconds()), "ns", 2)

	}, int(willRunTimes))
})
