### IOCP理论
https://docs.microsoft.com/en-us/windows/win32/fileio/i-o-completion-ports

### 说明
运行此服务，必须在windows系统上，否则无法编译运行

### 状态：尚未完成...

### 概念解释

```
// 直译为重叠的意思，但不好理解。实际表达的是一次异步I/O过程中的相关状态信息
// 所有未使用的成员字段都应该初始化为零值，否函数调用将报错ERROR_INVALID_PARAMETER
type Overlapped struct {
        // 内核使用字段: 表示I/O请求的状态码
	Internal     uintptr
        // 内核使用字段: 表示I/O请求过程中传输的bytes number
	InternalHigh uintptr
	// 启动I/O请求的文件位置的低位部分，用户指定；仅当在支持offset概念的搜索设备上执行I/O请求时，此字段才非零
	Offset       uint32
	// 启动I/O请求的文件位置的高位部分，用户指定；仅当在支持offset概念的搜索设备上执行I/O请求时，此字段才非零
	OffsetHigh   uint32
	// 操作完成时系统将设置为信号状态的事件句柄，用户使用0或CreateEvent()返回值进行初始化
	HEvent       Handle
}
```