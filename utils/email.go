package utils

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"gopkg.in/gomail.v2"
	"html/template"
	"strconv"
)

type emailType struct{}

var Email = emailType{}

type EmailConfig struct {
	Enable   bool   // 邮箱启用状态
	Username string // 邮箱用户名
	Address  string // 邮箱地址
	Host     string // 邮箱服务器地址
	Port     string // 邮箱服务器端口
	Password string // 邮箱密码
	SSL      bool   // 是否使用SSL
}

// SendTemplate 发送HTML模板，从配置文件中读取邮箱配置，支持上下文控制
func SendTemplate(ctx context.Context, emailConfig *EmailConfig, target, subject, htmlTemplate string, data map[string]interface{}) error {
	// 使用Go的模板系统处理HTML模板
	tmpl, err := template.New("email").Parse(htmlTemplate)
	if err != nil {
		return fmt.Errorf("解析模板失败: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("执行模板失败: %w", err)
	}

	// 发送处理后的HTML内容
	return SendEmail(ctx, emailConfig, target, subject, buf.String(), true)
}

// SendEmail 使用gomail库发送邮件
func SendEmail(ctx context.Context, emailConfig *EmailConfig, target, subject, content string, isHTML bool) error {
	if !emailConfig.Enable {
		return nil
	}

	// 创建新邮件
	m := gomail.NewMessage()
	m.SetHeader("From", emailConfig.Address)
	m.SetHeader("To", target)
	m.SetHeader("Subject", subject)

	// 设置内容类型
	if isHTML {
		m.SetBody("text/html", content)
	} else {
		m.SetBody("text/plain", content)
	}

	// 转换端口号为整数
	port, err := strconv.Atoi(emailConfig.Port)
	if err != nil {
		return fmt.Errorf("端口号格式不正确: %w", err)
	}

	// 创建发送器
	d := gomail.NewDialer(emailConfig.Host, port, emailConfig.Username, emailConfig.Password)

	// 配置SSL/TLS
	if emailConfig.SSL {
		d.SSL = true
	} else {
		// 对于非SSL但需要STARTTLS的情况
		d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}

	// 发送邮件
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("发送邮件失败: %w", err)
	}

	return nil
}
