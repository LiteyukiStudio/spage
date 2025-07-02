package router

import (
	"github.com/LiteyukiStudio/spage/spage/handlers"
	"github.com/LiteyukiStudio/spage/spage/middle"
	"github.com/cloudwego/hertz/pkg/route"
)

func registerUserGroup(group *route.RouterGroup, groupWithoutAuth *route.RouterGroup) {
	userGroupWithoutAuth := groupWithoutAuth.Group("/user")
	{
		userGroupWithoutAuth.GET("/captcha", handlers.User.GetCaptchaConfig)
		userGroupWithoutAuth.POST("/logout", handlers.User.Logout)
		userGroupWithoutAuth.GET("/oidc/config", handlers.Oidc.ListOidcConfig)
		userGroupWithoutAuth.GET("/oidc/login/:name", handlers.Oidc.LoginOidcConfig)

		userGroupWithoutAuthNeedCaptcha := userGroupWithoutAuth.Group("").Use(middle.Captcha.UseCaptcha())
		{
			userGroupWithoutAuthNeedCaptcha.POST("/register", middle.Auth.NeedVerifyEmail(), handlers.User.Register)
			userGroupWithoutAuthNeedCaptcha.POST("/login", handlers.User.Login)
			userGroupWithoutAuthNeedCaptcha.POST("/verify-code", handlers.User.SendVerifyCode) // 发送邮件验证码
		}
	}

	userInfoGroup := group.Group("/user-info")
	{
		userInfoGroup.PUT("", handlers.User.UpdateUser)               // 更新用户信息 Update user info
		userInfoGroup.GET("", handlers.User.GetUser)                  // 获取用户信息 Get user info
		userInfoGroup.GET("/:id", handlers.User.GetUser)              // 获取用户信息 Get user info
		userInfoGroup.GET("/:id/projects", handlers.User.GetProjects) // 获取用户项目 Get user projects
		userInfoGroup.GET("/:id/orgs", handlers.User.GetOrgs)         // 获取用户组织 Get user orgs
	}

	userTokenGroup := group.Group("/user-token")
	{
		userTokenGroup.POST("", handlers.User.CreateApiToken)
		userTokenGroup.GET("/list", handlers.User.ListApiToken)
		userTokenGroup.DELETE("/:id", handlers.User.RevokeApiToken)
	}
}
