package dto

type SysBaseUserLoginReq struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
type SysBaseUserLoginRes struct {
	Success  bool   `json:"-"`
	Message  string `json:"-"`
	Token    string `json:"token"`
	ExpireAt int64  `json:"expireAt"`
}

type SysBaseGetUserInfoReq struct{}
type SysBaseGetUserInfoRes struct {
	Id       int64    `json:"id"`
	Username string   `json:"username"`
	NickName string   `json:"nickname"`
	Email    string   `json:"email"`
	Tel      string   `json:"tel"`
	Roles    []string `json:"roles"`
	Avatar   string   `json:"avatar"`
	Source   int64    `json:"source"`
}

type SysBaseGetUserPermCodeReq struct {
}

type SysBaseGetUserMenuReq struct {
}

type RefreshTokenReq struct {
}

type SysBaseGetStarPortalLoginUrlReq struct {
	RedirectUri string `json:"redirectUri" uri:"redirectUri" form:"redirectUri" validate:"required"`
}
type SysBaseGetStarPortalLoginUrlRes struct {
	LoginUrl string `json:"loginUrl"`
}

type SysBaseStarPortalLoginReq struct {
	Code        string `json:"code" validate:"required"`
	State       string `json:"state" validate:"required"`
	RedirectUri string `json:"redirectUri" validate:"required"`
}
