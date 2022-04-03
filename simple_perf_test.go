package main

import (
	"github.com/spf13/nitro"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
)

func someFunc(timer *nitro.B) {
	//buf := bytes.NewBuffer([]byte(`ssd`))
	//_ = buf
	resp := &http.Response{
		StatusCode: 303,
		Body:       io.NopCloser(strings.NewReader("111")),
	}
	_ = resp
	f, _ := os.CreateTemp(".", "*")
	timer.Step("step1")
	//b, _ := ioutil.ReadAll(resp.Body)
	f.ReadFrom(resp.Body)
	//f.Write(b)
	timer.Step("step1")
	//f.Close()
}

func TestXx(t *testing.T) {
	timer := nitro.Initalize()
	nitro.AnalysisOn = true

	someFunc(timer)
}
