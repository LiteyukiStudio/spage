package utils

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/tidwall/gjson"
)

// GetJsonFieldFromCtx  高效获取JSON字段
func GetJsonFieldFromCtx(c *app.RequestContext, path string) (string, bool) {
	body := c.Request.Body()
	if len(body) == 0 {
		return "", false
	}
	result := gjson.GetBytes(body, path)
	if result.Exists() {
		return result.String(), true
	}
	return "", false
}
