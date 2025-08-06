package config

import (
	"embed"
	"errors"
	"fmt"
	"github.com/LiteyukiStudio/spage/pkg/constants"
	"github.com/LiteyukiStudio/spage/pkg/utils/filedriver"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	ServerPort string
	// Mode 运行模式 dev/prod
	Mode      = constants.ModeProd
	JwtSecret string
	LogLevel  = "info"
	// 日志级别 Log Level

	AdminUsername = "admin"
	AdminPassword = "admin"

	// BaseUrl 基础路径
	BaseUrl = "http://localhost:3000"
	OidcUri = "/api/v1/user/oidc/login"
	// EmailEnable Email 相关配置项 Email Configuration
	EmailEnable      bool   // 是否启用邮箱发送 Enable Email Sending
	EmailUsername    string // 邮箱用户名 Email Username
	EmailAddress     string // 邮箱地址 Email Address
	EmailHost        string // 邮箱服务器地址 Email Server Address
	EmailPort        string // 邮箱服务器端口 Email Server Port
	EmailPassword    string // 邮箱密码 Email Password
	EmailSSL         bool   // 是否启用SSL Enable SSL
	EmailTestAddress string // 测试邮箱地址，用于测试邮件发送功能 Test Email Address
	// DomainVerifyPolice 域名相关配置
	DomainVerifyPolice = constants.DomainVerifyPolicyLoose // 域名验证策略，默认为宽松验证 Loose Domain Verification Policy
	// PageLimit 每页显示的文章数量，默认为40
	PageLimit = 40
	// CaptchaType 验证码类型，支持turnstile、recaptcha和hcaptcha
	CaptchaType      = constants.CaptchaTypeDisable
	CaptchaSiteKey   string // reCAPTCHA v3的站点密钥
	CaptchaSecretKey string // reCAPTCHA v3的密钥
	CaptchaUrl       string // for mcaptcha

	TokenExpireTime        = 60 * 5
	RefreshTokenExpireTime = 3600 * 24

	BuildTime  = "0000-00-00 00:00:00" // 构建时间 Build Time
	Version    = "0.0.0"               // 版本号 Version
	CommitHash = "unknown"             // 提交哈希 Commit Hash

	AllowRegisterByOidc = true // 是否允许通过OIDC注册 Allow Register by OIDC
	AllowRegister       = true // 是否允许注册 Allow Register

	// Meta

	Icon = "/apage.svg"
	Name = "liteyuki spage"

	ReleaseSavePath  = "./data/releases"
	UploadsPath      = "./data/uploads"
	Title            = "spage"           // 网站标题 Website Title
	FileMaxSize      = 100 * 1024 * 1024 // 文件最大大小，单位字节 File Max Size
	FileDriverConfig = &filedriver.DriverConfig{
		Type:     constants.FileDriverLocal,
		BasePath: UploadsPath,
	}

	//go:embed config.example.yaml
	configExample embed.FS
)

// InitConfig 初始化配置文件
func InitConfig() error {
	configPath := "config.yaml"
	// 目标配置文件路径
	// 如果 config.yaml 已存在，直接返回
	if _, err := os.Stat(configPath); err == nil {
		return nil
	}
	// 读取嵌入的示例配置
	data, err := configExample.ReadFile("config.example.yaml")
	if err != nil {
		return errors.New("failed to read embedded config: " + err.Error())
	}
	// 确保目录存在（如果 config.yaml 不在当前目录）
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return errors.New("failed to create config directory: " + err.Error())
	}
	// 写入 config.yaml
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return errors.New("failed to write config file: " + err.Error())
	}
	return nil
}

// Init 初始化配置文件和常量
func Init() error {
	configPath := os.Getenv("CONFIG")
	if configPath != "" {
		configPath = filepath.Clean(configPath)
	} else {
		configPath = "./config.yaml"
	}
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			err := InitConfig()
			if err != nil {
				return err
			}
		}
	}
	// 初始化配置常量
	ServerPort = GetString("server.port", ServerPort)
	Mode = GetString("mode", Mode)
	LogLevel = GetString("log.level", LogLevel)
	BaseUrl = GetString("base-url", BaseUrl)

	// Admin配置项
	AdminUsername = GetString("admin.username", AdminUsername)
	AdminPassword = GetString("admin.password", AdminPassword)
	// 注册相关配置项
	AllowRegisterByOidc = GetBool("allow-register-by-oidc", AllowRegisterByOidc)
	AllowRegister = GetBool("allow-register", AllowRegister)

	// Captcha配置项
	CaptchaType = GetString("captcha.type", CaptchaType)
	CaptchaSiteKey = GetString("captcha.site-key", CaptchaSiteKey)
	CaptchaSecretKey = GetString("captcha.secret-key", CaptchaSecretKey)
	CaptchaUrl = GetString("captcha.url", CaptchaUrl)

	// Email配置项
	EmailEnable = GetBool("email.enable", EmailEnable)
	EmailUsername = GetString("email.username", EmailUsername)
	EmailAddress = GetString("email.address", EmailAddress)
	EmailHost = GetString("email.host", EmailHost)
	EmailPort = GetString("email.port", EmailPort)
	EmailPassword = GetString("email.password", EmailPassword)
	EmailSSL = GetBool("email.ssl", EmailSSL)
	EmailTestAddress = GetString("email.test-address", EmailTestAddress)

	// 域名配置项
	DomainVerifyPolice = GetString("domain.verify-policy", DomainVerifyPolice)

	// Meta配置项
	Icon = GetString("meta.icon", Icon)
	Name = GetString("meta.name", Name)

	// File存储配置项
	ReleaseSavePath = GetString("file.release-path", ReleaseSavePath)
	UploadsPath = GetString("file.uploads-path", UploadsPath)
	FileMaxSize = GetInt("file.max-size", FileMaxSize)

	FileDriverConfig = &filedriver.DriverConfig{
		Type:           GetString("file.driver.type", constants.FileDriverLocal),
		BasePath:       UploadsPath,
		WebDavUrl:      GetString("file.driver.webdav.url", ""),
		WebDavUser:     GetString("file.driver.webdav.user", ""),
		WebDavPassword: GetString("file.driver.webdav.password", ""),
		WebDavPolicy:   GetString("file.driver.webdav.policy", constants.WebDavPolicyProxy),
	}
	Title = GetString("title", Title)

	// 分页查询限制
	PageLimit = GetInt("page-limit", PageLimit)

	// Session过期时间
	TokenExpireTime = GetInt("token.expire", TokenExpireTime)
	RefreshTokenExpireTime = GetInt("token.refresh-expire", RefreshTokenExpireTime)
	JwtSecret = GetString("token.secret", JwtSecret)
	logrus.Info("Configuration loaded successfully, mode: ", Mode)

	// 设置日志级别
	logLevel, err := logrus.ParseLevel(LogLevel)
	if err != nil {
		logrus.Error("Invalid log level, using default level: info")
		logLevel = logrus.InfoLevel
	}
	logrus.SetLevel(logLevel)
	logrus.Info("Log level set to: ", logLevel)
	logrus.Debugln("LogLevel is: ", LogLevel)

	// 储存相关配置
	// 创建上传目录和release目录
	for _, path := range []string{UploadsPath, ReleaseSavePath} {
		if err := os.MkdirAll(path, 0755); err != nil {
			return errors.New("failed to create directory: " + path + ", error: " + err.Error())
		}
	}
	return nil

}

type agentConfig struct {
	Server struct {
		Token string `yaml:"token"`
		Host  string `yaml:"host"`
		Port  string `yaml:"port"`
	} `yaml:"server"` // 服务器配置
	Service struct {
		Host   string `yaml:"host"`
		Port   string `yaml:"port"`
		Static string `yaml:"static"`
	} `yaml:"service"` // 当前服务配置
	Caddy struct {
		Point string `yaml:"point"`
	} `yaml:"caddy"` // Caddy 配置
}

var AgentConfig = agentConfig{}

func (a *agentConfig) LoadByCmd() {
	Cmd.ParseArgs()
	a.Server.Token = Cmd.GetString("server.token", "")
	a.Server.Host = Cmd.GetString("server.host", "127.0.0.1")
	a.Server.Port = Cmd.GetString("server.port", "9526")
	a.Service.Host = Cmd.GetString("service.host", "")
	a.Service.Port = Cmd.GetString("service.port", "9527")
	a.Service.Static = Cmd.GetString("service.static", "./static")
	a.Caddy.Point = Cmd.GetString("caddy.point", "")
}

// InitAgentConfig 初始化 Agent 配置文件，如果不存在则创建；存在则读取并赋值给 a
func (a *agentConfig) InitAgentConfig() error {
	// 获取配置路径
	configPath := os.Getenv("CONFIG")
	if configPath == "" {
		configPath = "./config.yaml"
	} else {
		configPath = filepath.Clean(configPath)
	}

	// 设置 viper 的读取目标
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	// 尝试读取配置
	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			// 文件不存在：序列化当前结构体 a 并写入配置文件
			configData, err := yaml.Marshal(a)
			if err != nil {
				return fmt.Errorf("序列化配置失败: %w", err)
			}
			// 创建文件夹（若存在上级路径）
			if err := os.MkdirAll(filepath.Dir(configPath), os.ModePerm); err != nil {
				return fmt.Errorf("创建配置文件目录失败: %w", err)
			}
			if err := os.WriteFile(configPath, configData, 0644); err != nil {
				return fmt.Errorf("写入配置文件失败: %w", err)
			}
			return nil
		}
		// 如果是其他读取错误
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 如果成功读取配置，则将其绑定到结构体 a
	if err := viper.Unmarshal(a); err != nil {
		return fmt.Errorf("解析配置文件失败: %w", err)
	}

	return nil
}

// Get 返回配置项的值，如果不存在则返回默认值
func Get[T any](key string, defaultValue T) T {
	if !viper.IsSet(key) {
		return defaultValue
	}

	value := viper.Get(key)
	if v, ok := value.(T); ok {
		return v
	}
	return defaultValue
}

// GetString 返回配置项的字符串值
func GetString(key string, defaultValue ...string) string {
	if len(defaultValue) > 0 {
		return Get(key, defaultValue[0])
	}
	return viper.GetString(key)
}

// GetInt 返回配置项的整数值
func GetInt(key string, defaultValue ...int) int {
	if len(defaultValue) > 0 {
		return Get(key, defaultValue[0])
	}
	return viper.GetInt(key)
}

// GetBool 返回配置项的布尔值
func GetBool(key string, defaultValue ...bool) bool {
	if len(defaultValue) > 0 {
		return Get(key, defaultValue[0])
	}
	return viper.GetBool(key)
}

// GetFloat64 返回配置项的浮点数值
func GetFloat64(key string, defaultValue ...float64) float64 {
	if len(defaultValue) > 0 {
		return Get(key, defaultValue[0])
	}
	return viper.GetFloat64(key)
}

// GetStringSlice 返回配置项的字符串切片值
func GetStringSlice(key string, defaultValue ...[]string) []string {
	if len(defaultValue) > 0 {
		return Get(key, defaultValue[0])
	}
	return viper.GetStringSlice(key)
}
