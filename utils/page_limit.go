package utils

import (
	"github.com/LiteyukiStudio/spage/config"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
)

type ctxType struct{}

var Ctx = ctxType{}

// GetPageLimit 封装从上下文中的query获取查询参数并转换为整数的函数，出错返回默认值
func (ctxType) GetPageLimit(c *app.RequestContext) (page, limit int) {
	pageString := c.Query("page")
	limitString := c.Query("limit")
	page, err := strconv.Atoi(pageString)
	if err != nil || page <= 0 {
		page = 1 // 默认第一页
	}

	limit, err = strconv.Atoi(limitString)
	if err != nil || limit <= 0 || limit > config.PageLimit {
		limit = config.PageLimit
	}
	return
}
