package dto

import "biz-auto-api/pkg/dto"

type GetResourceListReq struct {
	dto.PaginationReq
	Name     *string `json:"name" uri:"name" form:"name"`
	Identify *string `json:"identify" uri:"identify" form:"identify"`
	Table    *string `json:"table" uri:"table" form:"table"`
}

type PermissionType struct {
	Code *string `json:"code" uri:"code" form:"code" validate:"required,min=4"`
	Name *string `json:"name" uri:"name" form:"name" validate:"required,min=2"`
}
type ResourceItem struct {
	Id              *int64                  `json:"id"`
	Name            *string                 `json:"name"`
	Identify        *string                 `json:"identify"`
	Table           *string                 `json:"table"`
	Filter          *string                 `json:"filter"`
	PermissionTypes []*PermissionType       `json:"permissionTypes"`
	Fields          []*SysResourceTypeField `json:"fields"`
	UpdatedAt       *string                 `json:"updatedAt"`
	Sort            *int64                  `json:"sort"`
}

type SysResourceTypeField struct {
	Id            *int64  `json:"id" uri:"id" form:"id"`
	FieldName     *string `json:"fieldName" uri:"fieldName" form:"fieldName" validate:"required"`
	ColumnName    *string `json:"columnName" uri:"columnName" form:"columnName" validate:"required,min=2"`
	IsDict        *int64  `json:"isDict" uri:"isDict" form:"isDict" validate:"required,oneof=0 1"`
	ShowWithTag   *int64  `json:"showWithTag" uri:"showWithTag" form:"showWithTag" validate:"omitnil,oneof=0 1"`
	DictKey       *string `json:"dictKey" uri:"dictKey" form:"dictKey" validate:"omitnil,required_if=IsDict 1,min=1"`
	SupportFilter *int64  `json:"supportFilter" uri:"supportFilter" form:"supportFilter" validate:"oneof=0 1"`
	Sort          *int64  `json:"sort" uri:"sort" form:"sort"`
}
type GetResourceListRes struct {
	Items []*ResourceItem `json:"items"`
	Total int64           `json:"total"`
}

type GetResourceTableListReq struct {
	dto.PaginationReq
	Name *string `json:"name" uri:"name" form:"name"`
}
type TableItem struct {
	Name *string `json:"name"`
}
type GetResourceTableFieldListRes struct {
	Items []*TableItem `json:"items"`
	Total int64        `json:"total"`
}

type GetResourceTableFieldReq struct {
	TableName *string `json:"tableName" uri:"tableName" form:"tableName" validate:"required,min=1"`
}
type GetResourceTableFieldRes struct {
	Fields []string `json:"fields"`
}

type AddResourceReq struct {
	Name            *string                 `json:"name" validate:"required,min=2"`
	Identify        *string                 `json:"identify" validate:"required,min=5"`
	Table           *string                 `json:"table" validate:"required,min=1"`
	Filter          *string                 `json:"filter"`
	PermissionTypes []*PermissionType       `json:"permissionTypes" validate:"min=1,dive"`
	Fields          []*SysResourceTypeField `json:"fields" validate:"min=1,dive"`
	Sort            *int64                  `json:"sort" uri:"sort" form:"sort" validate:"required"`
}

type AddResourceRes struct {
}
type UpdateResourceReq struct {
	Id              *int64                  `json:"id" uri:"id" form:"id" validate:"required,gt=0"`
	Name            *string                 `json:"name" validate:"required,min=2"`
	Identify        *string                 `json:"identify" validate:"required,min=5"`
	Table           *string                 `json:"table" validate:"required,min=1"`
	Filter          *string                 `json:"filter"`
	PermissionTypes []*PermissionType       `json:"permissionTypes" validate:"min=1,dive"`
	Fields          []*SysResourceTypeField `json:"fields" validate:"min=1,dive"`
	Sort            *int64                  `json:"sort" uri:"sort" form:"sort" validate:"required"`
}
type UpdateResourceRes struct {
}

type DeleteResourceReq struct {
	Id *int64 `json:"id" uri:"id" form:"id" validate:"required,gt=0"`
}
type DeleteResourceRes struct {
}

type GetResourceViewFormSchemasReq struct {
	Id *int64 `json:"id" uri:"id" form:"id" validate:"required,gt=0"`
}
type FormSchema struct {
	Label     *string `json:"label,omitempty"`
	Field     *string `json:"field,omitempty"`
	Component *string `json:"component,omitempty"`
	ColProps  *struct {
		Span *int64 `json:"span,omitempty"`
	} `json:"colProps,omitempty"`
	ComponentProps map[string]interface{} `json:"componentProps,omitempty"`
}
type GetResourceViewFormSchemasRes struct {
	Schemas []*FormSchema `json:"schemas"`
}

type GetResourceViewTableColumnsReq struct {
	Id *int64 `json:"id" uri:"id" form:"id" validate:"required,gt=0"`
}
type TableColumn struct {
	Title     *string `json:"title,omitempty"`
	DataIndex *string `json:"dataIndex,omitempty"`
	Width     *int64  `json:"width,omitempty"`
	Resizable *bool   `json:"resizable,omitempty"`
}
type GetResourceViewTableColumnsRes struct {
	Columns        []*TableColumn    `json:"columns"`
	ShowEnumFields map[string]string `json:"showEnumFields"`
	ShowTagFields  map[string]string `json:"showTagFields"`
}
type GetResourceDetailListExtraParams map[string]interface{}
type GetResourceDetailListReq struct {
	Id *int64 `json:"id" uri:"id" form:"id" validate:"required,gt=0"`
	dto.PaginationReq
	ExtraParams GetResourceDetailListExtraParams `json:"-" form:"-" uri:"-"`
}

type ResourceDetailItem map[string]interface{}
type GetResourceDetailListRes struct {
	Items []ResourceDetailItem `json:"items"`
	Total int64                `json:"total"`
}

type GetRoleResourceInfoReq struct {
	RoleId         *int64  `json:"roleId" uri:"roleId" form:"roleId" validate:"required,gt=0"`
	ResourceTypeId *int64  `json:"resourceTypeId" uri:"resourceTypeId" form:"resourceTypeId" validate:"required_without=Identify,omitnil,gt=0"`
	Identify       *string `json:"identify" uri:"identify" form:"identify" validate:"required_without=ResourceTypeId,omitnil,min=1"`
}
type GetRoleResourceInfoRes struct {
	PermissionTypes                  []*PermissionType `json:"permissionTypes"`
	AuthedAllResourcePermissionTypes []string          `json:"authedAllResourcePermissionTypes"` // 授权所有资源 0:否 1:是
}

type GetRoleResourceDetailListReq struct {
	RoleId         *int64 `json:"roleId" uri:"roleId" form:"roleId" validate:"required,gt=0"`
	ResourceTypeId *int64 `json:"resourceTypeId" uri:"resourceTypeId" form:"resourceTypeId" validate:"required,gt=0"`
	dto.PaginationReq
	ExtraParams GetResourceDetailListExtraParams `json:"-" form:"-" uri:"-"`
}
type GetRoleResourceDetailListRes GetResourceDetailListRes

type UpdateRoleResourceReq struct {
	RoleId                           *int64                   `json:"roleId" uri:"roleId" form:"roleId" validate:"required,gt=0"`
	ResourceTypeId                   *int64                   `json:"resourceTypeId" uri:"resourceTypeId" form:"resourceTypeId" validate:"required,gt=0"`
	AuthedAllResourcePermissionTypes []string                 `json:"authedAllResourcePermissionTypes" uri:"authedAllResourcePermissionTypes" form:"authedAllResourcePermissionTypes" validate:"dive,min=1"`
	ChangedPermissionTypes           []*ChangedPermissionType `json:"changedPermissionTypes" validate:"dive"`
}
type ChangedPermissionType struct {
	ResourceId      *int64   `json:"resourceId" validate:"required,gt=0"`
	PermissionTypes []string `json:"permissionTypes"`
}

type UpdateRoleResourceRes struct {
}

type GetRoleAuthedResourceReq struct {
	RoleId         *int64  `json:"roleId" uri:"roleId" form:"roleId" validate:"required,gt=0"`
	ResourceTypeId *int64  `json:"resourceTypeId" uri:"resourceTypeId" form:"resourceTypeId" validate:"required_without=Identify,omitnil,gt=0"`
	Identify       *string `json:"identify" uri:"identify" form:"identify" validate:"required_without=ResourceTypeId,omitnil,min=1"`
	PermissionCode *string `json:"permissionCode" uri:"permissionCode" form:"permissionCode" validate:"required,min=1"`
}

type GetRoleAuthedResourceRes struct {
	AuthAllResource *int64  `json:"authAllResource"`
	ResourceIds     []int64 `json:"resourceIds"`
}
type Category struct {
	Id            *int64         `json:"id"`
	Name          *string        `json:"name"`
	Subcategories []*Subcategory `json:"subcategories"`
}
type Subcategory struct {
	Id   *int64  `json:"id"`
	Name *string `json:"name"`
}
type GetBusinessResourceReq struct {
}

type GetBusinessResourceRes struct {
	Data []*Category `json:"data"`
}

const (
	AuthAllResourceYes = 1
	AuthAllResourceNo  = 0
)

type RoleResourceAuthReq struct {
	RoleId          *int64   `json:"roleId" uri:"roleId" form:"roleId" validate:"required,gt=0"`
	ResourceTypeId  *int64   `json:"resourceTypeId" uri:"resourceTypeId" form:"resourceTypeId" validate:"required_without=Identify,omitnil,gt=0"`
	Identify        *string  `json:"identify" uri:"identify" form:"identify" validate:"required_without=ResourceTypeId,omitnil,min=1"`
	PermissionCode  *string  `json:"permissionCode" uri:"permissionCode" form:"permissionCode" validate:"required,min=1"`
	ResourceIds     []string `json:"resourceIds" uri:"resourceIds" form:"resourceIds" validate:"dive,min=1"`
	AuthAllResource *int64   `json:"authAllResource" validate:"required,oneof=0 1"`
}
type AuthSubcategoryPermissionRes struct {
}
