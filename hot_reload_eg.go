package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

var listener net.Listener
var reloading = flag.Bool("reload", false, "")
var hi = flag.String("hi", "hello", "")
var serv = &http.Server{Addr: ":801"}

func handler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(5 * time.Second)
	_, _ = w.Write([]byte(*hi))
}

func main() {
	http.HandleFunc("/", handler)
	var err error
	flag.Parse()
	log.Printf("args: reloading:%v hi:%s | ALL:%s\n", *reloading, *hi, os.Args[1:])
	if *reloading {
		f := os.NewFile(3, "")
		listener, err = net.FileListener(f)
	} else {
		listener, err = net.Listen("tcp", serv.Addr)
	}
	if err != nil {
		log.Fatalf("err:%v", err)
		return
	}
	go func() {
		err = serv.Serve(listener)
		if err == http.ErrServerClosed {
			log.Println("reload succ!")
		} else {
			log.Fatalf("Serve fail %v", err)
		}
	}()
	OnSignal()
	log.Printf("old process quit\n")
}

func OnSignal() {
	var c = make(chan os.Signal)

	// syscall.SIGUSR2 是linux特有的
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2)

	for {
		s := <-c
		ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)
		switch s {
		case syscall.SIGINT, syscall.SIGTERM:
			_ = serv.Shutdown(ctx)
			signal.Stop(c)
			os.Exit(0)
		case syscall.SIGUSR2:
			err := reload()
			// reload 失败则不会关闭旧的server，继续监听信号（需要修复新的代码）
			if err != nil {
				log.Printf("reload fail %v\n", err)
				continue
			}
			err = serv.Shutdown(ctx)
			signal.Stop(c)
			if err != nil {
				log.Printf("serv.Shutdown fail %v\n", err)
				return
			}
			time.Sleep(300 * time.Second)
			return
		}
	}
}

func reload() error {
	tl, ok := listener.(*net.TCPListener)
	if !ok {
		return errors.New("listener is not tcp listener")
	}
	// 获取父进程 socket fd
	f, err := tl.File()
	if err != nil {
		return err
	}
	// 设置传递给子进程的参数（包含 socket 描述符）
	args := []string{"-reload", "-hi='reload ok'"}
	cmd := exec.Command(os.Args[0], args...)
	cmd.Stdout = os.Stdout         // 标准输出
	cmd.Stderr = os.Stderr         // 错误输出
	cmd.ExtraFiles = []*os.File{f} // 文件描述符
	// 新建并执行子进程
	return cmd.Start()
}

// kill -USR2 PID to reload
// kill -INT PID  to kill
