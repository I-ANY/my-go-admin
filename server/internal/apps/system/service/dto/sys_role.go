package dto

import "biz-auto-api/pkg/dto"

type Role struct {
	Id       int64   `json:"id" `
	Name     string  `json:"name"`
	Identify string  `json:"identify"`
	Status   int64   `json:"status"`
	Remark   string  `json:"remark"`
	MenuIds  []int64 `json:"menuIds"`
}
type GetRoleListRes struct {
	Items []*Role `json:"items"`
	Total int64   `json:"total"`
}

type GetRoleListReq struct {
	dto.PaginationReq
	Status int64  `json:"status" form:"status" validate:"oneof=0 1 2"`
	Search string `json:"search" form:"search"`
}

type AddRoleReq struct {
	Name     string  `json:"name" validate:"required"`
	Identify string  `json:"identify" validate:"required,min=5"`
	Status   int64   `json:"status" validate:"required,oneof=1 2"`
	Remark   string  `json:"remark"`
	MenuIds  []int64 `json:"menuIds" validate:"dive,gt=0"`
}
type AddRoleRes struct {
}
type UpdateRoleReq struct {
	Id      int64   `json:"id" form:"id" uri:"id" validate:"required,gt=0"`
	Name    *string `json:"name" validate:"required"`
	Status  *int64  `json:"status" validate:"required,oneof=1 2"`
	Remark  *string `json:"remark"`
	MenuIds []int64 `json:"menuIds" validate:"dive,gt=0"`
	//Identify string  `json:"identify" validate:"required"`

}
type UpdateRoleRes struct {
}
type DeleteRoleReq struct {
	Id int64 `json:"id" form:"id" uri:"id" validate:"required,gt=0"`
}
type DeleteRoleRes struct{}
