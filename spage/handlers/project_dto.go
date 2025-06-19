package handlers

// ProjectDTO 项目信息数据传输对象
type ProjectDTO struct {
	ID          uint      `json:"id"`           // 项目ID
	Name        string    `json:"name"`         // 项目名称
	DisplayName *string   `json:"display_name"` // 项目显示名称
	Description string    `json:"description"`  // 项目描述
	OwnerType   string    `json:"owner_type"`   // 项目拥有者类型
	OwnerID     uint      `json:"owner_id"`     // 项目拥有者ID
	Owners      []UserDTO `json:"owners"`       // 项目拥有者列表
	SiteLimit   int       `json:"site_limit"`   // 项目站点数量限制
}

// CreateProjectReq 创建项目请求参数
type CreateProjectReq struct {
	Name        string  `json:"name" binding:"required"`                                         // 项目名称
	DisplayName *string `json:"display_name"`                                                    // 项目显示名称
	Description string  `json:"description"`                                                     // 项目描述
	OwnerType   string  `json:"owner_type" binding:"required"  vd:"in($,'user','organization')"` // 项目拥有者类型
	OwnerID     uint    `json:"owner_id" binding:"required"`                                     // 项目拥有者ID

}

// UpdateProjectReq 更新项目请求参数
type UpdateProjectReq struct {
	Name        *string `json:"name"`         // 项目名称
	DisplayName *string `json:"display_name"` // 项目显示名称
	Description *string `json:"description"`  // 项目描述
}

// ProjectUserReq 项目用户请求参数
type ProjectUserReq struct {
	UserID uint `json:"user_id" binding:"required"` // 用户ID
}

// GetProjectListReq 获取项目列表请求参数
type GetSiteListReq struct {
	Page     int    `form:"page" binding:"required"`    // 页码
	Limit    int    `form:"limit" binding:"required"`   // 每页数量
	Project  string `form:"project" binding:"required"` // 项目名称
	SiteName string `form:"site_name"`                  // 站点名称
}
