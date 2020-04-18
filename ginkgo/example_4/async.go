package example_4

import (
	"github.com/fatih/color"
	"io"
	"net/http"
)

var httpSvr *http.Server

func ReadMsgFromIOReader(reader io.ReadCloser) (string, error) {
	var buf = make([]byte, 10)
	_, err := reader.Read(buf)
	if err == io.EOF {
		err = nil
	}
	return string(buf), err
}

func StartHttpServer(httpAddr string) {
	var RedPrintFn = color.New(color.Bold, color.FgRed).PrintfFunc()
	http.HandleFunc("/hi", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("HI"))
		RedPrintFn("httpSvr%s recv req -- /hi\n", httpAddr)
	})
	http.HandleFunc("/close_me", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("OK"))
		err := httpSvr.Close()
		if err != nil {
			panic(err)
		}
		RedPrintFn("httpSvr%s closed!\n", httpAddr)
	})
	httpSvr = &http.Server{Addr: httpAddr}
	err := httpSvr.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}
