package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func main() {
	http.HandleFunc("/hi", func(writer http.ResponseWriter, request *http.Request) {
		defer request.Body.Close()
		b, _ := ioutil.ReadAll(request.Body)
		_, _ = writer.Write([]byte(fmt.Sprintf(`{"received": "%s"}`, string(b))))
	})
	defer log.Println("simple_http exited!")
	log.Println("simple_http started...")

	portStart := 8080
	httpSvr := &http.Server{Addr: fmt.Sprintf(":%d", portStart)}
	for {
		err := httpSvr.ListenAndServe()
		if err != nil && strings.Contains(err.Error(), "already in use") {
			portStart++
			httpSvr.Addr = fmt.Sprintf(":%d", portStart)
			continue
		}
		log.Println("err:", err)
		break
	}
}

/*
$curl -d 'xxx' -i http://localhost:8080/hi
HTTP/1.1 200 OK
Date: Mon, 07 Sep 2020 02:42:37 GMT
Content-Length: 19
Content-Type: text/plain; charset=utf-8

{"received": "xxx"}
*/
