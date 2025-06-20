package utils

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/smtp"
	"strings"
)

type EmailType struct{}

var Email = EmailType{}

type EmailConfig struct {
	Enable   bool   // 邮箱启用状态
	Username string // 邮箱用户名
	Address  string // 邮箱地址
	Host     string // 邮箱服务器地址
	Port     string // 邮箱服务器端口
	Password string // 邮箱密码
	SSL      bool   // 是否使用SSL
}

// SendTemplate 发送HTML模板，从配置文件中读取邮箱配置
func SendTemplate(emailConfig *EmailConfig, target, htmlTemplate string, placeholders map[string]string) error {
	for placeholder, value := range placeholders {
		htmlTemplate = strings.ReplaceAll(htmlTemplate, placeholder, value)
	}
	err := SendEmail(emailConfig, target, htmlTemplate, true)
	if err != nil {
		return err
	}
	return nil
}

// SendEmail 发送邮件
func SendEmail(emailConfig *EmailConfig, target, content string, isHTML bool) error {
	// 如果配置未启用，则直接返回nil
	if !emailConfig.Enable {
		return nil
	}

	auth := smtp.PlainAuth("", emailConfig.Username, emailConfig.Password, emailConfig.Host)

	var conn net.Conn
	var err error

	if emailConfig.SSL {
		conn, err = tls.Dial("tcp", net.JoinHostPort(emailConfig.Host, emailConfig.Port), &tls.Config{ServerName: emailConfig.Host})
	} else {
		conn, err = net.Dial("tcp", net.JoinHostPort(emailConfig.Host, emailConfig.Port))
	}

	if err != nil {
		return err
	}

	// 在函数退出时关闭连接
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			// todo: 处理连接关闭时的错误
		}
	}(conn)

	client, err := smtp.NewClient(conn, emailConfig.Host)
	if err != nil {
		return err
		// todo: 处理连接时的错误
	}

	defer func(client *smtp.Client) {
		err := client.Quit()
		if err != nil {
			// todo: 处理关闭连接时的错误
		}
	}(client)

	if err = client.Auth(auth); err != nil {
		return err
		// todo: 处理身份验证时的错误
	}

	if err = client.Mail(emailConfig.Address); err != nil {
		return err
		// todo: 处理发件人时的错误
	}

	if err = client.Rcpt(target); err != nil {
		return err
		// todo: 处理收件人时的错误
	}

	writer, err := client.Data()
	if err != nil {
		return err
		// todo: 处理数据写入器创建时的错误
	}

	defer func(writer io.WriteCloser) {
		err := writer.Close()
		if err != nil {
			// todo: 处理关闭写入器时的错误
		}
	}(writer)
	var message string
	if isHTML {
		message = fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: Notification\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s", emailConfig.Address, target, content)
	} else {
		message = fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: Notification\r\n\r\n%s", emailConfig.Address, target, content)
	}
	_, err = writer.Write([]byte(message))
	return err
}
