package handlers

import (
	"context"
	"strconv"

	"github.com/LiteyukiStudio/spage/constants"
	"github.com/LiteyukiStudio/spage/resps"
	"github.com/LiteyukiStudio/spage/spage/middle"
	"github.com/LiteyukiStudio/spage/spage/models"
	"github.com/LiteyukiStudio/spage/spage/store"
	"github.com/cloudwego/hertz/pkg/app"
)

// TODO 给项目name更新添加可用性检查，函数在store
type ProjectApi struct {
}

var Project = ProjectApi{}

// toDTO 项目信息数据传输对象
func (ProjectApi) toDTO(project *models.Project, full bool) ProjectDTO {
	projectDto := ProjectDTO{
		Description: project.Description,
		DisplayName: project.DisplayName,
		ID:          project.ID,
		Name:        project.Name,
		OwnerType:   project.OwnerType,
	}
	if full {
		projectDto.OwnerID = project.OwnerID
		projectDto.Owners = func([]models.User) (owners []UserDTO) {
			for _, owner := range project.Owners {
				owners = append(owners, User.ToDTO(&owner, false))
			}
			return
		}(project.Owners)
		projectDto.SiteLimit = project.SiteLimit
	}
	return projectDto
}

// getProject 获取项目信息
func getProject(ctx context.Context) *models.Project {
	project, ok := ctx.Value("userProject").(*models.Project)
	if !ok {
		return nil
	}
	return project
}

// UserProjectAuth 用户项目权限认证
func (ProjectApi) UserProjectAuth(ctx context.Context, c *app.RequestContext) {
	user := middle.Auth.GetUserWithBlock(ctx, c)
	projectIdStr := c.Param("id")
	// 当id为空默认为创建
	if projectIdStr == "" && string(c.Method()) == "POST" {
		return
	}
	// 查询项目
	projectId, err := strconv.Atoi(projectIdStr)
	if err != nil {
		resps.BadRequest(c, resps.ParameterError)
		c.Abort()
		return
	}
	project, err := store.Project.GetByID(uint(projectId))
	if err != nil || project == nil {
		resps.NotFound(c, resps.TargetNotFound)
		c.Abort()
		return
	}
	// 项目权限判断
	if store.Project.UserIsOwner(project, user.ID) {
		context.WithValue(ctx, "userProject", project)
		return
	}
	// 组织权限判断
	if project.OwnerType == constants.OwnerTypeOrg {
		// 组织查询
		org, err := store.Org.GetOrgById(project.OwnerID)
		if err != nil || org == nil {
			resps.InternalServerError(c, resps.ParameterError)
			c.Abort()
			return
		}
		// 请求类型判断
		var authType string
		if string(c.Method()) == "GET" {
			authType = "member"
		} else {
			authType = "owner"
		}
		if authType == store.Org.GetUserRole(org, user.ID) {
			context.WithValue(ctx, "userOrg", org)
			context.WithValue(ctx, "userProject", project)
			return
		} else {
			resps.BadRequest(c, resps.ParameterError)
			c.Abort()
			return
		}
	} else {
		resps.BadRequest(c, resps.ParameterError)
		c.Abort()
	}
}

// Create 创建项目
func (ProjectApi) Create(ctx context.Context, c *app.RequestContext) {
	req := CreateProjectReq{}
	if err := c.BindAndValidate(&req); err != nil {
		resps.BadRequest(c, resps.ParameterError)
		return
	}
	// 校验项目名称是否合法
	if !store.Project.CheckNameAvailable(req.OwnerType, req.OwnerID, req.Name) {
		resps.BadRequest(c, "Project name is invalid or already exists")
		return
	}
	currentUser := middle.Auth.GetUserWithBlock(ctx, c)
	switch req.OwnerType {
	case constants.OwnerTypeUser:
		if req.OwnerID != currentUser.ID {
			resps.BadRequest(c, resps.ParameterError)
			return
		}
	case constants.OwnerTypeOrg:
		isOrgMemberPass, err := isOrgMember(currentUser.ID, req.OwnerID)
		if !isOrgMemberPass || err != nil {
			resps.BadRequest(c, resps.ParameterError)
			return
		}
	default:
		resps.BadRequest(c, resps.ParameterError)
		return
	}
	project := &models.Project{
		Description: req.Description,
		DisplayName: req.DisplayName,
		Name:        req.Name,
		OwnerID:     req.OwnerID,
		OwnerType:   req.OwnerType,
		Owners:      []models.User{*currentUser},
		Members:     []*models.User{currentUser},
		IsPrivate:   req.IsPrivate,
	}
	if err := store.Project.Create(project); err != nil {
		resps.InternalServerError(c, resps.ParameterError)
		return
	}
	resps.Ok(c, resps.OK, map[string]any{
		"project": Project.toDTO(project, true),
	})
}

// Update 更新项目
func (ProjectApi) Update(ctx context.Context, c *app.RequestContext) {
	req := UpdateProjectReq{}
	if err := c.BindAndValidate(&req); err != nil {
		resps.BadRequest(c, resps.ParameterError)
		return
	}
	project := getProject(ctx)
	if project == nil {
		resps.NotFound(c, resps.TargetNotFound)
		return
	}
	// 校验项目名称是否合法
	if req.Name != nil && !store.Project.CheckNameAvailable(project.OwnerType, project.OwnerID, *req.Name) {
		resps.BadRequest(c, "Project name is invalid or already exists")
		return
	}
	// 更新数据
	project.Description = *req.Description
	project.DisplayName = req.DisplayName
	project.Name = *req.Name
	if err := store.Project.Update(project); err != nil {
		resps.InternalServerError(c, resps.ParameterError)
		return
	}
	resps.Ok(c, resps.OK, map[string]any{
		"project": Project.toDTO(project, true),
	})
}

// Delete 删除项目
func (ProjectApi) Delete(ctx context.Context, c *app.RequestContext) {
	project := getProject(ctx)
	if project == nil {
		resps.NotFound(c, resps.TargetNotFound)
		return
	}
	if err := store.Project.Delete(project); err != nil {
		resps.InternalServerError(c, resps.ParameterError)
		return
	}
	resps.Ok(c, resps.OK)
}

// Info 获取项目信息
func (ProjectApi) Info(ctx context.Context, c *app.RequestContext) {
	project := getProject(ctx)
	if project == nil {
		resps.NotFound(c, resps.TargetNotFound)
		return
	}
	resps.Ok(c, resps.OK, map[string]any{
		"project": Project.toDTO(project, true),
	})
}

// GetOwners 获取项目所有者列表
func (ProjectApi) GetOwners(ctx context.Context, c *app.RequestContext) {
	project := getProject(ctx)
	if project == nil {
		resps.NotFound(c, resps.TargetNotFound)
		return
	}
	resps.Ok(c, resps.OK, map[string]any{
		"owners": func([]models.User) (owners []UserDTO) {
			for _, owner := range project.Owners {
				owners = append(owners, User.ToDTO(&owner, false))
			}
			return
		}(project.Owners),
	})
}

// AddOwner 添加项目所有者
func (ProjectApi) AddOwner(ctx context.Context, c *app.RequestContext) {
	project := getProject(ctx)
	if project == nil {
		resps.NotFound(c, resps.TargetNotFound)
		return
	}
	req := ProjectUserReq{}
	if err := c.BindAndValidate(&req); err != nil {
		resps.BadRequest(c, resps.ParameterError)
		return
	}
	// 查询用户	查询用户
	user, err := store.User.GetByID(req.UserID)
	if err != nil || user == nil {
		resps.NotFound(c, resps.TargetNotFound)
		return
	}
	// 判断用户是否已存在权限列表
	for _, owner := range project.Owners {
		if owner.ID == user.ID {
			resps.BadRequest(c, "User already exists in the permission list")
			return
		}
	}
	// 添加用户
	if err := store.Project.AddOwner(project, user); err != nil {
		resps.InternalServerError(c, resps.ParameterError)
		return
	}
	resps.Ok(c, resps.OK, map[string]any{
		"project": Project.toDTO(project, true),
	})
}

// DeleteOwner 删除项目所有者
func (ProjectApi) DeleteOwner(ctx context.Context, c *app.RequestContext) {
	project := getProject(ctx)
	if project == nil {
		resps.NotFound(c, resps.TargetNotFound)
		return
	}
	req := ProjectUserReq{}
	if err := c.BindAndValidate(&req); err != nil {
		resps.BadRequest(c, resps.ParameterError)
		return
	}
	// 查询用户
	user, err := store.User.GetByID(req.UserID)
	if err != nil || user == nil {
		resps.NotFound(c, resps.TargetNotFound)
		return
	}
	// 删除用户
	if err := store.Project.DeleteOwner(project, user); err != nil {
		resps.InternalServerError(c, resps.ParameterError)
		return
	}
	resps.Ok(c, resps.OK, map[string]any{
		"project": Project.toDTO(project, true),
	})
}

// GetSites 获取站点列表
func (ProjectApi) GetSites(ctx context.Context, c *app.RequestContext) {
	req := GetSiteListReq{}
	if err := c.BindAndValidate(&req); err != nil {
		resps.BadRequest(c, resps.ParameterError)
		return
	}
	project := getProject(ctx)
	if project == nil {
		resps.NotFound(c, resps.TargetNotFound)
		return
	}
	sites, total, err := store.Project.GetSiteList(project, req.Page, req.Limit)
	if err != nil {
		resps.InternalServerError(c, "Failed to get sites")
		return
	}
	resps.Ok(c, resps.OK, map[string]any{
		"sites": func([]models.Site) (siteDTOs []SiteDTO) {
			for _, site := range sites {
				siteDTOs = append(siteDTOs, Site.ToDTO(&site, false))
			}
			return
		}(sites),
		"total": total,
	})
}

// IsProjectOwnerMiddleware 项目所有者中间件
func (ProjectApi) IsProjectOwnerMiddleware(ctx context.Context, c *app.RequestContext) {
	currentUser := middle.Auth.GetUserWithBlock(ctx, c)
	projectIdStr := c.Param("id")
	if projectIdStr == "" {
		resps.BadRequest(c, resps.ParameterError)
		c.Abort()
		return
	}
	projectId, err := strconv.Atoi(projectIdStr)
	if err != nil {
		resps.BadRequest(c, resps.ParameterError)
		c.Abort()
		return
	}
	isOwner, err := isProjectOwner(currentUser.ID, uint(projectId))
	if err != nil {
		resps.InternalServerError(c, err.Error())
		c.Abort()
		return
	}
	if !isOwner {
		resps.Forbidden(c, resps.PermissionDenied)
		c.Abort()
		return
	}
	c.Next(ctx)
}

func isProjectOwner(userId, projectId uint) (bool, error) {
	project, err := store.Project.GetByID(projectId)
	if err != nil {
		return false, err
	}
	// 情况1: 用户是项目直接拥有者
	if project.OwnerType == constants.OwnerTypeUser && project.OwnerID == userId {
		return true, nil
	}
	// 情况2: 用户在项目所有者列表中
	for _, owner := range project.Owners {
		if owner.ID == userId {
			return true, nil
		}
	}
	// 情况3: 组织项目 - 检查用户是否为组织所有者
	if project.OwnerType == constants.OwnerTypeOrg {
		org, _ := store.Org.GetOrgById(project.OwnerID)
		if org != nil && store.Org.GetUserRole(org, userId) == constants.OrgRoleOwner {
			return true, nil
		}
	}
	return false, nil
}

// IsProjectMemberMiddleware 项目成员中间件
func (ProjectApi) IsProjectMemberMiddleware(ctx context.Context, c *app.RequestContext) {
	currentUser := middle.Auth.GetUserWithBlock(ctx, c)
	projectIdStr := c.Param("id")
	if projectIdStr == "" {
		resps.BadRequest(c, resps.ParameterError)
		c.Abort()
		return
	}
	projectId, err := strconv.Atoi(projectIdStr)
	if err != nil {
		resps.BadRequest(c, resps.ParameterError)
		c.Abort()
		return
	}
	isMember, err := isProjectMember(currentUser.ID, uint(projectId))
	if err != nil {
		resps.InternalServerError(c, err.Error())
		c.Abort()
		return
	}
	if !isMember {
		resps.Forbidden(c, resps.PermissionDenied)
		c.Abort()
		return
	}
	c.Next(ctx)
}

func isProjectMember(userId, projectId uint) (bool, error) {
	// 情况1: 用户是项目直接拥有者
	isOwner, _ := isProjectOwner(userId, projectId)
	if isOwner {
		return true, nil
	}
	// 用户在项目成员列表中
	project, err := store.Project.GetByID(projectId)
	if err != nil {
		return false, err
	}
	for _, member := range project.Members {
		if member.ID == userId {
			return true, nil
		}
	}
	return false, nil
}

// HasReadPermissionMiddleware 项目读取权限中间件
func (ProjectApi) HasReadPermissionMiddleware(ctx context.Context, c *app.RequestContext) {
	currentUser := middle.Auth.GetUserWithBlock(ctx, c)
	projectIdStr := c.Param("id")
	if projectIdStr == "" {
		resps.BadRequest(c, resps.ParameterError)
		c.Abort()
		return
	}
	projectId, err := strconv.Atoi(projectIdStr)
	if err != nil {
		resps.BadRequest(c, resps.ParameterError)
		c.Abort()
		return
	}
	hasPermission, err := hasReadPermission(currentUser.ID, uint(projectId))
	if err != nil {
		resps.InternalServerError(c, err.Error())
		c.Abort()
		return
	}
	if !hasPermission {
		resps.Forbidden(c, resps.PermissionDenied)
		c.Abort()
		return
	}
	c.Next(ctx)
}

func hasReadPermission(userId, projectId uint) (bool, error) {
	isMember, _ := isProjectMember(userId, projectId)
	if isMember {
		return true, nil
	}
	// 用户是组织的成员，也对组织项目可见，但没有其他权限
	org, err := store.Org.GetOrgById(projectId)
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
