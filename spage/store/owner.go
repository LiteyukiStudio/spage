package store

import (
	"fmt"
	"github.com/LiteyukiStudio/spage/constants"
	"github.com/LiteyukiStudio/spage/spage/models"
	"github.com/LiteyukiStudio/spage/utils"
)

type ownerType struct{}

var Owner = &ownerType{}

// IsNameAvailable 检查所有者名称是被注册且合法
func (o *ownerType) IsNameAvailable(name string) bool {
	if !utils.IsValidEntityName(name) {
		return false
	}
	var userCount, orgCount int64
	if err := DB.Model(&models.User{}).Where("name = ?", name).Count(&userCount).Error; err != nil {
		return false
	}
	if err := DB.Model(&models.Organization{}).Where("name = ?", name).Count(&orgCount).Error; err != nil {
		return false
	}
	return (userCount + orgCount) == 0
}

// GetOwner 根据所有者类型和ID获取用户或组织，通过检查不为nil的返回值来判断所有者类型
func (o *ownerType) GetOwner(ownerType string, ownerID uint) (*models.User, *models.Organization, error) {
	switch ownerType {
	case constants.OwnerTypeUser:
		var user models.User
		if err := DB.First(&user, ownerID).Error; err != nil {
			return nil, nil, err
		}
		return &user, nil, nil
	case constants.OwnerTypeOrg:
		var org models.Organization
		if err := DB.First(&org, ownerID).Error; err != nil {
			return nil, nil, err
		}
		return nil, &org, nil
	default:
		return nil, nil, fmt.Errorf("不支持的所有者类型: %s", ownerType)
	}
}
