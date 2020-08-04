package mock_test

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/mock"
	"testing"
)

/*
mock包提供一种机制，方便我们模拟某个对象（struct），最终测试的是使用该对象的函数
*/

// 最终测试的函数的入参（最好是定义好的interface类型）
type targetAPI interface {
	DoAdd(num int) (int, error)
}

// 待测试对象
type targetS struct {
	root   int
	maxArg int
}

var tooLargeErr = errors.New("arg too large")

// 待测试对象的行为，实际并不会运行
func (t *targetS) DoAdd(num int) (int, error) {
	if num > t.maxArg {
		return -1, tooLargeErr
	}
	return t.root + num, nil
}

// 【重要】这是要测试的函数
func targetFuncThatUseTargetS(ts targetAPI) {
	r, _ := ts.DoAdd(1)
	fmt.Println("add 1 =>", r)
	r, err := ts.DoAdd(11)
	if err != nil {
		fmt.Println("add 11 err:", err)
		return
	}
	fmt.Println("add 11 =>", r)
}

// mock 对象【模拟上面的targetS】
type MyMockedObject struct {
	mock.Mock
}

// mock方法【模拟上面的targetS的行为】
func (m *MyMockedObject) DoAdd(number int) (int, error) {
	// 这里的代码是固定的
	args := m.Called(number)
	return args.Int(0), args.Error(1)
}

func TestSomething(t *testing.T) {

	// create an instance of our test object
	testObj := new(MyMockedObject)

	// 设置预期行为的入参和出参
	testObj.On("DoAdd", 1).Return(2, nil)
	testObj.On("DoAdd", 11).Return(-1, tooLargeErr)

	// 调用测试函数
	targetFuncThatUseTargetS(testObj)

	// 断言上面的两处出入参设置都有执行
	testObj.AssertExpectations(t)

}

func TestAnyTypeArgs(t *testing.T) {
	testObj := new(MyMockedObject)

	// 上面设置的是指定参数，我们还可以设置任意参数类型
	// 注：这里的出参我们可以任意设置（需类型一致），因为目的就是测试目标函数能不能处理各种情况
	// （在实际情况中，我们也不知道【待测试对象】的行为是怎样的）
	testObj.On("DoAdd", mock.Anything).Return(1, nil)

	// 调用测试函数
	targetFuncThatUseTargetS(testObj)

	// 断言上面的两处出入参设置都有执行
	testObj.AssertExpectations(t)
}
