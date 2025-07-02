package middle

import (
	"context"
	"github.com/LiteyukiStudio/spage/config"
	"github.com/LiteyukiStudio/spage/constants"
	"github.com/LiteyukiStudio/spage/utils"
	"github.com/cloudwego/hertz/pkg/app"
)

type verifyEmailReq struct {
	EmailVerifyKey string `json:"email_verify_key"` // 验证电子邮件的密钥
	VerifyCode     string `json:"verify_code"`
}

// NeedVerifyEmail 需要验证电子邮件的中间件
func (authType) NeedVerifyEmail() app.HandlerFunc {
	if config.EmailEnable {
		return func(ctx context.Context, c *app.RequestContext) {
			// TODO 电子邮件验证思路
			// 从请求中获取验证key和验证码，在kv中验证，验证通过后也不知道邮箱是什么
		}
	} else {
		return func(ctx context.Context, c *app.RequestContext) {
			c.Next(ctx)
		}
	}
}

func verifyEmailVerifyCode(key, code string) bool {
	kvStore := utils.GetKVStore()
	value, exists := kvStore.Get(constants.KVPrefixEmailVerifyCode + key)
	if !exists {
		return false
	}
	verifyCode, ok := value.(string)
	if !ok || verifyCode != code {
		return false
	}
	// 删除验证码，防止重复使用
	kvStore.Delete(constants.KVPrefixEmailVerifyCode + key)
	return true
}
