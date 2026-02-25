package dto

type MenuMeta struct {
	Title               string `json:"title" uri:"title" form:"title"`
	IgnoreKeepAlive     bool   `json:"ignoreKeepAlive" uri:"ignoreKeepAlive" form:"ignoreKeepAlive"`
	Icon                string `json:"icon,omitempty" uri:"icon" form:"icon"`
	HideChildrenInMenu  bool   `json:"hideChildrenInMenu" uri:"hideChildrenInMenu" form:"hideChildrenInMenu"`
	HideMenu            bool   `json:"hideMenu" uri:"hideMenu" form:"hideMenu"`
	OrderNo             int64  `json:"orderNo" uri:"orderNo" form:"orderNo"`
	IgnoreRoute         bool   `json:"ignoreRoute" uri:"ignoreRoute" form:"ignoreRoute"`
	HidePathForChildren bool   `json:"hidePathForChildren" uri:"hidePathForChildren" form:"hidePathForChildren"`
}
type Menu struct {
	Id         int64     `json:"id" uri:"id" form:"id"`
	Path       string    `json:"path" uri:"path" form:"path"`
	Name       string    `json:"name" uri:"name" form:"name"`
	Component  string    `json:"component" uri:"component" form:"component"`
	Redirect   string    `json:"redirect,omitempty" uri:"redirect" form:"redirect"`
	Meta       *MenuMeta `json:"meta" uri:"meta" form:"meta"`
	Children   []*Menu   `json:"children" uri:"children" form:"children"`
	ParentMenu int64     `json:"parentMenu"`
	Permission string    `json:"permission" uri:"permission" form:"permission"`
	Type       int64     `json:"type" uri:"type" form:"type"`
	MenuTitle  string    `json:"menuTitle"`
	ApiIds     []int64   `json:"apiIds,omitempty"`
	Status     int64     `json:"status" uri:"status" form:"status"`
}

type GetMenuTreeReq struct {
}

type GetMenuTreeRes struct {
	Items []*Menu `json:"menus"`
	Total int64   `json:"total"`
}

type UpdateMenuReq struct {
	Id                  int64   `json:"id" uri:"id" form:"id"  validate:"required,gt=0"`
	Path                *string `json:"path" uri:"path" form:"path"`
	Name                *string `json:"name" uri:"name" form:"name" validate:"required"`
	Component           *string `json:"component" uri:"component" form:"component"`
	Redirect            *string `json:"redirect,omitempty" uri:"redirect" form:"redirect"`
	ParentMenu          *int64  `json:"parentMenu"`
	Permission          *string `json:"permission" uri:"permission" form:"permission"`
	Type                *int64  `json:"type" uri:"type" form:"type" validate:"required"`
	MenuTitle           string  `json:"menuTitle"`
	ApiIds              []int64 `json:"apiIds,omitempty"`
	Title               *string `json:"title" uri:"title" form:"title" validate:"required"`
	IgnoreKeepAlive     *bool   `json:"ignoreKeepAlive" uri:"ignoreKeepAlive" form:"ignoreKeepAlive"`
	Icon                *string `json:"icon,omitempty" uri:"icon" form:"icon"`
	HideChildrenInMenu  *bool   `json:"hideChildrenInMenu" uri:"hideChildrenInMenu" form:"hideChildrenInMenu"`
	HideMenu            *bool   `json:"hideMenu" uri:"hideMenu" form:"hideMenu"`
	OrderNo             *int64  `json:"orderNo" uri:"orderNo" form:"orderNo" validate:"required"`
	IgnoreRoute         *bool   `json:"ignoreRoute" uri:"ignoreRoute" form:"ignoreRoute"`
	HidePathForChildren *bool   `json:"hidePathForChildren" uri:"hidePathForChildren" form:"hidePathForChildren"`
	Status              *int64  `json:"status" uri:"status" form:"status"`
}
type UpdateMenuRes struct {
}
type AddMenuReq struct {
	Path                *string `json:"path" uri:"path" form:"path"`
	Name                *string `json:"name" uri:"name" form:"name" validate:"required"`
	Component           *string `json:"component" uri:"component" form:"component"`
	Redirect            *string `json:"redirect,omitempty" uri:"redirect" form:"redirect"`
	ParentMenu          *int64  `json:"parentMenu"`
	Permission          *string `json:"permission" uri:"permission" form:"permission"`
	Type                *int64  `json:"type" uri:"type" form:"type" validate:"required"`
	MenuTitle           string  `json:"menuTitle"`
	ApiIds              []int64 `json:"apiIds,omitempty"`
	Title               *string `json:"title" uri:"title" form:"title" validate:"required"`
	IgnoreKeepAlive     *bool   `json:"ignoreKeepAlive" uri:"ignoreKeepAlive" form:"ignoreKeepAlive"`
	Icon                *string `json:"icon,omitempty" uri:"icon" form:"icon"`
	HideChildrenInMenu  *bool   `json:"hideChildrenInMenu" uri:"hideChildrenInMenu" form:"hideChildrenInMenu"`
	HideMenu            *bool   `json:"hideMenu" uri:"hideMenu" form:"hideMenu"`
	OrderNo             *int64  `json:"orderNo" uri:"orderNo" form:"orderNo" validate:"required"`
	IgnoreRoute         *bool   `json:"ignoreRoute" uri:"ignoreRoute" form:"ignoreRoute"`
	HidePathForChildren *bool   `json:"hidePathForChildren" uri:"hidePathForChildren" form:"hidePathForChildren"`
	Status              *int64  `json:"status" uri:"status" form:"status"`
}
type AddMenuRes struct {
}

type DeleteMenuReq struct {
	Id int64 `json:"id" uri:"id" form:"id"  validate:"required,gt=0"`
}
type DeleteMenuRes struct {
	Success bool   `json:"-"`
	Message string `json:"-"`
}
