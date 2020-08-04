package testify

import (
	"github.com/chaseSpace/go-common-pkg-exmaples/testify/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

/*
使用mock tool生成测试代码，节省我们的时间(这是另一个库)
https://github.com/vektra/mockery

我这里下载开发版本以便演示，正确的应该是下载发布的版本
go get github.com/vektra/mockery/cmd/mockery
*/

// 定义具有特定行为的测试对象
type SubstractOwner interface {
	Substract(a, b int) int
}

/*
$cd $CURR_DIR
$mockery -name=substractOwner    // 生成文件 mocks/substractOwner.go, name是interface名，通过正则搜索
*/

// 待测试的函数
func targetFuncThatUseSubstractOwner(s SubstractOwner, a, b int) int {
	r := s.Substract(a, b)
	if r < 0 {
		return 0
	}
	return r
}

func TestSubstractOwner(t *testing.T) {

	testObj := &mocks.SubstractOwner{}
	testObj.On("Substract", 1, 1).Return(2)

	r := targetFuncThatUseSubstractOwner(testObj, 1, 1)
	testObj.AssertExpectations(t)
	assert.Equal(t, r, 2)
}
