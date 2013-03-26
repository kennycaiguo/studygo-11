package main

import (
	"net/mail"
	"net/smtp"
	"encoding/base64"
	"fmt"
)

func main(){
	host := "192.168.0.243:25"

	from := mail.Address{"发件人", "service@tlt.cn"}
	to := mail.Address{"收件人", "tzm529@163.com"}
	header := make(map[string]string)
	header["From"] = from.String()
	header["To"] = to.String()
	b64 := base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")
	header["Subject"] = fmt.Sprintf("=?UTF-8?B?%s?=", b64.EncodeToString([]byte("标题测试")))
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=UTF-8"
	header["Content-Transfer-Encoding"] = "base64"

	body := "邮件正文"

	message := ""

	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + b64.EncodeToString([]byte(body))
	auth := smtp.PlainAuth("", "", "", host)
	err := smtp.SendMail(host, auth, "service@tlt.cn", []string{to.Address}, []byte(message))
	fmt.Println(err)
}
