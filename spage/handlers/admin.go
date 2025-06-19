package handlers

import (
	"context"
	"fmt"
	"github.com/LiteyukiStudio/spage/resps"
	"github.com/LiteyukiStudio/spage/spage/models"
	"github.com/LiteyukiStudio/spage/spage/store"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/go-resty/resty/v2"
	"strconv"
	"time"
)

type adminApi struct{}

var Admin = adminApi{}

// CreateUser 创建用户
func (adminApi) CreateUser(ctx context.Context, c *app.RequestContext) {
	var userDTO *UserDTO
	err := c.BindJSON(&userDTO)
	if err != nil {
		resps.BadRequest(c, resps.ParameterError)
		return
	}
}

func (adminApi) CreateOidcConfig(ctx context.Context, c *app.RequestContext) {
	var oidcDtp = OidcDto{}
	err := c.BindJSON(&oidcDtp)
	if err != nil {
		resps.BadRequest(c, resps.ParameterError)
		return
	}
	oidcConfig := &models.OIDCConfig{
		Name:             oidcDtp.Name,
		AdminGroups:      oidcDtp.AdminGroups,
		AllowedGroups:    oidcDtp.AllowedGroups,
		ClientID:         oidcDtp.ClientID,
		ClientSecret:     oidcDtp.ClientSecret,
		DisplayName:      oidcDtp.DisplayName,
		GroupsClaim:      &oidcDtp.GroupClaims,
		Icon:             &oidcDtp.Icon,
		OidcDiscoveryUrl: oidcDtp.OidcDiscoveryUrl,
		Enabled:          true,
	}

	err = updateOidcConfigFromUrl(oidcConfig.OidcDiscoveryUrl, oidcConfig)
	if err != nil {
		resps.BadRequest(c, fmt.Sprintf("请求OIDC发现端点失败: %v", err))
		return
	}

	if err = store.Oidc.CreateOidcConfig(oidcConfig); err != nil {
		resps.InternalServerError(c, err.Error())
		return
	}
	resps.Ok(c, resps.OK, map[string]any{"id": oidcConfig.ID})
}

func (adminApi) UpdateOidcConfig(ctx context.Context, c *app.RequestContext) {
	idString := c.Param("id")
	id, err := strconv.Atoi(idString)
	// 验证ID是否存在
	existingConfig, err := store.Oidc.GetByID(uint(id))
	if err != nil || existingConfig == nil {
		existingConfig = &models.OIDCConfig{}
	}

	if err != nil {
		resps.BadRequest(c, resps.ParameterError)
		return
	}
	var oidcDtp = OidcDto{}
	err = c.BindJSON(&oidcDtp)
	if err != nil {
		resps.BadRequest(c, resps.ParameterError)
		return
	}
	existingConfig = &models.OIDCConfig{
		Name:          oidcDtp.Name,
		AdminGroups:   oidcDtp.AdminGroups,
		AllowedGroups: oidcDtp.AllowedGroups,
		ClientID:      oidcDtp.ClientID,
		ClientSecret:  oidcDtp.ClientSecret,
		DisplayName:   oidcDtp.DisplayName,
		GroupsClaim:   &oidcDtp.GroupClaims,
		Icon:          &oidcDtp.Icon,
		OidcDiscoveryUrl: func() string {
			if oidcDtp.OidcDiscoveryUrl != "" {
				return oidcDtp.OidcDiscoveryUrl
			}
			return existingConfig.OidcDiscoveryUrl
		}(),
		Enabled: true,
	}
	existingConfig.ID = uint(id)
	err = updateOidcConfigFromUrl(existingConfig.OidcDiscoveryUrl, existingConfig)
	if err != nil {
		resps.BadRequest(c, fmt.Sprintf("请求OIDC发现端点失败: %v", err))
		return
	}
	if err = store.Oidc.UpdateOidcConfig(existingConfig); err != nil {
		resps.InternalServerError(c, err.Error())
		return
	}
	resps.Ok(c, resps.OK, map[string]any{"id": existingConfig.ID})
}

type oidcDiscoveryResp struct {
	Issuer                string `json:"issuer" validate:"required"`
	AuthorizationEndpoint string `json:"authorization_endpoint" validate:"required"`
	TokenEndpoint         string `json:"token_endpoint" validate:"required"`
	UserInfoEndpoint      string `json:"userinfo_endpoint" validate:"required"`
	JwksUri               string `json:"jwks_uri" validate:"required"`
	// 可选字段
	RegistrationEndpoint             string   `json:"registration_endpoint,omitempty"`
	ScopesSupported                  []string `json:"scopes_supported,omitempty"`
	ResponseTypesSupported           []string `json:"response_types_supported,omitempty"`
	GrantTypesSupported              []string `json:"grant_types_supported,omitempty"`
	SubjectTypesSupported            []string `json:"subject_types_supported,omitempty"`
	IdTokenSigningAlgValuesSupported []string `json:"id_token_signing_alg_values_supported,omitempty"`
	ClaimsSupported                  []string `json:"claims_supported,omitempty"`
	EndSessionEndpoint               string   `json:"end_session_endpoint,omitempty"`
}

func updateOidcConfigFromUrl(url string, config *models.OIDCConfig) error {
	// 创建 resty 客户端
	client := resty.New()
	client.SetTimeout(10 * time.Second) // 设置超时时间
	// 声明结果变量
	var discovery oidcDiscoveryResp
	// 发起 GET 请求并直接将结果解析到 discovery 变量
	resp, err := client.R().
		SetHeader("Accept", "application/json").
		SetResult(&discovery).
		Get(url)
	if err != nil {
		return err
	}
	if resp.StatusCode() != 200 {
		return fmt.Errorf("请求OIDC发现端点失败，状态码: %d", resp.StatusCode())
	}
	// 验证必要字段
	if discovery.Issuer == "" ||
		discovery.AuthorizationEndpoint == "" ||
		discovery.TokenEndpoint == "" ||
		discovery.UserInfoEndpoint == "" ||
		discovery.JwksUri == "" {
		return fmt.Errorf("OIDC发现响应缺少必要字段")
	}
	config.Issuer = discovery.Issuer
	config.AuthorizationEndpoint = discovery.AuthorizationEndpoint
	config.TokenEndpoint = discovery.TokenEndpoint
	config.UserInfoEndpoint = discovery.UserInfoEndpoint
	config.JwksUri = discovery.JwksUri
	return nil
}
