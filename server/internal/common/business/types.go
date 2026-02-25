package business

import "biz-auto-api/pkg/dto"

type GetCategoryListReq struct {
	dto.PaginationReq
	Name  *string  `json:"name" uri:"name" form:"name"`
	Code  *string  `json:"code" uri:"code" form:"code"`
	Names []string `json:"names" uri:"names[]" form:"names[]"`
	//BizType *int64  `json:"bizType"  uri:"bizType" form:"bizType"`
	Status   *int64 `json:"status"  uri:"status" form:"status"`
	Paginate *bool  `json:"paginate" uri:"paginate" form:"paginate"`
}

type CategoryItem struct {
	Id            *int64             `json:"id"`
	Name          *string            `json:"name"`
	Code          *string            `json:"code"`
	Status        *int64             `json:"status"`
	Subcategories []*SubcategoryItem `json:"subcategories,omitempty"`
}
type GetCategoryListRes struct {
	Items []*CategoryItem `json:"items"`
	Total int64           `json:"total"`
}

type GetSubcategoryListReq struct {
	dto.PaginationReq
	Paginate     *bool    `json:"paginate" uri:"paginate" form:"paginate"`
	Name         *string  `json:"name" uri:"name" form:"name"`
	Names        []string `json:"names" uri:"names[]" form:"names[]"`
	Status       *int64   `json:"status"  uri:"status" form:"status"`
	CategoryIds  []int64  `json:"categoryIds" uri:"categoryId[]" form:"categoryId[]" validate:"omitempty,dive,gt=0"`
	CategoryCode *string  `json:"categoryCode" uri:"categoryCode" form:"categoryCode"  `
	IDCType      *int64   `json:"idcType"  uri:"idcType" form:"idcType" validate:"omitempty,oneof=1 2"`
	AddCategory  *bool    `json:"addCategory" uri:"addCategory" form:"addCategory"`
}

type SubcategoryItem struct {
	Id       *int64        `json:"id"`
	Name     *string       `json:"name"`
	Status   *int64        `json:"status"`
	Category *CategoryItem `json:"category,omitempty"`
}
type GetSubcategoryListRes struct {
	Items []*SubcategoryItem `json:"items"`
	Total int64              `json:"total"`
}
