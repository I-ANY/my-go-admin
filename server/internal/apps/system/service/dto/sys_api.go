package dto

import "biz-auto-api/pkg/dto"

type AddOrUpdateApiReq struct {
	Apis []*SysApi `json:"apis"`
}

type SysApi struct {
	Id      int64  `json:"id"`
	Path    string `json:"path"`
	Method  string `json:"method"`
	Handler string `json:"handler"`
}

type GetApiListReq struct {
	dto.PaginationReq `json:"-"`
	Search            string `json:"search"`
}

type GetApiListRes struct {
	Items []*SysApi `json:"apis"`
	Total int64     `json:"total"`
}
