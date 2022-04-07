package conn

import (
	"bytes"
	"context"
	"fmt"
	"github.com/chaseSpace/go-common-pkg-exmaples/tcpcomm/socket"
	"io"
	"log"
	"syscall"
	"time"
)

type TcpConnET struct {
	sock         *socket.Socket
	reuseTempBuf []byte
	streamBuf    *bytes.Buffer
	err          error
	CloseCtx     struct {
		cancel context.CancelFunc
		context.Context
	}
}

func NewTcpConnET(s *socket.Socket) *TcpConnET {
	ctx, cancel := context.WithCancel(context.TODO())
	return &TcpConnET{
		sock:         s,
		reuseTempBuf: make([]byte, 10), // 设置一个较小的buf 方便观察到多次从底层buffer中读取流并组成一个pack的现象
		streamBuf:    bytes.NewBuffer(make([]byte, 0)),
		CloseCtx: struct {
			cancel context.CancelFunc
			context.Context
		}{cancel: cancel, Context: ctx},
	}
}

// ET模式下，需要一次读完buffer再return，再等待下次内核通知;
// LT模式下，则不需要一次读完buffer，因为buffer只要非空就会收到内核通知
// (不论哪种模式，都要自行缓存stream)

func (c *TcpConnET) ReadLoop() {
	for {
		n, err := c.sock.Read(c.reuseTempBuf)
		if err != nil {
			if n > 0 {
				log.Println(111, n)
				bs := make([]byte, n)
				copy(bs, c.reuseTempBuf[:n])
				c.streamBuf.Write(bs) // 中断了，读出来的要缓存到buf
			}
			// read buffer已读完
			if err == syscall.EAGAIN {
				select {
				case <-c.CloseCtx.Done():
					c.err = fmt.Errorf("ReadLoop: ctx.Done(), conn close")
					goto LoopEnd
				default:
				}
				// 若client不发送数据，这句会一直打印，即读空转问题！
				//log.Println("ReadLoop: no data, try next time~")
				time.Sleep(time.Millisecond * 100)
				continue
			}
			if err == syscall.EINTR {
				continue
			}
			log.Println("ReadLoop: socket.read other err", err)
			c.err = err
			goto LoopEnd
		}
		if n < 1 {
			c.err = fmt.Errorf("ReadLoop: socket.read 0, client closed")
			log.Println(c.err)
			goto LoopEnd
		}
		bs := make([]byte, n)
		copy(bs, c.reuseTempBuf[:n])
		c.streamBuf.Write(bs)
		println("ReadLoop: read stream in one loop:", c.streamBuf.String())
	}
LoopEnd:
	c.CloseCtx.cancel() // 也终止WriteLoop()
	log.Println("ReadLoop end")
	return
}

func (c *TcpConnET) WriteLoop() {
	// 检查：收到一个pack，回复一个pack（示例中，回复内容与收到内容无关）
	for {
		pack, err := c.streamBuf.ReadBytes(packDelimiter)
		if err == io.EOF {
			// 无数据可写时/接收数据的间隙时，判断ctx是否结束
			select {
			case <-c.CloseCtx.Done():
				c.err = fmt.Errorf("WriteLoop: ctx.Done(), conn close")
				log.Println(c.err)
				goto LoopEnd
			default:
			}
			if len(pack) > 0 {
				log.Println("WriteLoop: incomplete pack, write back~")
				c.streamBuf.Write(pack) // 找不到分隔符，说明是不完整的pack，要把读出来的写回去，下次再读
			}
			// 写空转问题，简单sleep处理
			time.Sleep(time.Millisecond * 100)
			continue
		}
		reply := []byte(fmt.Sprintf("server reply: [%s]", pack[:len(pack)-1]))
		reply = append(reply, packDelimiter) // to be a Pack

		_, err = c.sock.Write(reply)
		if err == syscall.EAGAIN {
			log.Println("WriteLoop：write buffer fulled, wait for next time...")
			continue
		}
		// 若write时conn被某一方关闭，必然会报错
		if err != nil {
			c.err = fmt.Errorf("WriteLoop err:%v", err)
			log.Println(c.err)
			break
		}
	}
LoopEnd:
	// 也终止ReadLoop()
	log.Println("WriteLoop: end")
}

func (c *TcpConnET) SockFd() int {
	return c.sock.Fd
}

func (c *TcpConnET) Err() error {
	return c.err
}

func (c *TcpConnET) Close() error {
	c.CloseCtx.cancel()
	c.streamBuf.Reset()
	return c.sock.Close()
}
