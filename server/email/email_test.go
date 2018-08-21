package email

import (
	"fmt"
	"testing"

	"github.com/chinput/InputMethodService/server/config"
)

/*
func Test_Send(t *testing.T) {
	config.Init("./config.toml")
	h, u, p, _ := config.EmailConfig()

	data := "您正在注册多多输入法云服务，以下是你的注册码：<br/><strong style=\"font-size:2em\">ABCDE</strong>"

	log.Println(h, u, p)
	err := sendToMail(u, p, h, "garfeng_gu@163.com", "hello world", data, "html")
	log.Println(err)
}
*/
func Test_IsEmailValid(t *testing.T) {
	config.Init("./config.toml")

	fmt.Println(IsEmailValid("hhh"),
		IsEmailValid("abc@163.com"), IsEmailValid("abc@abc"))
}
