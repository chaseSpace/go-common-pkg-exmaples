package example_2

import (
	"net/http"
	"time"
)

/*
1. create request.go
2. CMD:ginkgo bootstrap (got a example_1_request_test.go)
3. CMD:ginkgo generate sum (got a request_test.go)
*/
var defaultTimeout = 20 * time.Second

func Get(url string, timeout ...time.Duration) (*http.Response, error) {
	var (
		timeOut = defaultTimeout
		rsp     *http.Response
		err     error
	)
	if len(timeout) > 0 {
		timeOut = timeout[0]
	}
	req := http.Client{Timeout: timeOut}
	rsp, err = req.Get(url)
	return rsp, err
}
