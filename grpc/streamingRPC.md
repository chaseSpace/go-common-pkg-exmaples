
### 流式RPC
-	用于大量数据包持续发送的场景

普通的RPC是REQ<->RSP的模式，但是一次来回也涉及到序列化和反序列化的开销，流式RPC将多个REQ/RSP压缩发送，
提供发送效率，减少解压次数，并且在一次请求响应过程中，客户端可以一次发送多个REQ，服务端也可以一次接收多个RSP；
普通RPC是来一个请求就启动一个goroutine，流式RPC是在一个goroutine中处理多个REQ/RSP；

就类似MQ的处理模式了

参考代码：
https://eddycjy.com/posts/go/grpc/2018-09-24-stream-client-server

备用链接：
https://github.com/EDDYCJY/go-grpc-example

