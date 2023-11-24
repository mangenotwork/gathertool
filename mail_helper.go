/*
*	Description : 邮件的方法，应用场景有抓取完成通知. TODO 测试
*	Author 		: ManGe
*	Mail 		: 2912882908@qq.com
**/

package gathertool

import (
	"gopkg.in/gomail.v2"
)

type Mail struct {
	From string

	// 邮件授权码
	AuthCode string

	// QQ 邮箱： SMTP 服务器地址：smtp.qq.com（SSL协议端口：465/587, 非SSL协议端口：25）
	// 163 邮箱：SMTP 服务器地址：smtp.163.com（SSL协议端口：465/994，非SSL协议端口：25）
	Host string
	Port int
	C    *gomail.Dialer
	Msg  *gomail.Message
}

func NewMail(host, from, auth string, port int) *Mail {
	m := &Mail{
		From:     from,
		AuthCode: auth,
		Host:     host,
		Port:     port,
	}
	m.C = gomail.NewDialer(
		m.Host,
		m.Port,
		m.From,
		m.AuthCode,
	)
	m.Msg = gomail.NewMessage()
	m.Msg.SetHeader("From", from)
	return m
}

func (m *Mail) Title(title string) *Mail {
	m.Msg.SetHeader("Subject", title)
	return m
}

func (m *Mail) HtmlBody(body string) *Mail {
	m.Msg.SetBody("text/html", body)
	return m
}

func (m *Mail) Send(to string) error {
	m.Msg.SetHeader("To", to)
	return m.C.DialAndSend(m.Msg)
}

func (m *Mail) SendMore(to []string) error {
	mgs := make([]*gomail.Message, 0)
	for _, v := range to {
		newMsg := m.Msg
		newMsg.SetHeader("To", v)
		mgs = append(mgs, newMsg)
	}
	return m.C.DialAndSend(mgs...)
}
