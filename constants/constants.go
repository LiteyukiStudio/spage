package constants

const (
	GlobalRoleAdmin   = "admin"        // 管理员
	GlobalRoleUser    = "user"         // 普通用户
	OrgRoleOwner      = "owner"        // 组织所有者
	OrgRoleMember     = "member"       // 组织成员
	ProjectRoleOwner  = "owner"        // 项目所有者
	ProjectRoleMember = "member"       // 项目成员
	FlagSystemAdmin   = "system_admin" // 系统管理员标志

	// DB

	PreloadFieldMembers = "Members"
	PreloadFieldOwners  = "Owners"
	PreloadFieldProject = "Project"
	PreloadFieldFile    = "File"

	CaptchaTypeDisable   = "disable"     // 禁用验证码
	CaptchaTypeTurnstile = "turnstile"   //  turnstile
	CaptchaTypeReCaptcha = "recaptcha"   // Google reCaptcha
	CaptchaTypeHCaptcha  = "hcaptcha"    // HCaptcha
	CaptchaDevPasscode   = "dev-captcha" // 开发者验证码

	ModeDev  = "dev"  // 开发者模式
	ModeProd = "prod" // 生产模式

	OwnerTypeUser = "user"         // 个人用户
	OwnerTypeOrg  = "organization" // 组织用户

	FileDriverLocal      = "local"
	FileDriverWebdav     = "webdav"
	FileDriverS3         = "s3"
	WebDavPolicyProxy    = "proxy"
	WebDavPolicyRedirect = "redirect"

	DomainVerifyPolicyNone   = "none"   // 不验证域名 - 适合私域实例
	DomainVerifyPolicyLoose  = "loose"  // 宽松验证 - 主域名验证后子域名无需验证
	DomainVerifyPolicyStrict = "strict" // 严格验证 - 所有域名都需要验证

	KVPrefixEmailVerifyCode = "email-verify:"  // KV存储前缀 - 用于存储验证信息
	KVValueEmailVerified    = "email-verified" // KV存储值 - 用于标记已验证
	KVValueEmailPending     = "email-pending"  // KV存储值 - 用于标记验证中
)
