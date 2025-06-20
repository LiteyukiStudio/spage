package handlers

import (
	"context"
	"github.com/LiteyukiStudio/spage/config"
	"github.com/LiteyukiStudio/spage/constants"
	"github.com/LiteyukiStudio/spage/spage/middle"
	"github.com/LiteyukiStudio/spage/spage/models"
	"github.com/LiteyukiStudio/spage/spage/store"
	"github.com/LiteyukiStudio/spage/utils"
	"strconv"
	"time"

	"github.com/LiteyukiStudio/spage/resps"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol"
)

type userApi struct{}

var User = userApi{}

// ToDTO 用户信息数据传输对象
func (userApi) ToDTO(user *models.User, self bool) UserDTO {
	userDTO := UserDTO{
		ID:          user.ID,
		Name:        user.Name,
		DisplayName: user.DisplayName,
		Email:       user.Email,
		Description: user.Description,
		Avatar:      user.AvatarURL,
	}
	if self {
		userDTO.Role = user.Role
		userDTO.Language = user.Language
	}
	return userDTO
}

// Login 用户登录
func (userApi) Login(ctx context.Context, c *app.RequestContext) {
	loginReq := &LoginReq{}
	err := c.BindJSON(loginReq)
	if err != nil {
		resps.BadRequest(c, "Parameter error")
		return
	}
	if loginReq.Username == "" || loginReq.Password == "" {
		resps.BadRequest(c, "Username or password cannot be empty")
		return
	}
	user, err := store.User.GetByName(loginReq.Username)
	if err != nil {
		user, err = store.User.GetByEmail(loginReq.Username)
		if err != nil {
			resps.BadRequest(c, "User does not exist")
			return
		}
	}
	if user.Password == nil {
		resps.Forbidden(c, "Password not set, please use another login method")
		return
	} else {
		if utils.Password.VerifyPassword(loginReq.Password, *user.Password, config.JwtSecret) {
			middle.Auth.SetTokenForCookie(c, user, true, loginReq.Remember)
		} else {
			resps.Forbidden(c, "Incorrect password")
			return
		}
	}
}

// Logout 用户登出
func (userApi) Logout(ctx context.Context, c *app.RequestContext) {
	// 删除cookie
	c.SetCookie("token", "", -1, "/", "", protocol.CookieSameSiteLaxMode, true, true)
	c.SetCookie("refresh_token", "", -1, "/", "", protocol.CookieSameSiteLaxMode, true, true)
	resps.Ok(c, "Logout successful")
}

// GetCaptchaConfig 获取验证码
func (userApi) GetCaptchaConfig(ctx context.Context, c *app.RequestContext) {
	resps.Ok(c, "ok", map[string]any{
		"provider": config.CaptchaType,
		"site_key": config.CaptchaSiteKey,
		"url":      config.CaptchaUrl,
	})
}

// GetOrgs 获取用户的组织
func (userApi) GetOrgs(ctx context.Context, c *app.RequestContext) {
	userID := c.Param("id")
	crtUser := middle.Auth.GetUserWithBlock(ctx, c)
	if userID != strconv.Itoa(int(crtUser.ID)) {
		resps.Forbidden(c, resps.PermissionDenied)
		return
	}
	page, limit := utils.Ctx.GetPageLimit(c)

	orgs, err := store.Org.ListByUserID(userID, page, limit)
	if err != nil {
		resps.InternalServerError(c, "Failed to get organizations")
		return
	}

	resps.Ok(c, resps.OK, map[string]any{
		"organizations": func() (orgDTOs []OrganizationDTO) {
			for _, org := range orgs {
				orgDTOs = append(orgDTOs, Org.ToDTO(&org))
			}
			return
		}(),
	})
}

// GetProjects 获取用户的项目
func (userApi) GetProjects(ctx context.Context, c *app.RequestContext) {
	userID := c.Param("id")
	crtUser := middle.Auth.GetUserWithBlock(ctx, c)
	if userID != strconv.Itoa(int(crtUser.ID)) {
		resps.Forbidden(c, resps.PermissionDenied)
	}
	page, limit := utils.Ctx.GetPageLimit(c)

	projects, total, err := store.Project.ListByOwner(constants.OwnerTypeUser, userID, page, limit)
	if err != nil {
		resps.InternalServerError(c, "Failed to get projects")
	}
	resps.Ok(c, resps.OK, map[string]any{
		"projects": func() (projectDTOs []ProjectDTO) {
			for _, project := range projects {
				projectDTOs = append(projectDTOs, Project.toDTO(&project, false)) // 这里naloveyuki尝试修复了下gitd7b49ff的问题
			}
			return
		}(),
		"total": total,
	})
}

// GetUser 获取用户信息
func (userApi) GetUser(ctx context.Context, c *app.RequestContext) {
	userIDString := c.Param("id")
	crtUser := middle.Auth.GetUserWithBlock(ctx, c)
	// 如果未提供用户ID或ID为空，则获取当前登录用户的信息
	if userIDString == "" {
		resps.Ok(c, "ok", map[string]any{
			"user": User.ToDTO(crtUser, true),
		})
		return
	}
	// 尝试将ID转换为整数
	userID, err := strconv.Atoi(userIDString)
	if err != nil {
		resps.BadRequest(c, "Invalid user ID")
		return
	}
	// 如果请求的是当前用户自己的信息
	if uint(userID) == crtUser.ID {
		resps.Ok(c, "ok", map[string]any{
			"user": User.ToDTO(crtUser, true),
		})
		return
	}
	// 获取其他用户的信息
	targetUser, err := store.User.GetByID(uint(userID))
	if err != nil {
		resps.NotFound(c, "User not found")
		return
	}
	resps.Ok(c, "ok", map[string]any{
		"user": User.ToDTO(targetUser, false),
	})
}

// Register 用户注册
func (userApi) Register(ctx context.Context, c *app.RequestContext) {
	// 接收参数
	request := &RegisterReq{}
	err := c.BindJSON(request)
	if err != nil {
		resps.BadRequest(c, "Parameter error")
		return
	}
	// TODO 校验邮箱验证码
	// 校验密码复杂度
	passwordLevel := config.GetInt("password_complexity", 3)
	if !utils.Password.CheckPasswordComplexity(request.Password, passwordLevel) {
		resps.BadRequest(c, "Password complexity is too low")
		return
	}
	// 判断用户名是否存在
	if store.Owner.IsNameAvailable(request.Username) {
		resps.BadRequest(c, "Username already exists")
		return
	}
	// 创建用户
	hashPassword, err := utils.Password.HashPassword(request.Password, config.JwtSecret)
	if err != nil {
		resps.InternalServerError(c, "Failed to hash password")
		return
	}
	err = store.User.Create(&models.User{
		Name:     request.Username,
		Email:    &request.Email,
		Password: &hashPassword,
	})
	if err != nil {
		resps.InternalServerError(c, "Failed to create user")
		return
	}
	resps.Ok(c, "Register successful", map[string]any{
		"user": UserDTO{
			Name:  request.Username,
			Email: &request.Email,
		},
	})
}

// UpdateUser 更新用户信息
func (userApi) UpdateUser(ctx context.Context, c *app.RequestContext) {
	userDTO := &UserDTO{}
	if err := c.BindJSON(userDTO); err != nil {
		resps.BadRequest(c, resps.ParameterError)
		return
	}
	// 校验用户名称是否合法
	if !store.Owner.IsNameAvailable(userDTO.Name) {
		resps.BadRequest(c, "Username is invalid")
		return
	}
	crtUser := middle.Auth.GetUserWithBlock(ctx, c)
	crtUser.Name = userDTO.Name
	crtUser.DisplayName = userDTO.DisplayName
	crtUser.Email = userDTO.Email
	crtUser.Description = userDTO.Description
	crtUser.AvatarURL = userDTO.Avatar
	crtUser.Language = userDTO.Language
	if err := store.User.Update(crtUser); err != nil {
		resps.InternalServerError(c, "Failed to update user")
		return
	}
	resps.Ok(c, resps.OK, map[string]any{})
}

// CreateApiToken 创建API token
func (userApi) CreateApiToken(ctx context.Context, c *app.RequestContext) {
	createTokenReq := &CreateTokenReq{}
	err := c.BindJSON(createTokenReq)
	if err != nil {
		resps.BadRequest(c, resps.ParameterError)
		return
	}
	crtUser := middle.Auth.GetUserWithBlock(ctx, c)
	token, err := utils.Token.CreateApiToken(crtUser.ID, time.Duration(createTokenReq.Expire)*time.Second, middle.ApiTokenPersistentHandler)
	if err != nil {
		resps.InternalServerError(c, "Failed to create token: "+err.Error())
		return
	}
	resps.Ok(c, resps.OK, map[string]any{"token": token})
}

// RevokeApiToken 吊销删除Api Token
func (userApi) RevokeApiToken(ctx context.Context, c *app.RequestContext) {
	tokenIDString := c.Param("id")
	tokenId, err := strconv.Atoi(tokenIDString)
	if err != nil {
		resps.BadRequest(c, resps.ParameterError)
		return
	}
	if tokenId == 0 {
		resps.BadRequest(c, "Invalid token ID")
		return
	}
	crtUser := middle.Auth.GetUserWithBlock(ctx, c)
	token, err := store.Token.GetApiTokenByID(uint(tokenId))
	if err != nil || token == nil {
		resps.InternalServerError(c, resps.TargetNotFound)
		return
	}
	if crtUser.ID != token.UserID {
		resps.Forbidden(c, resps.PermissionDenied)
		return
	}

	err = store.Token.RevokeApiTokenByID(uint(tokenId))
	if err != nil {
		resps.InternalServerError(c, "Failed to revoke token: "+err.Error())
		return
	}
	resps.Ok(c, resps.OK, map[string]any{})
}

// ListApiToken 列出ApiToken
func (userApi) ListApiToken(ctx context.Context, c *app.RequestContext) {
	page, limit := utils.Ctx.GetPageLimit(c)
	crtUser := middle.Auth.GetUserWithBlock(ctx, c)
	tokens, total, err := store.Token.ListApiTokens(crtUser.ID, page, limit)
	if err != nil {
		resps.BadRequest(c, resps.ParameterError)
		return
	}
	resps.Ok(c, resps.OK, map[string]any{
		"tokens": func() []map[string]any {
			var tokenDTOs []map[string]any
			for _, token := range tokens {
				tokenDTOs = append(tokenDTOs, map[string]any{
					"id":         token.ID,
					"created_at": token.CreatedAt.Unix(),
					"expires_at": token.ExpiresAt.Unix(),
				})
			}
			return tokenDTOs
		}(),
		"total": total,
	})
}
