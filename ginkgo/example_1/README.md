1. create sum.go
2. CMD:ginkgo bootstrap (got a example_1_suite_test.go)
3. CMD:ginkgo generate sum (got a sum_test.go)
4. CMD:go test # OR ginkgo, that is equal to `go test -run=TestExample1`

套件测试表示对这个目录的全部代码做测试，所以这个目录下只能有一个xxx_suite_test.go
如果要对新增的代码文件做测试，使用 ginkgo generate xxx生成xxx_test.go,然后在生成的xxx_test.go中
写测试用例即可

执行套件测试(所有测试用例):  
```bash
$ go test # OR ginkgo  只打印失败的用例详情
$ go test -ginkgo.v #ginkgo -v 会打印所有用例的执行结果（成功和失败）
```

输出
```text
=== RUN   TestExample1
Running Suite: Example1 Suite
=============================
Random Seed: 1587117037
Will run 2 of 2 specs

2020/04/17 17:50:37 do BeforeSuite===
++2020/04/17 17:50:37 do AfterSuite===

Ran 2 of 2 Specs in 0.048 seconds
SUCCESS! -- 2 Passed | 0 Failed | 0 Pending | 0 Skipped
--- PASS: TestExample1 (0.05s)
PASS
```