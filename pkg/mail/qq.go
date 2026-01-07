package mail

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"os"
	"strings"

	"github.com/jordan-wright/email"
)

type qqMail struct {
	smtpHost string
	smtpPort string
	addr     string
	username string
	password string
	from     string
}

func NewQQMail() (Manager, error) {
	username := os.Getenv("QQ_USERNAME")
	password := os.Getenv("QQ_PASSWORD")
	if username == "" || password == "" {
		return nil, fmt.Errorf("QQ_USERNAME or QQ_PASSWORD is empty")
	}
	return &qqMail{
		smtpHost: "smtp.qq.com",
		smtpPort: "465",
		addr:     "smtp.qq.com:465",
		username: username,
		from:     "web-chat <" + username + ">",
		password: password,
	}, nil
}
func (q *qqMail) Send(msg Message, targets []string) error {
	e := &email.Email{
		From:    q.from,
		To:      targets,
		Subject: msg.Title,
		Text:    []byte(msg.Content.(string)),
	}
	if msg.Appendix != "" {
		if _, err := e.AttachFile(msg.Appendix); err != nil {
			return err
		}
	}

	auth := smtp.PlainAuth("", q.username, q.password, q.smtpHost)

	tlsCfg := &tls.Config{
		ServerName: q.smtpHost,
		MinVersion: tls.VersionTLS12,
	}

	err := e.SendWithTLS(q.addr, auth, tlsCfg)
	if err != nil {
		s := err.Error()
		if strings.Contains(s, "short response") {
			return nil
		}
		return err
	}
	return nil

}
