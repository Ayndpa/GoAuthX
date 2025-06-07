package email

import (
	"fmt"
	"goauthx/internal/config"
	"gopkg.in/gomail.v2"
)

// SendEmail 发送邮件
func SendEmail(to []string, subject, body string) error {
	cfg := config.GetConfig()
	m := gomail.NewMessage()

	from := fmt.Sprintf("%s <%s>", cfg.Name, cfg.SMTP.Username)
	m.SetHeader("From", from)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(cfg.SMTP.Host, cfg.SMTP.Port, cfg.SMTP.Username, cfg.SMTP.Password)

	if err := d.DialAndSend(m); err != nil {
		fmt.Printf("Gomail send error: %v\n", err)
		return err
	}
	return nil
}
