package handlers

import (
	"context"
	"github.com/LiteyukiStudio/spage/constants"
	"github.com/LiteyukiStudio/spage/spage/middle"
	"github.com/LiteyukiStudio/spage/spage/models"
	store "github.com/LiteyukiStudio/spage/spage/store"
	"github.com/LiteyukiStudio/spage/utils"
	"strconv"

	"github.com/LiteyukiStudio/spage/resps"
	"github.com/cloudwego/hertz/pkg/app"
)

type OrgApi struct{}

var Org = OrgApi{}

// OrganizationDTO 组织信息数据传输对象
func getOrg(ctx context.Context) *models.Organization {
	org, ok := ctx.Value("userOrg").(*models.Organization)
	if !ok {
		return nil
	}
	return org
}

// ToDTO 组织信息数据传输对象
func (OrgApi) ToDTO(org *models.Organization) OrganizationDTO {
	return OrganizationDTO{
		ID:           org.ID,
		Name:         org.Name,
		DisplayName:  org.DisplayName,
		Email:        org.Email,
		Description:  org.Description,
		AvatarURL:    org.AvatarURL,
		ProjectLimit: org.ProjectLimit,
	}
}

// CreateOrganization 创建组织
func (OrgApi) CreateOrganization(ctx context.Context, c *app.RequestContext) {
	// 绑定参数
	req := &CreateOrgReq{}
	if err := c.BindAndValidate(&req); err != nil {
		resps.BadRequest(c, resps.ParameterError)
		return
	}
	// 校验组织名称是否合法
	if !utils.IsValidEntityName(req.Name) {
		resps.BadRequest(c, "organization name is invalid")
		return
	}
	// 检验组织名称是否存在
	if store.Org.OrgNameIsExist(req.Name) {
		resps.BadRequest(c, "organization name already exists")
		return
	}
	// 创建组织
	user := middle.Auth.GetUserWithBlock(ctx, c)
	org := models.Organization{
		Name:        req.Name,
		DisplayName: &req.DisplayName,
		Email:       req.Email,
		Description: req.Description,
		AvatarURL:   req.AvatarURL,
		Members:     []*models.User{user},
		Owners:      []models.User{*user},
	}
	if err := store.Org.CreateOrg(&org); err != nil {
		resps.InternalServerError(c, err.Error())
		return
	}
	resps.Ok(c, resps.OK, map[string]any{
		"organization": Org.ToDTO(&org),
	})
}

// UpdateOrganization 更新组织信息
func (OrgApi) UpdateOrganization(ctx context.Context, c *app.RequestContext) {
	req := &UpdateOrgReq{}
	if err := c.BindAndValidate(&req); err != nil {
		resps.BadRequest(c, resps.ParameterError)
		return
	}
	// 校验组织名称是否合法
	if !utils.IsValidEntityName(req.Name) {
		resps.BadRequest(c, "organization name is invalid")
		return
	}
	org := getOrg(ctx)
	// 更新 Update
	org.DisplayName = req.DisplayName
	org.Email = req.Email
	org.Description = *req.Description
	org.AvatarURL = req.AvatarURL
	if err := store.Org.UpdateOrg(org); err != nil {
		resps.InternalServerError(c, err.Error())
		return
	}
	resps.Ok(c, resps.OK, map[string]any{
		"organization": Org.ToDTO(org),
	})
}

// GetOrganizationProject 获取组织项目
func (OrgApi) GetOrganizationProject(ctx context.Context, c *app.RequestContext) {
	req := &GetOrgProjectReq{}
	if err := c.BindAndValidate(&req); err != nil {
		resps.BadRequest(c, resps.ParameterError)
		return
	}
	org := getOrg(ctx)
	// 查询 Query
	projects, total, err := store.Project.ListByOwner(constants.OwnerTypeOrg, strconv.Itoa(int(org.ID)), req.Page, req.Limit)
	if err != nil {
		resps.InternalServerError(c, "Failed to get projects")
		return
	}
	resps.Ok(c, resps.OK, map[string]any{
		"projects": func() (projectDTOs []ProjectDTO) {
			for _, project := range projects {
				projectDTOs = append(projectDTOs, Project.toDTO(&project, false))
			}
			return
		}(),
		"total": total,
	})
}

// DeleteOrganization 删除组织
func (OrgApi) DeleteOrganization(ctx context.Context, c *app.RequestContext) {
	org := getOrg(ctx)
	// 删除组织
	if err := store.Org.DeleteOrg(org); err != nil {
		resps.InternalServerError(c, err.Error())
		return
	}
	resps.Ok(c, resps.OK)
}

// GetOrganization 获取组织信息
func (OrgApi) GetOrganization(ctx context.Context, c *app.RequestContext) {
	org := getOrg(ctx)
	resps.Ok(c, resps.OK, map[string]any{
		"organization": Org.ToDTO(org),
	})
}

// GetOrganizationUsers 获取组织用户
func (OrgApi) GetOrganizationUsers(ctx context.Context, c *app.RequestContext) {
	org := getOrg(ctx)
	resps.Ok(c, resps.OK, map[string]any{
		"members": func() (users []UserDTO) {
			for _, user := range org.Members {
				users = append(users, User.ToDTO(user, false))
			}
			return
		},
		"owners": func() (users []UserDTO) {
			for _, user := range org.Owners {
				users = append(users, User.ToDTO(&user, false))
			}
			return
		},
	})
}

// AddOrganizationUser 添加组织用户
func (OrgApi) AddOrganizationUser(ctx context.Context, c *app.RequestContext) {
	req := &OrgUserReq{}
	if err := c.BindAndValidate(&req); err != nil {
		resps.BadRequest(c, resps.ParameterError)
		return
	}
	// 查询 Query
	user, err := store.User.GetByID(req.UserID)
	if err != nil || user == nil {
		resps.NotFound(c, resps.TargetNotFound)
		return
	}
	// 新增 Add
	org := getOrg(ctx)
	if req.Role == "member" {
		org.Members = append(org.Members, user)
		if err = store.Org.UpdateOrg(org); err != nil {
			resps.InternalServerError(c, err.Error())
			return
		}
	} else if req.Role == "owner" {
		org.Owners = append(org.Owners, *user)
		if err = store.Org.UpdateOrg(org); err != nil {
			resps.InternalServerError(c, err.Error())
			return
		}
	}
	resps.Ok(c, resps.OK)
}

// DeleteOrganizationUser 删除组织用户
func (OrgApi) DeleteOrganizationUser(ctx context.Context, c *app.RequestContext) {
	req := &OrgUserReq{}
	if err := c.BindAndValidate(&req); err != nil {
		resps.BadRequest(c, resps.ParameterError)
		return
	}
	// 查询 Query
	user, err := store.User.GetByID(req.UserID)
	if err != nil || user == nil {
		resps.NotFound(c, resps.TargetNotFound)
		return
	}
	// 删除 Delete
	org := getOrg(ctx)
	if req.Role == "member" {
		for i, member := range org.Members {
			if member.ID == user.ID {
				org.Members = append(org.Members[:i], org.Members[i+1:]...)
				break
			}
		}
		if err = store.Org.UpdateOrg(org); err != nil {
			resps.InternalServerError(c, err.Error())
			return
		}
	} else if req.Role == "owner" {
		for i, owner := range org.Owners {
			if owner.ID == user.ID {
				org.Owners = append(org.Owners[:i], org.Owners[i+1:]...)
				break
			}
			if err = store.Org.UpdateOrg(org); err != nil {
				resps.InternalServerError(c, err.Error())
				return
			}
		}
	}
	resps.Ok(c, resps.OK)
}

// IsOrgMemberMiddleware 组织成员中间件
func (OrgApi) IsOrgMemberMiddleware(ctx context.Context, c *app.RequestContext) {
	orgIdString := c.Param("id")
	orgId, err := strconv.Atoi(orgIdString)
	if err != nil {
		resps.BadRequest(c, resps.ParameterError)
		return
	}
	crtUser := middle.Auth.GetUserWithBlock(ctx, c)
	pass, err := isOrgMember(crtUser.ID, uint(orgId))
	if err != nil || !pass {
		resps.Forbidden(c, resps.PermissionDenied)
		return
	}
	c.Next(ctx)
	return
}

func isOrgMember(userId, orgId uint) (bool, error) {
	org, err := store.Org.GetOrgById(orgId)
	if err != nil {
		return false, err
	}
	for _, member := range org.Members {
		if member.ID == userId {
			return true, nil
		}
	}
	return false, nil
}

// IsOrgOwnerMiddleware 组织所有者中间件
func (OrgApi) IsOrgOwnerMiddleware(ctx context.Context, c *app.RequestContext) {
	orgIdString := c.Param("id")
	orgId, err := strconv.Atoi(orgIdString)
	if err != nil {
		resps.BadRequest(c, resps.ParameterError)
		return
	}
	crtUser := middle.Auth.GetUserWithBlock(ctx, c)
	pass, err := isOrgOwner(crtUser.ID, uint(orgId))
	if err != nil || !pass {
		resps.Forbidden(c, resps.PermissionDenied)
		return
	}
	c.Next(ctx)
	return
}

func isOrgOwner(userId, orgId uint) (bool, error) {
	// owners 是一定在 members 子集中
	org, err := store.Org.GetOrgById(orgId)
	if err != nil {
		return false, err
	}
	for _, owner := range org.Owners {
		if owner.ID == userId {
			return true, nil
		}
	}
	return false, nil
}
