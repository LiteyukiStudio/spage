package handlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/LiteyukiStudio/spage/config"
	"github.com/LiteyukiStudio/spage/constants"
	"github.com/LiteyukiStudio/spage/resps"
	"github.com/LiteyukiStudio/spage/spage/middle"
	"github.com/LiteyukiStudio/spage/spage/store"
	"github.com/LiteyukiStudio/spage/utils"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/url"
	"time"
)

type oidcType struct{}

var Oidc = oidcType{}

func (oidcType) ListOidcConfig(ctx context.Context, c *app.RequestContext) {
	oidcConfigs, err := store.Oidc.ListEnabledOidcConfig()
	if err != nil {
		resps.InternalServerError(c, err.Error())
		return
	}
	resps.Ok(c, resps.OK, map[string]any{
		"oidc_configs": func() []map[string]any {
			var configsDto []map[string]any
			for _, oidcConfig := range oidcConfigs {
				state := utils.GenerateRandomString(32)
				// TODO 使用utils的键值内存储存和验证state
				kvStore := utils.GetKVStore()
				kvStore.Set(constants.KVKeyOidcState+state, oidcConfig.Name, 5*time.Minute)
				configsDto = append(configsDto, map[string]any{
					"id":           oidcConfig.ID,
					"display_name": oidcConfig.DisplayName,
					"icon":         oidcConfig.Icon,
					"login_url": buildURL(oidcConfig.AuthorizationEndpoint, map[string]string{
						"client_id":     oidcConfig.ClientID,
						"redirect_uri":  config.BaseUrl + config.OidcUri + "/" + oidcConfig.Name,
						"response_type": "code",
						"scope":         "openid email profile",
						"state":         state,
					}),
					"name": oidcConfig.Name,
				})
			}
			return configsDto
		}(),
	})
}

// requestToken 请求访问令牌
func requestToken(client *resty.Client, tokenEndpoint, clientID, clientSecret, code, redirectURI string) (*TokenResponse, error) {
	tokenResp, err := client.R().
		SetFormData(map[string]string{
			"grant_type":    "authorization_code",
			"client_id":     clientID,
			"client_secret": clientSecret,
			"code":          code,
			"redirect_uri":  redirectURI,
		}).
		SetHeader("Accept", "application/json").
		SetResult(&TokenResponse{}).
		Post(tokenEndpoint)

	if err != nil {
		return nil, err
	}

	if tokenResp.StatusCode() != 200 {
		return nil, fmt.Errorf("状态码: %d，响应: %s", tokenResp.StatusCode(), tokenResp.String())
	}
	return tokenResp.Result().(*TokenResponse), nil
}

// requestUserInfo 请求用户信息
func requestUserInfo(client *resty.Client, userInfoEndpoint, accessToken string) (*UserInfo, error) {
	userInfoResp, err := client.R().
		SetHeader("Authorization", "Bearer "+accessToken).
		SetHeader("Accept", "application/json").
		SetResult(&UserInfo{}).
		Get(userInfoEndpoint)

	if err != nil {
		return nil, err
	}

	if userInfoResp.StatusCode() != 200 {
		return nil, fmt.Errorf("状态码: %d，响应: %s", userInfoResp.StatusCode(), userInfoResp.String())
	}

	return userInfoResp.Result().(*UserInfo), nil
}

// LoginOidcConfig 主函数
func (oidcType) LoginOidcConfig(ctx context.Context, c *app.RequestContext) {
	name := c.Param("name")
	code := c.Query("code")
	state := c.Query("state")
	kvStore := utils.GetKVStore()
	v, ok := kvStore.Get(constants.KVKeyOidcState + state)
	if !ok || name != v {
		resps.BadRequest(c, "无效的OIDC state")
		return
	}

	oidcConfig, err := store.Oidc.GetByName(name)
	if err != nil || oidcConfig == nil {
		resps.NotFound(c, "OIDC配置未找到: "+name)
		return
	}
	if code == "" {
		resps.BadRequest(c, "缺少授权码")
		return
	}

	client := resty.New()
	tokenResult, err := requestToken(
		client,
		oidcConfig.TokenEndpoint,
		oidcConfig.ClientID,
		oidcConfig.ClientSecret,
		code,
		config.BaseUrl+config.OidcUri+"/"+oidcConfig.Name,
	)
	if err != nil {
		logrus.Errorf("获取访问令牌失败: %v", err)
		resps.InternalServerError(c, "获取访问令牌失败")
		return
	}
	userInfo, err := requestUserInfo(client, oidcConfig.UserInfoEndpoint, tokenResult.AccessToken)
	if err != nil {
		logrus.Errorf("获取用户信息失败: %v", err)
		resps.InternalServerError(c, "获取用户信息失败")
		return
	}
	if !store.Owner.IsNameAvailable(userInfo.Name) {
		userInfo.Name = utils.GenerateRandomString(4) + userInfo.Name
		logrus.Warnf("用户名 %s 已存在，已更改为 %s", userInfo.Name, userInfo.Name)
	}
	user, err := store.User.GetByEmail(userInfo.Email)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if !config.AllowRegisterByOidc {
			logrus.Warnf("用户 %s 不存在且不允许通过OIDC注册", userInfo.Email)
			resps.Forbidden(c, "不允许通过OIDC注册")
			return
		}
	}
	if !matchGroups(userInfo.Groups, oidcConfig.AllowedGroups, true) {
		resps.Forbidden(c, "用户不在允许的组中")
		return
	}
	user, err = store.User.FindOrCreateByEmail(userInfo.Email, userInfo.Name)
	if err != nil {
		logrus.Errorf("用户处理失败: %v", err)
		resps.InternalServerError(c, "用户处理失败")
		return
	}
	if matchGroups(userInfo.Groups, oidcConfig.AdminGroups, false) {
		user.Role = constants.GlobalRoleAdmin
		err = store.User.Update(user)
		if err != nil {
			logrus.Errorf("更新用户角色失败: %v", err)
			resps.InternalServerError(c, "更新用户角色失败")
			return
		}
	}
	middle.Auth.SetTokenForCookie(c, user, false, false)
	resps.Redirect(c, config.BaseUrl)
}

func buildURL(baseURL string, queryParams map[string]string) string {
	u, err := url.Parse(baseURL)
	if err != nil {
		return baseURL
	}
	q := u.Query()
	for key, value := range queryParams {
		q.Set(key, value)
	}
	u.RawQuery = q.Encode()
	return u.String()
}

// TokenResponse 定义访问令牌响应结构
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	IDToken      string `json:"id_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

// UserInfo 定义用户信息结构
type UserInfo struct {
	Sub     string   `json:"sub"`
	Name    string   `json:"name"`
	Email   string   `json:"email"`
	Picture string   `json:"picture,omitempty"`
	Groups  []string `json:"groups,omitempty"` // 可选字段，OIDC提供的用户组信息
}

// 匹配用户是否属于指定组列表中
func matchGroups(userGroups []string, groups []string, defaultMatch bool) bool {
	if len(groups) == 0 || len(groups) >= 1 && groups[0] == "*" {
		return defaultMatch
	}
	for _, userGroup := range userGroups {
		for _, group := range groups {
			if userGroup == group {
				return true
			}
		}
	}
	return false
}
