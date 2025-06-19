package handlers

// OrganizationDTO 组织信息数据传输对象
type OrganizationDTO struct {
	ID           uint      `json:"id"`            // 组织ID
	Name         string    `json:"name"`          // 组织名称
	DisplayName  *string   `json:"display_name"`  // 显示名称
	Email        *string   `json:"email"`         // 邮箱地址
	Description  string    `json:"description"`   // 描述信息
	AvatarURL    *string   `json:"avatar_url"`    // 头像URL
	ProjectLimit int       `json:"project_limit"` // 项目数量限制
	Members      []UserDTO `json:"members"`       // 组织成员
	Owners       []UserDTO `json:"owners"`        // 组织所有者
}

type CreateOrgReq struct {
	Name        string  `json:"name" binding:"required"`         // 组织名称
	DisplayName string  `json:"display_name" binding:"required"` // 显示名称
	Email       *string `json:"email"`                           // 邮箱地址
	Description string  `json:"description" binding:"required"`  // 描述信息
	AvatarURL   *string `json:"avatar_url"`                      // 头像URL
}

// UpdateOrgReq 用于更新组织信息的请求体
type UpdateOrgReq struct {
	Name        string  `json:"name"`         // 组织名称
	DisplayName *string `json:"display_name"` // 显示名称
	Email       *string `json:"email"`        // 邮箱地址
	Description *string `json:"description"`  // 描述信息
	AvatarURL   *string `json:"avatar_url"`   // 头像URL
}

// GetOrgProjectReq 用于获取组织项目的请求体
type GetOrgProjectReq struct {
	Page    int    `json:"page" binding:"required"`  // 页码
	Limit   int    `json:"limit" binding:"required"` // 每页项目数量
	OrderBy string // 排序字段
}

// OrgUserReq 用于添加或删除组织用户的请求体
type OrgUserReq struct {
	UserID uint   `json:"user_id" binding:"required"` // 用户ID
	Role   string `json:"role" binding:"required"`    // 角色
}
