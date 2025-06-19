package utils

import (
	"fmt"
	"github.com/LiteyukiStudio/spage/constants"
	"github.com/go-resty/resty/v2"
)

type captchaType struct{}

var Captcha = captchaType{}

type CaptchaConfig struct {
	Type        string
	SiteSecrete string
	SecretKey   string
}

// VerifyCaptcha 根据提供的配置和令牌验证验证码
func (captchaType) VerifyCaptcha(restyClient *resty.Client, captchaConfig *CaptchaConfig, captchaToken string) (bool, error) {
	switch captchaConfig.Type {
	case constants.CaptchaTypeDisable:
		return true, nil
	case constants.CaptchaTypeHCaptcha:
		result := make(map[string]any)
		resp, err := restyClient.R().
			SetFormData(map[string]string{
				"secret":   captchaConfig.SecretKey,
				"response": captchaToken,
			}).SetResult(&result).Post("https://hcaptcha.com/siteverify")
		if err != nil {
			return false, err
		}
		if resp.IsError() {
			return false, nil
		}
		fmt.Printf("%#v\n", result)
		if success, ok := result["success"].(bool); ok && success {
			return true, nil
		} else {
			return false, nil
		}
	case constants.CaptchaTypeTurnstile:
		result := make(map[string]any)
		resp, err := restyClient.R().
			SetFormData(map[string]string{
				"secret":   captchaConfig.SecretKey,
				"response": captchaToken,
			}).SetResult(&result).Post("https://challenges.cloudflare.com/turnstile/v0/siteverify")
		if err != nil {
			return false, err
		}
		if resp.IsError() {
			return false, nil
		}
		fmt.Printf("%#v\n", result)
		if success, ok := result["success"].(bool); ok && success {
			return true, nil
		} else {
			return false, nil
		}
	case constants.CaptchaTypeReCaptcha:
		result := make(map[string]any)
		resp, err := restyClient.R().
			SetFormData(map[string]string{
				"secret":   captchaConfig.SecretKey,
				"response": captchaToken,
			}).SetResult(&result).Post("https://www.google.com/recaptcha/api/siteverify")
		if err != nil {
			return false, err
		}
		if resp.IsError() {
			return false, nil
		}
		fmt.Printf("%#v\n", result)
		if success, ok := result["success"].(bool); ok && success {
			return true, nil
		} else {
			return false, nil
		}
	default:
		return false, fmt.Errorf("invalid captcha type: %s", captchaConfig.Type)
	}
}
