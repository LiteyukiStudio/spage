package store

import (
	"fmt"
	"github.com/LiteyukiStudio/spage/constants"
	"github.com/LiteyukiStudio/spage/spage/models"
)

type projectType struct {
}

var Project = projectType{}

// Create 创建项目
func (p *projectType) Create(project *models.Project) (err error) {
	return DB.Create(project).Error
}

// GetByID 通过项目ID获取项目
func (p *projectType) GetByID(id uint) (project *models.Project, err error) {
	err = DB.First(&project, id).Preload(constants.PreloadFieldOwners).Preload(constants.PreloadFieldMembers).Error
	return
}

// UserIsOwner 判断用户是否是项目的所有者
func (p *projectType) UserIsOwner(project *models.Project, userID uint) bool {
	if project.OwnerType == constants.OwnerTypeUser && project.OwnerID == userID {
		return true
	}
	for _, owner := range project.Owners {
		if owner.ID == userID {
			return true
		}
	}
	return false
}

// ListByOwner 通过用户ID获取项目列表，支持分页和从新到旧排序
func (p *projectType) ListByOwner(ownerType, ownerID string, page, limit int) (projects []models.Project, total int64, err error) {
	if ownerType != constants.OwnerTypeUser && ownerType != constants.OwnerTypeOrg {
		return nil, 0, fmt.Errorf("owner type not allowed")
	}
	projects, total, err = Paginate[models.Project](
		DB,
		page,
		limit,
		"owner_type = ? AND owner_id = ?",
		ownerType,
		ownerID,
	)
	return
}

// Update 更新项目
func (p *projectType) Update(project *models.Project) (err error) {
	return DB.Updates(project).Error
}

// Delete 删除项目
func (p *projectType) Delete(project *models.Project) (err error) {
	return DB.Delete(project).Error
}

// AddOwner 为项目添加所有者
func (p *projectType) AddOwner(project *models.Project, user *models.User) (err error) {
	return DB.Model(project).Association(constants.PreloadFieldMembers).Append(user)
}

// DeleteOwner 从项目删除所有者
func (p *projectType) DeleteOwner(project *models.Project, user *models.User) (err error) {
	return DB.Model(project).Association(constants.PreloadFieldOwners).Delete(user)
}

// GetSiteList 获取项目下的站点列表
func (p *projectType) GetSiteList(project *models.Project, page, limit int) (sites []models.Site, total int64, err error) {
	sites, total, err = Paginate[models.Site](
		DB,
		page,
		limit,
		"project_id = ?",
		project.ID,
	)
	return
}
