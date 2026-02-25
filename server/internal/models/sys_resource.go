package models

import (
	"biz-auto-api/pkg/models"
	"time"
)

const (
	// ResourceFieldIsDictYes 字段为字典
	ResourceFieldIsDictYes = 1
	// ResourceFieldIsDictNo 字段不为字典
	ResourceFieldIsDictNo = 0
)

const (
	// ResourceFieldShowWithTagYes 字段使用字典显示
	ResourceFieldShowWithTagYes = 1
	// ResourceFieldShowWithTagNo 字段不使用字典
	ResourceFieldShowWithTagNo = 0
)

const (
	// ResourceFieldSupportFilterYes 字段支持过滤
	ResourceFieldSupportFilterYes = 1
	// ResourceFieldSupportFilterNo 字段不支持过滤
	ResourceFieldSupportFilterNo = 0
)

type SysRoleResource struct {
	models.BseId
	RoleId               *int64    `gorm:"column:role_id;type:int(11);not null;uniqueIndex:uk_role_res_action,priority:3;comment:角色ID"` // 角色ID
	Role                 *SysRole  `gorm:"constraint:OnUpdate:SET NULL,OnDelete:SET NULL;references:Id;foreignKey:RoleId"`
	ResourceTypeIdentify *string   `gorm:"column:resource_type_identify;size:255;not null;uniqueIndex:uk_role_res_action,priority:5;comment:资源类型唯一标识"` // 资源类型
	ResourceId           *string   `gorm:"column:resource_id;size:255;not null;uniqueIndex:uk_role_res_action,priority:10;comment:资源ID，或者*"`           // 资源ID（支持通配符*）
	Action               *string   `gorm:"column:action;uniqueIndex:uk_role_res_action,priority:15;not null;comment:权限类型"`
	CreatedAt            time.Time `json:"createdAt" gorm:"comment:创建时间"`
	UpdatedAt            time.Time `json:"updatedAt" gorm:"comment:最后更新时间"`
}

func (*SysRoleResource) TableName() string {
	return "sys_role_resource"
}

type PermissionType struct {
	models.BseId
	Code *string `json:"code,omitempty"`
	Name *string `json:"name,omitempty"`
}

type SysResourceType struct {
	models.BseId
	Name            *string                 `gorm:"column:name;size:255;not null;uniqueIndex:uk_name;comment:资源名称"` // 资源名称
	Identify        *string                 `gorm:"column:identify;size:255;not null;column:identify;uniqueIndex:uk_identify;comment:资源唯一标识"`
	Table           *string                 `gorm:"column:table;size:255;not null;comment:资源对应的数据库表名"` // 授权时从这个表查询数据
	Filter          *string                 `gorm:"column:filter;type:text;comment:授权查询时过滤的SQL"`       // 过滤条件，授权时过滤数据
	PermissionTypes []*PermissionType       `gorm:"column:permission_types;type:json;serializer:json;not null;comment:权限类型"`
	Fields          []*SysResourceTypeField `gorm:"constraint:OnUpdate:SET NULL,OnDelete:SET NULL;references:Id;foreignKey:ResourceTypeId"`
	Sort            *int64                  `json:"sort" gorm:"column:sort;size:11;type:int;default:10;not null;comment:排序"`
	CreatedAt       time.Time               `json:"createdAt" gorm:"comment:创建时间"`
	UpdatedAt       time.Time               `json:"updatedAt" gorm:"comment:最后更新时间"`
	models.ControlBy
}

func (*SysResourceType) TableName() string {
	return "sys_resource_type"
}

type SysResourceTypeField struct {
	models.BseId
	ResourceTypeId *int64           `gorm:"column:resource_type_id;type:int(11);not null;index:idx_res_type_id;comment:资源类型ID"`
	ResourceType   *SysResourceType `gorm:"constraint:OnUpdate:SET NULL,OnDelete:SET NULL;references:Id;foreignKey:ResourceTypeId"`
	FieldName      *string          `gorm:"column:field_name;size:255;not null;comment:字段名称"`
	ColumnName     *string          `gorm:"column:column_name;size:255;not null;comment:资源对应表格的列名"`
	IsDict         *int64           `gorm:"column:is_dict;default:0;type:tinyint(1);not null;comment:是否使用字典,1-是，0-否"`
	ShowWithTag    *int64           `gorm:"column:show_with_tag;default:0;type:tinyint(1);not null;comment:是否使用标签显示,1-是，0-否"`
	DictKey        *string          `gorm:"column:dict_key;size:255;comment:如果这个字段使用的是枚举：展示/筛选时使用的字典key"`
	SupportFilter  *int64           `gorm:"column:support_filter;type:tinyint(1);comment:是否支持页面过滤，如果dict_key存在则使用select，否则使用输入框，1-是，0-否"`
	Sort           *int64           `gorm:"column:sort;size:11;type:int;default:1;not null;comment:字段展示顺序"`
	CreatedAt      time.Time        `json:"createdAt" gorm:"comment:创建时间"`
	UpdatedAt      time.Time        `json:"updatedAt" gorm:"comment:最后更新时间"`
}

func (*SysResourceTypeField) TableName() string {
	return "sys_resource_type_field"
}
