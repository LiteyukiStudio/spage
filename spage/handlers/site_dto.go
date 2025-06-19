package handlers

// SiteDTO 网站详情
type SiteDTO struct {
	ID          uint       `json:"id"`          // 网站ID
	Name        string     `json:"name"`        // 网站名称
	Description string     `json:"description"` // 网站描述
	ProjectID   uint       `json:"project_id"`  // 项目ID
	Project     ProjectDTO `json:"project"`     // 项目详情
	SubDomain   *string    `json:"sub_domain"`  // 子域名
	Domains     []string   `json:"domains"`     // 域名
}

// CreateSiteReq 创建网站请求参数
type CreateSiteReq struct {
	Name        string   `json:"name" binding:"required"`       // 网站名称
	Description string   `json:"description"`                   // 网站描述
	ProjectID   uint     `json:"project_id" binding:"required"` // 项目ID
	SubDomain   *string  `json:"sub_domain"`                    // 子域名
	Domains     []string `json:"domains"`                       // 域名
}

type UpdateSiteReq struct {
	Name        *string  `json:"name"`        // 网站名称
	Description *string  `json:"description"` // 网站描述
	SubDomain   *string  `json:"sub_domain"`  // 子域名
	Domains     []string `json:"domains"`     // 域名
}
