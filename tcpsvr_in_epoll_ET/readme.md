
### 说明
运行此服务，必须内核版本>=2.6.8的linux系统或类unix系统上，否则无法编译运行

### 测试方式
1. 启动服务：
```
>go run .
Server started. Waiting for incoming connections. ^C to exit.
```

2.使用本项目下的tcp client发送请求(不限平台)
```
>cd tcpclient
>go run .
```

观察到client端日志如下：
```
conn ok...
write 28 err:<nil>
server reply: [client sent 0]
server reply: [client sent 1]
...省略
nothing to do...
```

同时会观察到server端程序日志如下：
```
root@4d455a521d36:/home/gocode/tcpsvr_in_epollLT# go run .
Server started. Waiting for incoming connections. ^C to exit.
2022/04/05 05:20:03 eventLoop new 1 events ...
2022/04/05 05:20:03 event: new Conn
2022/04/05 05:20:03 eventLoop new 1 events ...
2022/04/05 05:20:03 event: Readable fd:5
read stream: client sen
read stream: client sent 0
client
read stream: client sent 0
client sent 1

2022/04/05 05:20:03 read stream end
2022/04/05 05:20:03 eventLoop new 1 events ...
2022/04/05 05:20:03 event: Writeable
2022/04/05 05:20:03 server reply: [client sent 0]
2022/04/05 05:20:03 server reply: [client sent 1]
2022/04/05 05:20:03 WriteReply: end
...省略
```

### 备注
本项目仅作为使用epoll(ET模式)实现tcp server的最小参考示例，相关逻辑并不是最优解

### 参考资料
https://man7.org/linux/man-pages/man2/