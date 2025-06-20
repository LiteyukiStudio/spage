package store

import (
	"github.com/LiteyukiStudio/spage/constants"
	"github.com/LiteyukiStudio/spage/spage/models"
)

type orgType struct {
}

var Org = orgType{}

// ListByUserID 通过UserID获取用户组织，支持分页和预加载关系
func (o *orgType) ListByUserID(userID string, page, limit int) (orgs []models.Organization, err error) {
	// 使用连接查询
	query := DB.Joins("JOIN organization_members ON organizations.id = organization_members.organization_id").
		Where("organization_members.user_id = ?", userID)
	// 预加载关系
	query = WithPreloads(query, constants.PreloadFieldMembers, constants.PreloadFieldOwners)
	// 使用通用分页方法
	orgs, _, err = Paginate[models.Organization](
		query,
		page,
		limit,
	)

	return
}

// GetOrgById 通过ID获取组织
func (o *orgType) GetOrgById(id uint) (org *models.Organization, err error) {
	err = DB.Model(&models.Organization{}).Where("id = ?", id).Preload(constants.PreloadFieldMembers).Preload(constants.PreloadFieldOwners).First(&org).Error
	return
}

// GetOrgByName 通过名称获取组织
func (o *orgType) GetOrgByName(name string) (org *models.Organization, err error) {
	err = DB.Model(&models.Organization{}).Where("name = ?", name).Preload(constants.PreloadFieldMembers).Preload(constants.PreloadFieldOwners).First(&org).Error
	return
}

// OrgNameIsExist 判断组织名称是否存在
func (o *orgType) OrgNameIsExist(name string) bool {
	var count int64
	DB.Model(&models.Organization{}).Where("name = ?", name).Count(&count)
	return count > 0
}

// GetUserRole 获取用户在组织中的权限
func (o *orgType) GetUserRole(org *models.Organization, userID uint) (auth string) {
	for _, owner := range org.Owners {
		if owner.ID == userID {
			return constants.OrgRoleOwner
		}
	}
	for _, member := range org.Members {
		if member.ID == userID {
			return constants.OrgRoleMember
		}
	}
	return ""
}

// CreateOrg 创建组织
func (o *orgType) CreateOrg(org *models.Organization) error {
	return DB.Create(org).Error
}

// UpdateOrg 更新组织
func (o *orgType) UpdateOrg(org *models.Organization) error {
	return DB.Updates(org).Error
}

// DeleteOrg 删除组织
func (o *orgType) DeleteOrg(org *models.Organization) error {
	return DB.Model(org).Delete(org).Error
}
