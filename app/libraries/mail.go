package libraries

import (
	"github.com/astaxie/beego"
	"net/smtp"
	"strings"
)

func SendMail(to, subject, body, mailType string) error {
	host := beego.AppConfig.String("MailHost")
	auth := smtp.PlainAuth("", beego.AppConfig.String("MailAccount"), beego.AppConfig.String("MailPassword"), host)
	var contentType string
	if mailType == "html" {
		contentType = "Content-Type: text/" + mailType + "; charset=UTF-8"
	} else {
		contentType = "Content-Type: text/plain" + "; charset=UTF-8"
	}
	msg := []byte("To: " + to + "\r\nFrom: " + beego.AppConfig.String("MailUserName") + "\r\nSubject: " + subject + "\r\n" + contentType + "\r\n\r\n" + body)
	sendTo := strings.Split(to, ";")
	err := smtp.SendMail(host+":"+beego.AppConfig.String("MailPort"), auth, beego.AppConfig.String("MailAccount"), sendTo, msg)
	return err
}
