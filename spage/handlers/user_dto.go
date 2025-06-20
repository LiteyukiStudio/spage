package handlers

// RegisterReq 注册请求结构体
type RegisterReq struct {
	Username string `json:"username" binding:"required"` // 用户名
	Password string `json:"password" binding:"required"` // 密码
	Email    string `json:"email" binding:"required"`    // 邮箱

}

// LoginReq 登录请求结构体
type LoginReq struct {
	Username     string `json:"username" binding:"required"`      // 用户名
	Password     string `json:"password" binding:"required"`      // 密码
	CaptchaToken string `json:"captcha_token" binding:"required"` // 验证码
	Remember     bool   `json:"remember"`                         // 记住我
}

// CreateTokenReq 创建Token请求
type CreateTokenReq struct {
	Expire uint `json:"expire" binding:"required"`
}

// UserDTO 组织信息数据传输对象
type UserDTO struct {
	ID            uint              `json:"id"`            // 用户ID
	Name          string            `json:"name"`          // 用户名
	DisplayName   *string           `json:"display_name"`  // 显示名称
	Email         *string           `json:"email"`         // 邮箱
	Description   string            `json:"description"`   // 描述
	Avatar        *string           `json:"avatar_url"`    // 头像
	Role          string            `json:"role"`          // 角色
	Organizations []OrganizationDTO `json:"organizations"` // 组织
	Language      string            `json:"language"`      // 语言
	//Password      string            `json:"password"` // 密码
}
