package dto

import "biz-auto-api/pkg/dto"

type GetListUserReq struct {
	dto.PaginationReq
	Status   int64  `json:"status" form:"status" validate:"oneof=0 1 2"`
	DeptId   int64  `json:"deptId" form:"deptId" validate:"gte=0"`
	Username string `json:"username" form:"username"`
	Search   string `json:"search" form:"search"`
	NickName string `json:"nickName" form:"nickName"`
	Email    string `json:"email" form:"email"`
	Tel      string `json:"tel" form:"tel"`
	Source   int64  `json:"source" form:"source" validate:"oneof=0 1 2"`
	RoleId   int64  `json:"roleId" form:"roleId" validate:"gte=0"`
}
type UserBaseInfo struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	NickName string `json:"nickName"`
}

type User struct {
	Id       int64   `json:"id"`
	Username string  `json:"username"`
	NickName string  `json:"nickName"`
	Email    string  `json:"email"`
	Tel      string  `json:"tel"`
	Status   int64   `json:"status"`
	RoleIds  []int64 `json:"roleIds"`
	Roles    []*Role `json:"roles"`
	DeptId   int64   `json:"deptId"`
	Dept     *Dept   `json:"dept"`
	Source   int64   `json:"source"`
}

type GetListUserRes struct {
	Items []*User `json:"items"`
	Total int64   `json:"total"`
}

type UpdateUserReq struct {
	Id       int64   `json:"id" form:"id" uri:"id" validate:"required,gte=1"`
	Username *string `json:"username" validate:"omitnil,min=4"`
	NickName *string `json:"nickName" validate:"omitnil,min=2"`
	Password *string `json:"password" validate:"omitnil,eq=|min=6"`
	Email    *string `json:"email" validate:"omitnil,eq=|email"`
	Tel      *string `json:"tel" validate:"omitnil,eq=|min=11,max=14"`
	Status   *int64  `json:"status" validate:"omitnil,oneof=1 2"`
	RoleIds  []int64 `json:"roleIds" validate:"omitnil,min=1,dive,gt=0"`
	DeptId   *int64  `json:"deptId" validate:"omitnil,gt=0"`
}
type UpdateUserRes struct {
}

type AddUserReq struct {
	Username string  `json:"username" validate:"required,min=4"`
	NickName string  `json:"nickName" validate:"required,min=2"`
	Password string  `json:"password" validate:"required,min=6"`
	Email    string  `json:"email" validate:"required,email"`
	Tel      string  `json:"tel" validate:"eq=|min=11,max=14"`
	Status   int64   `json:"status" validate:"required,oneof=1 2"`
	RoleIds  []int64 `json:"roleIds" validate:"required,min=1,dive,gt=0"`
	DeptId   int64   `json:"deptId" validate:"required,gt=0"`
}
type AddUserRes struct{}
type DeleteUserReq struct {
	Id int64 `json:"id" form:"id" uri:"id" validate:"required,gte=1"`
}
type DeleteUserRes struct{}
