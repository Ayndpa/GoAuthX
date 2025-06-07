package smtp

import (
	"bytes"
	gomail "gopkg.in/gomail.v2"
	"hub/internal/config"
	"log"
	"os"
	"path/filepath"
)

// SendVerificationCode 发送验证码邮件
func SendVerificationCode(toEmail, code string) error {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Printf("获取配置出错: %v", err)
		return err
	}

	// 读取邮件模板内容
	tmplPath := filepath.Join("static", "email.html")
	tmplBytes, err := os.ReadFile(tmplPath)
	if err != nil {
		log.Printf("读取邮件模板出错: %v", err)
		return err
	}
	bodyStr := string(tmplBytes)
	// 用简单的字符串替换
	bodyStr = string(bytes.ReplaceAll([]byte(bodyStr), []byte("{{CODE}}"), []byte(code)))

	// 使用gomail构造邮件
	m := gomail.NewMessage()
	m.SetHeader("From", cfg.SMTP.From+"<"+cfg.SMTP.Username+">")
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "验证码通知")
	m.SetBody("text/html", bodyStr)

	port := cfg.SMTP.Port
	host := cfg.SMTP.Host

	d := gomail.NewDialer(host, port, cfg.SMTP.Username, cfg.SMTP.Password)
	// gomail默认支持TLS

	if err := d.DialAndSend(m); err != nil {
		log.Printf("发送邮件出错: %v", err)
		return err
	}
	return nil
}
