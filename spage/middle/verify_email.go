package middle

import (
	"context"
	"github.com/LiteyukiStudio/spage/config"
	"github.com/LiteyukiStudio/spage/constants"
	"github.com/LiteyukiStudio/spage/resps"
	"github.com/LiteyukiStudio/spage/utils"
	"github.com/cloudwego/hertz/pkg/app"
)

// NeedVerifyEmail 需要验证电子邮件的中间件
func (authType) NeedVerifyEmail() app.HandlerFunc {
	if config.EmailEnable {
		return func(ctx context.Context, c *app.RequestContext) {
			email, ok1 := utils.GetJsonFieldFromCtx(c, "email")
			emailVerifyCode, ok2 := utils.GetJsonFieldFromCtx(c, "email_verify_code")
			if !ok1 || !ok2 || email == "" || emailVerifyCode == "" {
				resps.BadRequest(c, "缺少电子邮件或验证码")
				return
			}
			if !verifyEmailVerifyCode(email, emailVerifyCode) {
				resps.BadRequest(c, "电子邮件验证码错误或已过期")
				return
			}
			c.Next(ctx)
		}
	} else {
		return func(ctx context.Context, c *app.RequestContext) {
			c.Next(ctx)
		}
	}
}

func verifyEmailVerifyCode(email, code string) bool {
	kvStore := utils.GetKVStore()
	value, exists := kvStore.Get(constants.KVPrefixEmailVerifyCode + email)
	if !exists {
		return false
	}
	verifyCode, ok := value.(string)
	if !ok || verifyCode != code {
		return false
	}
	kvStore.Delete(constants.KVPrefixEmailVerifyCode + email)
	return true
}
