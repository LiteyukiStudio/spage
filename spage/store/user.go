package store

import (
	"errors"
	"github.com/LiteyukiStudio/spage/constants"
	"github.com/LiteyukiStudio/spage/spage/models"
	"github.com/sirupsen/logrus"

	"gorm.io/gorm"
)

type userType struct {
}

var User = userType{}

// Create 创建用户
func (u *userType) Create(user *models.User) (err error) {
	return DB.Create(user).Error
}

// GetByName 根据名称获取用户
func (u *userType) GetByName(name string) (user *models.User, err error) {
	user = &models.User{} // 初始化指针 // Initialize pointer
	err = DB.Where("name = ?", name).First(user).Error
	if err != nil {
		return nil, err // 出错时返回nil When an error occurs, return nil
	}
	return user, nil
}

// IsNameExist 判断用户名是否存在
func (u *userType) IsNameExist(name string) bool {
	var count int64
	err := DB.Model(&models.User{}).Where("name = ?", name).Count(&count).Error
	if err != nil {
		return false
	}
	return count > 0
}

// GetByID 根据ID获取用户
func (u *userType) GetByID(id uint) (user *models.User, err error) {
	user = &models.User{} // 初始化指针
	err = DB.Where("id = ?", id).First(user).Error
	if err != nil {
		return nil, err // 出错时返回nil
	}
	return user, nil
}

// GetByEmail 根据邮箱获取用户
func (u *userType) GetByEmail(email string) (user *models.User, err error) {
	user = &models.User{} // 初始化指针
	err = DB.Where("email = ?", email).First(user).Error
	if err != nil {
		return nil, err // 出错时返回nil
	}
	return user, nil
}

// FindOrCreateByEmail 根据邮箱查找或创建用户，仅在用户不存在时name才会生效
func (u *userType) FindOrCreateByEmail(email, name string) (*models.User, error) {
	user := &models.User{}
	// 尝试查找用户 Try to find user
	err := DB.Where("email = ?", email).First(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 用户不存在，创建新用户 User does not exist, create a new user
			user.Email = &email
			user.Name = name
			err = DB.Create(user).Error
			if err != nil {
				return nil, err // 创建失败时返回错误 Return error if creation fails
			}
		} else {
			return nil, err // 查找失败时返回错误 Return error if search fails
		}
	}
	return user, nil // 返回找到或创建的用户 Return found or created user
}

// Update 更新用户信息
func (u *userType) Update(user *models.User) error {
	logrus.Println("Updating user:", user.ID, user.Name, user.Email, user.Role)
	return DB.Updates(user).Error
}

// DeleteByID 根据ID删除用户
func (u *userType) DeleteByID(id uint) (err error) {
	err = DB.Delete(&models.User{}, id).Error
	if err != nil {
		return err
	}
	return nil
}

// UpdateSystemAdmin 更新系统管理员用户，不存在则创建
func (u *userType) UpdateSystemAdmin(user *models.User) (err error) {
	// 设置该用户为系统管理员 Set this user as a system admin
	user.Flag = constants.FlagSystemAdmin
	user.Role = constants.GlobalRoleAdmin

	// 尝试查找系统管理员 Try to find system admin
	existingAdmin := models.User{}
	result := DB.Where("flag = ?", constants.FlagSystemAdmin).First(&existingAdmin)

	if result.Error != nil {
		// 如果不存在系统管理员（记录未找到），则创建一个 If there is no system admin (record not found), create one
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// 创建新的系统管理员 Create new system admin
			return DB.Create(user).Error
		}
		// 其他错误则直接返回 Other errors are returned directly
		return result.Error
	}
	// 系统管理员已存在，更新信息 System admin exists, update information
	// 保留ID，更新其他字段 Keep ID, update other fields
	user.ID = existingAdmin.ID
	return DB.Model(&existingAdmin).Updates(user).Error
}

// IsMemberOfOrg 检查用户是否是组织成员
func (u *userType) IsMemberOfOrg(userID, orgID uint) (isMember bool, err error) {
	var count int64
	err = DB.Model(&models.Organization{}).
		Joins("JOIN organization_members ON organization_members.organization_id = organizations.id").
		Where("organization_members.user_id = ? AND organizations.id = ?", userID, orgID).
		Count(&count).Error

	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// IsOwnerOfOrg 检查用户是否是组织所有者
func (u *userType) IsOwnerOfOrg(userID, orgID uint) (isOwner bool, err error) {
	var count int64
	err = DB.Model(&models.Organization{}).
		Joins("JOIN organization_owners ON organization_owners.organization_id = organizations.id").
		Where("organization_owners.user_id = ? AND organizations.id = ?", userID, orgID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
