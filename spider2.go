package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"io"
	"log"
	"net/url"
	"strings"
)

// go get -u github.com/parnurzeal/gorequest

type loginRsp struct {
	Success bool   `json:"success"`
	Msg     string `json:"msg"`
}

// 【人事】列表数据
type RenshiRsp struct {
	Rows []struct {
		Data []interface{} `json:"data"`
		Id   string        `json:"id"`
	} `json:"rows"`
	Pos        int `json:"pos"`
	TotalCount int `json:"total_count"`
}

func createFormReader(data map[string]string) io.Reader {
	form := url.Values{}
	for k, v := range data {
		form.Add(k, v)
	}
	return strings.NewReader(form.Encode())
}

func getBody(body io.ReadCloser) string {
	b, _ := io.ReadAll(body)
	return string(b)
}

func prettyPrint(v interface{}) {
	b, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(b))
}

func main() {
	// 请求前两个URL是调用登录接口前的必须操作，什么NT设计，content还是固定的
	portalPwdEffectiveCheckUrl := "https://58.23.12.98:8034/portalPwdEffectiveCheck.do"
	form1 := map[string]string{"content": "YJS03Pf1Vilj3NNPtPm8W2zDGJocVmvaTbCCi0yfJANCyYqwAp.2F7Ap7BCEbd2ZvAnP14N8kVtnoh.2BPrWf9Jyj9toPRPVUtP4Nds.2FamdGNxwqs7PYR6IpGqoh7d0YfqoM"}
	portalCheckStrengthUrl := "https://58.23.12.98:8034/portalCheckStrength.do"
	form2 := map[string]string{"strength": "-1"}
	// 登录接口
	loginUrl := "https://58.23.12.98:8034/login.do"
	form3 := map[string]string{"username": "test", "password": "098f6bcd4621d373cade4e832627b4f6", "loginType": "NORMAL"}
	// 确认真实session的URL
	dashboardUrl := "https://58.23.12.98:8034/dashboard.do?dashboard"

	req := gorequest.New().TLSClientConfig(&tls.Config{InsecureSkipVerify: true}).Set("Content-Type", "application/x-www-form-urlencoded")

	_rsp, _bytes, err := req.Clone().Post(portalPwdEffectiveCheckUrl).Type(gorequest.TypeForm).Send(form1).EndBytes()
	if err != nil {
		log.Fatal(err)
	}
	println("S1 OK!", _rsp.Header.Get("Set-Cookie"))

	_rsp, _bytes, err = req.Clone().Post(portalCheckStrengthUrl).Type(gorequest.TypeForm).Send(form2).Set("Cookie", _rsp.Header.Get("Set-Cookie")).EndBytes()
	if err != nil {
		log.Fatal(err)
	}
	println("S2 OK!", string(_bytes))

	rsp := new(loginRsp)
	_rsp, _bytes, err = req.Clone().Post(loginUrl).Type(gorequest.TypeForm).Set("Cookie", _rsp.Header.Get("Set-Cookie")).Send(form3).EndStruct(rsp)
	if err != nil {
		log.Fatal(err)
	}
	if !rsp.Success {
		println("登录失败，返回:", string(_bytes))
	}
	println("登录成功！", _rsp.Header.Get("Set-Cookie"))

	_rsp, _bytes, err = req.Clone().Get(dashboardUrl).Set("Cookie", _rsp.Header.Get("Set-Cookie")).EndBytes()
	if err != nil {
		log.Fatal(err)
	}
	if _rsp.Request.URL.Path == "/bioLogin.do" { // 302
		println("得到了假session, FK!!! *******")
		return
	}
	println("真的登录成功！")

	// 尝试获取人事列表数据
	page := 1
	rspS := new(RenshiRsp)
	form4 := map[string]string{"posStart": fmt.Sprintf("%d", (page-1)*50), "count": "50", "pageSize": "50", "list": ""}
	_rsp, _bytes, err = req.Clone().Post("https://58.23.12.98:8034/persPerson.do").Send(form4).Set("Cookie", _rsp.Header.Get("Set-Cookie")).EndStruct(rspS)
	if err != nil {
		log.Fatal(err)
	}
	//prettyPrint(rspS)
	if rspS.TotalCount > 0 {
		println("获取人事列表成功！")
	}
}
