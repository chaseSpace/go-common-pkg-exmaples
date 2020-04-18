
![image](https://onsi.github.io/ginkgo/images/ginkgo.png)

Ginkgo allows you to write tests in Go using expressive Behavior-Driven Development ("BDD") style.

### [What is BDD ?](https://zh.wikipedia.org/wiki/行为驱动开发)

### Learn more about [Ginkgo](https://github.com/onsi/ginkgo)

### [中文文档](https://blog.gmem.cc/ginkgo-study-note)

---

### Getting started
```bash
go get -u github.com/onsi/ginkgo/ginkgo  # installs the ginkgo CLI
go get -u github.com/onsi/gomega/...     # fetches the matcher library

cd path/to/package/you/want/to/test

ginkgo bootstrap # set up a new ginkgo suite
ginkgo generate  # will create a sample test file.  edit this file and add your tests then...

go test # to run your tests

ginkgo  # also runs your tests

ginkgo -v # be equivalent to `go test -ginkgo.v`, -v represents -verbose 
```

### CLI
```bash
# 运行当前目录中的测试
ginkgo
# 运行其它目录中的测试
ginkgo /path/to/package /path/to/other/package ...
# 递归运行所有子目录中的测试
ginkgo -r ...
```

**Pass options by cli**

option | desc
---|---
--reportPassed | 打印测试通过的报告
--v | 冗长模式
--trace | 打印所有错误的调用栈
--progress | 打印进度
-race | 开启竞争条件检测
-cover | 覆盖率测试
-tags | 指定编译期标记
-timeout | 套件超时时间，默认24h0m0s, eg: ginkgo -timeout 1m1s