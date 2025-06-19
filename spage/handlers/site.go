package handlers

import (
	"context"
	"strconv"

	"github.com/LiteyukiStudio/spage/spage/models"
	"github.com/LiteyukiStudio/spage/spage/store"

	"github.com/LiteyukiStudio/spage/resps"
	"github.com/cloudwego/hertz/pkg/app"
)

type SiteApi struct {
}

var Site = SiteApi{}

// ToDTO 站点信息数据传输对象
func (SiteApi) ToDTO(site *models.Site, full bool) SiteDTO {
	siteDTO := SiteDTO{
		Description: site.Description,
		ID:          site.ID,
		Name:        site.Name,
	}
	if full {
		siteDTO.Project = Project.toDTO(&site.Project, full)
		siteDTO.SubDomain = &site.SubDomain
		siteDTO.Domains = site.Domains
	}
	return siteDTO
}

func getSite(ctx context.Context) *models.Site {
	site, ok := ctx.Value("userSite").(*models.Site)
	if !ok {
		return nil
	}
	return site
}

func (SiteApi) SiteAuth(ctx context.Context, c *app.RequestContext) {
	siteIDStr := c.Param("id")
	// 当id为空默认为创建
	if siteIDStr == "" && string(c.Method()) == "POST" {
		return
	}
	// 获取站点信息
	siteID, err := strconv.Atoi(siteIDStr)
	if err != nil {
		resps.BadRequest(c, resps.ParameterError)
		return
	}
	site, err := store.Site.GetByID(uint(siteID))
	if err != nil {
		resps.NotFound(c, "Site not found")
		return
	}
	context.WithValue(ctx, "userSite", site)
}

// Create 创建站点
func (SiteApi) Create(ctx context.Context, c *app.RequestContext) {
	req := CreateSiteReq{}
	if err := c.BindAndValidate(&req); err != nil {
		resps.BadRequest(c, resps.ParameterError)
		return
	}
	site := models.Site{
		Name:        req.Name,
		Description: req.Description,
		Domains:     req.Domains,
		ProjectID:   req.ProjectID,
		SubDomain:   *req.SubDomain,
	}
	if err := store.Site.Create(&site); err != nil {
		resps.InternalServerError(c, err.Error())
		return
	}
	// TODO 创建站点信息
	resps.Ok(c, resps.OK, map[string]any{
		"site": Site.ToDTO(&site, true),
	})
}

func (SiteApi) Update(ctx context.Context, c *app.RequestContext) {
	req := UpdateSiteReq{}
	if err := c.BindAndValidate(&req); err != nil {
		resps.BadRequest(c, resps.ParameterError)
		return
	}
	site := getSite(ctx)
	if site == nil {
		resps.NotFound(c, resps.TargetNotFound)
		return
	}
	site.Description = *req.Description
	site.Domains = req.Domains
	site.Name = *req.Name
	site.SubDomain = *req.SubDomain
	if err := store.Site.Update(site); err != nil {
		resps.InternalServerError(c, resps.ParameterError)
		return
	}
	// TODO 更新站点信息
	resps.Ok(c, resps.OK, map[string]any{
		"site": Site.ToDTO(site, true),
	})
}

func (SiteApi) Delete(ctx context.Context, c *app.RequestContext) {
	site := getSite(ctx)
	if site == nil {
		resps.NotFound(c, resps.TargetNotFound)
		return
	}
	if err := store.Site.Delete(site); err != nil {
		resps.InternalServerError(c, resps.ParameterError)
		return
	}
	// TODO 删除站点
	resps.Ok(c, resps.OK, map[string]any{
		"site": Site.ToDTO(site, true),
	})
}

func (SiteApi) Info(ctx context.Context, c *app.RequestContext) {
	site := getSite(ctx)
	if site == nil {
		resps.NotFound(c, resps.TargetNotFound)
		return
	}
	resps.Ok(c, resps.OK, map[string]any{
		"site": Site.ToDTO(site, true),
	})
}
