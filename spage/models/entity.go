package models

import (
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	gorm.Model
	Name          string          `gorm:"not null;uniqueIndex"`            // 用户的唯一名称 User's unique name
	DisplayName   *string         `gorm:"column:display_name"`             // 用户的显示名称 User's display name
	Email         *string         `gorm:"unique"`                          // 用户的电子邮件地址，只有用户的电子邮件地址是唯一的（用于 oidc 身份验证） User's email address, only the user's email address is unique (used for oidc authentication)
	Description   string          `gorm:"default:'No description.'"`       // 用户描述 User description
	AvatarURL     *string         `gorm:"column:avatar_url"`               // 留空以使用 Gravatar Leave blank to use Gravatar
	Role          string          `gorm:"not null;default:user"`           // 用户的全局角色 User's global role
	Organizations []*Organization `gorm:"many2many:organization_members;"` // 隶属于许多组织 Many organizations the user belongs to
	ProjectLimit  int             `gorm:"default:-1"`                      // 用户的项目限制，0 表示无限制 User's project limit, 0 means no limit
	Language      string          `gorm:"default:'zh-cn'"`                 // 用户的语言，默认为英语 User's language, default to English
	Flag          string          `gorm:"default:'0'"`                     // system_admin 的另一面旗帜 The other side of system_admin flag
	IsPrivate     bool            `gorm:"default:false"`                   // 用户是否为私有用户，默认为 false，表示公开用户 Whether the user is a private user, default is false, meaning public user
	Password      *string         `gorm:"column:password"`                 // 用户的密码（经过哈希处理），仅用于本地身份验证 User's password (hashed), only used for local authentication
}

func (User) TableName() string {
	return "users"
}

// Organization 组织模型
type Organization struct {
	gorm.Model
	Name         string  `gorm:"not null;uniqueIndex"`            // 组织的唯一名称 Organization's unique name
	DisplayName  *string `gorm:"column:display_name"`             // 组织的显示名称 Organization's display name
	Email        *string `gorm:"column:email"`                    // 组织的电子邮件地址 Organization's email address
	Description  string  `gorm:"default:'No description.'"`       // 组织描述 Organization description
	AvatarURL    *string `gorm:"column:avatar_url"`               // 留空以使用 Gravatar Leave blank to use Gravatar
	Members      []*User `gorm:"many2many:organization_members;"` // 组织的成员包含创建者 (including the creator)
	Owners       []User  `gorm:"many2many:organization_owners;"`  // 组织的所有者（无反向关系）包含创建者 (including the creator)
	IsPrivate    bool    `gorm:"default:false"`                   // 组织是否为私有组织，默认为 false，表示公开组织 Whether the organization is a private organization, default is false, meaning public organization
	ProjectLimit int     `gorm:"default:0"`                       // 组织的项目限制，0：遵循策略，-1：无限制 Organization's project limit, 0: follow the policy, -1: unlimited
}

func (Organization) TableName() string {
	return "organizations"
}

// Project 项目模型
type Project struct {
	gorm.Model
	Name        string  `gorm:"not null"`                   // 项目在一个主体下的唯一名称 Project's unique name
	DisplayName *string `gorm:"column:display_name"`        // 项目的显示名称 Project's display name
	Description string  `gorm:"default:'No description.'"`  // 项目描述 Project description
	OwnerID     uint    `gorm:"not null"`                   // 所有者 ID（用户 ID 或组织 ID） Owner ID (user ID or organization ID)
	OwnerType   string  `gorm:"not null"`                   // 所有者类型，可以是用户或组织 Owner type, can be user or organization
	Owners      []User  `gorm:"many2many:project_owners;"`  // 项目的所有者，无反向关系 Project's owners, no reverse relation
	Members     []*User `gorm:"many2many:project_members;"` // 项目的成员 Project's members
	SiteLimit   int     `gorm:"default:0"`                  // 项目的站点限制，0：遵循策略，-1：无限制 Project's site limit, 0: follow the policy, -1: unlimited
	IsPrivate   bool    `gorm:"default:false"`
}

func (Project) TableName() string {
	return "projects"
}
