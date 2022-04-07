
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
write 28 err:<nil>
server reply: [client sent 1]
server reply: [client sent 2]
write 28 err:<nil>
server reply: [client sent 2]
server reply: [client sent 3]
```

同时会观察到server端程序日志如下：
```
root@4d455a521d36:/home/gocode/tcpsvr_in_epollLT# go run .
Server started. Waiting for incoming connections. ^C to exit.
2022/04/07 01:47:45 eventLoop new 1 events ...
2022/04/07 01:47:45 event: new Conn
2022/04/07 01:47:45 eventLoop new 1 events ...
2022/04/07 01:47:45 event: Readable fd:5
read stream in one loop:: client sen
read stream in one loop:: client sent 0
client
read stream in one loop:: client sent 0
client sent 1

2022/04/07 01:47:45 read stream end
2022/04/07 01:47:45 eventLoop new 1 events ...
2022/04/07 01:47:45 event: Writeable
2022/04/07 01:47:45 server reply: [client sent 0]
2022/04/07 01:47:45 server reply: [client sent 1]
2022/04/07 01:47:45 WriteReply: end
2022/04/07 01:47:46 eventLoop new 1 events ...
2022/04/07 01:47:46 event: Readable fd:5
read stream in one loop:: client sen
read stream in one loop:: client sent 1
client
read stream in one loop:: client sent 1
client sent 2

2022/04/07 01:47:46 read stream end
2022/04/07 01:47:46 eventLoop new 1 events ...
2022/04/07 01:47:46 event: Writeable
2022/04/07 01:47:46 server reply: [client sent 1]
2022/04/07 01:47:46 server reply: [client sent 2]
2022/04/07 01:47:46 WriteReply: end
2022/04/07 01:47:47 eventLoop new 1 events ...
2022/04/07 01:47:47 event: Readable fd:5
read stream in one loop:: client sen
read stream in one loop:: client sent 2
client
read stream in one loop:: client sent 2
client sent 3

2022/04/07 01:47:47 read stream end
2022/04/07 01:47:47 eventLoop new 1 events ...
2022/04/07 01:47:47 event: Writeable
2022/04/07 01:47:47 server reply: [client sent 2]
2022/04/07 01:47:47 server reply: [client sent 3]
2022/04/07 01:47:47 WriteReply: end
```

### 备注
本项目仅作为使用epoll(LT模式)实现tcp server的最小参考示例，相关逻辑并不是最优解

### 参考资料
https://man7.org/linux/man-pages/man2/