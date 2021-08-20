package main

import (
	"github.com/chaseSpace/go-common-pkg-exmaples/ignore"
	"gopkg.in/gomail.v2"
)

func do() error {
	m := gomail.NewMessage()
	m.SetHeader("To", ignore.EmailRecver) //收件人列表
	//m.SetHeader("Cc", ...) // 抄送
	m.SetAddressHeader("From", ignore.EmailSender, "Li")
	m.SetHeader("Subject", "明天天气预报") // 主题
	body := `
	<ul>
	<li>Attachments</li>
	<li>Embedded images</li>
	</ul>
	`
	m.SetBody("text/html", body) // 正文
	m.Attach("go.mod")           // 附件

	// 配置发件人账户
	d := gomail.NewDialer("smtp.exmail.qq.com", 465, ignore.EmailSender, ignore.EmailPass)
	err := d.DialAndSend(m)
	return err
}

func main() {
	if err := do(); err != nil {
		panic(err)
	}
}
