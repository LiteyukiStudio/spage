package config

import (
	"embed"
	"errors"
	"github.com/LiteyukiStudio/spage/constants"
	"github.com/LiteyukiStudio/spage/utils/filedriver"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	ServerPort string
	// 服务器端口 Server Port

	Mode = constants.ModeProd
	// 运行模式，支持dev和prod

	JwtSecret string
	// JWT密钥 JWT Secret

	FrontEndURL string
	// 前端URL Frontend URL

	LogLevel = "info"
	// 日志级别 Log Level

	AdminUsername = "admin"
	// 管理员用户名 Admin Username

	AdminPassword = "admin"
	// 管理员密码 Admin Password

	EmailEnable   bool   // 是否启用邮箱发送 Enable Email Sending
	EmailUsername string // 邮箱用户名 Email Username
	EmailAddress  string // 邮箱地址 Email Address
	EmailHost     string // 邮箱服务器地址 Email Server Address
	EmailPort     string // 邮箱服务器端口 Email Server Port
	EmailPassword string // 邮箱密码 Email Password
	EmailSSL      bool   // 是否启用SSL Enable SSL

	PageLimit = 40
	// 每页显示的文章数量，默认为40

	CaptchaType = constants.CaptchaTypeDisable
	// 验证码类型，支持turnstile、recaptcha和hcaptcha

	CaptchaSiteKey   string // reCAPTCHA v3的站点密钥
	CaptchaSecretKey string // reCAPTCHA v3的密钥
	CaptchaUrl       string // for mcaptcha

	TokenExpireTime = 3600 * 24
	// session过期时间，单位秒

	RefreshTokenExpireTime = 3600 * 144
	// 刷新token过期时间，单位秒

	BuildTime  = "0000-00-00 00:00:00" // 构建时间 Build Time
	Version    = "0.0.0"               // 版本号 Version
	CommitHash = "unknown"             // 提交哈希 Commit MD5

	ReleaseSavePath  = "./data/releases"
	UploadsPath      = "./data/uploads"
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
	// 设置配置文件路径
	// Set the configuration file path
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	// 读取配置文件
	// Read the configuration file
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

	// Admin配置项
	AdminUsername = GetString("admin.username", AdminUsername)
	AdminPassword = GetString("admin.password", AdminPassword)

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

	// File存储配置项
	ReleaseSavePath = GetString("file.release-path", ReleaseSavePath)
	UploadsPath = GetString("file.uploads-path", UploadsPath)
	FileMaxSize = GetInt("file.max-size", FileMaxSize)

	FileDriverConfig = &filedriver.DriverConfig{
		Type:           GetString("file.driver.type", constants.FileDriverLocal),
		BasePath:       UploadsPath,
		WebDavUrl:      GetString("file.driver.webdav.url", UploadsPath),
		WebDavUser:     GetString("file.driver.webdav.user", ""),
		WebDavPassword: GetString("file.driver.webdav.password", ""),
	}

	// 分页查询限制
	PageLimit = GetInt("page-limit", PageLimit)

	// Session过期时间
	TokenExpireTime = GetInt("token.expire", TokenExpireTime)
	RefreshTokenExpireTime = GetInt("token.refresh-expire", RefreshTokenExpireTime)
	JwtSecret = GetString("token.secret", JwtSecret)

	// 从启动参数拿取一些配置项mode frontend-url
	argsMap := Cmd.GetArgsMap(os.Args[1:])
	queryKeys := []string{"mode", "frontend-url", "port"}
	for _, key := range queryKeys {
		if value, ok := argsMap[key]; ok {
			switch key {
			case "mode":
				Mode = value
			case "port":
				ServerPort = value
			case "frontend-url":
				FrontEndURL = value
			}
		}
	}
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
	// 其他配置项合法校验流程
	// ...
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
