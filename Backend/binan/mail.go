package binan

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

func PostEmail() {
	m := gomail.NewMessage()
	//m.SetHeader("From", m.FormatAddress("pchaibo@163.com", "zhang"))
	m.SetHeader("From", "zhang"+"<pchaibo@163.com>")
	m.SetHeader("To", "386378183@qq.com")
	m.SetHeader("Subject", "COIN")
	html := `Hello <b>Bob</b> <a href="http://www.baidu.com/55">登录</a> and <i>Cora</i>!`
	m.SetBody("text/html", html)
	//m.Attach("./ss.docx") //附件

	d := gomail.NewDialer("smtp.163.com", 465, "pchaibo", "CJFMKMYLVZZTWKTY")

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		fmt.Println("err:", err)
	}

	//	c.String(200, "postemail"+c.ClientIP())
}
