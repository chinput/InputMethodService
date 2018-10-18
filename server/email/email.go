package email

import (
	"net/smtp"
	"regexp"
	"strings"
	"time"

	"github.com/chinput/InputMethodService/server/config"
)

type EmailType struct {
	subject     string
	common_data string
}

const (
	KTypeRegisterCode  = 0
	KTypeResetPassword = 1
)

var (
	common_subject = "[多多输入法云服务]"
	emailType      = []EmailType{
		{
			subject:     common_subject + "注册码",
			common_data: "您正在注册" + common_subject + "，以下是您的注册码（15分钟内有效）：",
		}, {
			subject:     common_subject + "重置密码",
			common_data: "您正在充值" + common_subject + "的用户密码，以下是您的验证码：",
		},
	}
)

var emailRegexp = regexp.MustCompile(".+@.+\\..+")

func IsEmailValid(m string) bool {
	return emailRegexp.MatchString(m)
}

type OneEmail struct {
	Type int
	Data string
	To   string
}

type emailQueue struct {
	que chan *OneEmail
}

func (e *emailQueue) Send(emailType int, key_data string, to string) {
	e.que <- &OneEmail{
		Type: emailType,
		Data: key_data,
		To:   to,
	}
}

// 设置为最快多长时间发送一次
// 最低 10 ms
func (e *emailQueue) Guard() {
	max := len(emailType)
	_, _, _, timeout := config.EmailConfig()

	if timeout < 10 {
		timeout = 10
	}

	for {
		one := <-e.que
		if one.Type < max {
			cfg := emailType[one.Type]
			sendEmail(one.To, cfg.subject, cfg.common_data, one.Data)
		}
		<-time.After(time.Millisecond * time.Duration(timeout))
	}
}

var que *emailQueue

func init() {
	que = new(emailQueue)
	que.que = make(chan *OneEmail, 64)
	go que.Guard()
}

func SendRegisterCode(to, code string) {
	que.Send(KTypeRegisterCode, code, to)
}

func SendResetCode(to, code string) {
	que.Send(KTypeResetPassword, code, to)
}

func sendEmail(to, subject, common_data, key_data string) error {
	h, u, p, _ := config.EmailConfig()

	body := common_data
	if key_data != "" {
		body += "<br/><strong style=\"font-size:2em\">" + key_data + "</strong>"
	}

	return sendToMail(u, p, h, to, subject, body, "html")
}

func sendToMail(user, password, host, to, subject, body, mailtype string) error {
	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])
	var content_type string
	if mailtype == "html" {
		content_type = "Content-Type: text/" + mailtype + "; charset=UTF-8"
	} else {
		content_type = "Content-Type: text/plain" + "; charset=UTF-8"
	}

	msg := []byte("To: " + to + "\r\nFrom: " + user + ">\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)
	send_to := strings.Split(to, ";")
	err := smtp.SendMail(host, auth, user, send_to, msg)
	return err
}
