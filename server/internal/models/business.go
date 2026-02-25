package models

import (
	"biz-auto-api/pkg/models"
	"encoding/json"
)

// 业务状态
const (
	// BusinessDisable 禁用
	BusinessDisable = 0
	// BusinessEnable 启用
	BusinessEnable = 1
)

// 业务类型
const (
	// IDCTypeSelfBuild 自建
	IDCTypeSelfBuild = 1
	// IDCTypeRecruitment 招募
	IDCTypeRecruitment = 2
)

const (
	// CategoryTypeIsVirtualNo 非虚拟业务组
	CategoryTypeIsVirtualNo = 0
	// CategoryTypeIsVirtualYes 虚拟业务组
	CategoryTypeIsVirtualYes = 1
)

// BusinessCategory 业务大类表
type BusinessCategory struct {
	models.BseId
	Name                 *string                `gorm:"column:name;not null;comment:业务大类名称;size:255;index:uk_name,unique"`
	Code                 *string                `gorm:"column:code;not null;comment:业务大类编码;size:255;index:uk_code,unique"`
	Status               *int64                 `gorm:"type:tinyint(1);column:status;not null;default:1;comment:状态 0-禁用，1-启用"`
	Describe             *string                `gorm:"column:describe;type:text;comment:业务大类描述"`
	IsVirtual            *int64                 `gorm:"column:is_virtual;default: 0;type:tinyint(1);not null;comment:是否虚拟业务组，1-是 0-否"`
	Subcategories        []*BusinessSubcategory `gorm:"constraint:OnUpdate:SET NULL,OnDelete:SET NULL;references:Id;foreignKey:CategoryId"`
	VirtualSubcategories []*BusinessSubcategory `gorm:"many2many:business_virtual_subcategory;foreignKey:Id;joinForeignKey:CategoryId;References:Id;joinReferences:SubcategoryId"` // 虚拟子业务
	models.ModelTime
}

func (BusinessCategory) TableName() string {
	return "business_category"
}

// BusinessSubcategory 业务小类
type BusinessSubcategory struct {
	models.BseId
	Name              *string             `gorm:"column:name;not null;comment:业务名称;size:255;index:uk_name_category,unique,priority:3"`
	Status            *int64              `gorm:"type:tinyint(1);column:status;not null;default:1;comment:状态 0-禁用，1-启用"`
	CategoryId        *int64              `gorm:"column:category_id;index:idx_category_id;index:uk_name_category,unique,priority:3;comment:业务大类ID"`
	Category          *BusinessCategory   `gorm:"constraint:OnUpdate:SET NULL,OnDelete:SET NULL;references:Id;foreignKey:CategoryId"`
	EcdnId            *int64              `gorm:"column:ecdn_id;comment:ECDN对应ID;unique"`
	VirtualCategories []*BusinessCategory `gorm:"many2many:business_virtual_subcategory;foreignKey:Id;joinForeignKey:SubcategoryId;References:Id;joinReferences:CategoryId"` // 虚拟业务大类
	models.ModelTime
}

func (BusinessSubcategory) TableName() string {
	return "business_subcategory"
}

// 业务组对应关系表
type BusinessGroup struct {
	models.BseId
	Name        *string         `gorm:"column:name;not null;comment:业务组名称;size:255;index:uk_name,unique"`
	CategoryIds json.RawMessage `gorm:"column:category_ids;type:json;comment:业务大类ID列表"`
	models.ModelTime
}

func (BusinessGroup) TableName() string {
	return "business_group"
}
