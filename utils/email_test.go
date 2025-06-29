package utils

import (
	"context"
	"github.com/LiteyukiStudio/spage/config"
	"os"
	"testing"
)

// test
func TestSendTemplate(t *testing.T) {
	os.Setenv("CONFIG", "../config.yaml")
	config.Init()

	// 使用简单的模板
	tmpl := `<html><body><h1>Hello, {{.Name}}!</h1><p>Welcome to our service.</p></body></html>`

	emailConfig := &EmailConfig{
		Enable:   true,
		Username: config.EmailUsername,
		Address:  config.EmailAddress,
		Host:     config.EmailHost,
		Port:     config.EmailPort,
		Password: config.EmailPassword,
		SSL:      config.EmailSSL,
	}

	t.Logf("Email configuration: %+v", emailConfig)

	// 使用正确的模板数据格式
	data := map[string]any{
		"Name": "Liteyuki",
	}
	// 创建超时上下文
	ctx := context.Background()
	// 发送测试邮件
	err := SendTemplate(ctx, emailConfig, config.EmailTestAddress, "Test Email", tmpl, data)
	if err != nil {
		t.Errorf("Failed to send email: %v", err)
	} else {
		t.Log("Email sent successfully")
	}
}
