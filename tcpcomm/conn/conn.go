package conn

import (
	"bytes"
	"context"
	"fmt"
	"github.com/chaseSpace/go-common-pkg-exmaples/tcpcomm/socket"
	"io"
	"log"
	"syscall"
)

type Pack []byte

// 这里规定一个pack以\n结束
const packDelimiter = '\n'

type TcpConn struct {
	sock         *socket.Socket
	reuseTempBuf []byte
	streamBuf    *bytes.Buffer
	err          error
	CloseCtx     struct {
		cancel context.CancelFunc
		ctx    context.Context
	}
}

func NewTcpConn(s *socket.Socket) *TcpConn {
	return &TcpConn{
		sock:         s,
		reuseTempBuf: make([]byte, 10), // 设置一个较小的buf 方便观察到多次从底层buffer中读取流并组成一个pack的现象
		streamBuf:    bytes.NewBuffer(make([]byte, 100)),
	}
}

// ET模式下，需要一次读完buffer再return，再等待下次内核通知;
// LT模式下，则不需要一次读完buffer，因为buffer只要非空就会收到内核通知
// (不论哪种模式，都要自行缓存stream)

func (c *TcpConn) Read() {
	for {
		n, err := c.sock.Read(c.reuseTempBuf)
		if err != nil {
			if n > 0 {
				bs := make([]byte, n)
				copy(bs, c.reuseTempBuf[:n])
				c.streamBuf.Write(bs) // 中断了，读出来的要缓存到buf
			}
			// read buffer已读完
			if err == syscall.EAGAIN {
				break
			}
			if err == syscall.EINTR {
				continue
			}
			log.Println("socket.read other err", err)
			c.err = err
			return
		}
		if n < 1 {
			c.err = fmt.Errorf("socket.read 0, client closed")
			log.Println(c.err)
			return
		}
		bs := make([]byte, n)
		copy(bs, c.reuseTempBuf[:n])
		c.streamBuf.Write(bs)
		println("read stream in one loop:", c.streamBuf.String())
		// o(╥﹏╥)o遗留的问题：
		// 经测试，ET模式下，即使没有一次性读完底层buffer的数据，还是会持续的收到read event（与LT模式无差了）
		// 使得本例仍能跑通~~ 作者暂无法找到原因
		//break
	}
	log.Println("read stream end")
	return
}

func (c *TcpConn) WriteReply() {
	// 检查：收到一个pack，回复一个pack（示例中，回复内容与收到内容无关）
	for {
		pack, err := c.streamBuf.ReadString(packDelimiter)
		if err == io.EOF {
			if pack != "" {
				c.streamBuf.WriteString(pack) // 找不到分隔符，说明是不完整的pack，要把读出来的写回去，下次再读
				log.Println("incomplete pack, write back~")
			}
			break
		}
		reply := []byte(fmt.Sprintf("server reply: [%s]", pack[:len(pack)-1]))
		reply = append(reply, packDelimiter) // to be a Pack

		_, err = c.sock.Write(reply)
		if err == syscall.EAGAIN {
			log.Println("WriteReply：write buffer fulled, wait for next time...")
			return
		}
		if err != nil {
			c.err = fmt.Errorf("WriteLoop err:%v", err)
			log.Println(c.err)
			return
		}
		log.Println(string(reply[:len(reply)-1]))
	}
	log.Println("WriteReply: end")
}

func (c *TcpConn) WriteReplyET() {
	// TODO 修改
	// 检查：收到一个pack，回复一个pack（示例中，回复内容与收到内容无关）
	for {
		pack, err := c.streamBuf.ReadString(packDelimiter)
		if err == io.EOF {
			if pack != "" {
				c.streamBuf.WriteString(pack) // 找不到分隔符，说明是不完整的pack，要把读出来的写回去，下次再读
				log.Println("incomplete pack, write back~")
			}
			break
		}
		reply := []byte(fmt.Sprintf("server reply: [%s]", pack[:len(pack)-1]))
		reply = append(reply, packDelimiter) // to be a Pack

		_, err = c.sock.Write(reply)
		if err == syscall.EAGAIN {
			log.Println("WriteReply：write buffer fulled, wait for next time...")
			return
		}
		if err != nil {
			c.err = fmt.Errorf("WriteLoop err:%v", err)
			log.Println(c.err)
			return
		}
		log.Println(string(reply[:len(reply)-1]))
	}
	log.Println("WriteReply: end")
}

func (c *TcpConn) SockFd() int {
	return c.sock.Fd
}

func (c *TcpConn) Err() error {
	return c.err
}

func (c *TcpConn) Close() error {
	c.streamBuf.Reset()
	return c.sock.Close()
}
