package models

import "biz-auto-api/pkg/models"

const (
	RoleStatusEnable  = 1
	RoleStatusDisable = 2
)

type SysRole struct {
	models.BseId
	Name            string             `gorm:"comment:角色名称;size:255;not null"`
	Identify        string             `gorm:"唯一标识;size:50;not null;unique"`
	Status          int64              `gorm:"角色状态;default:1;size:1"`
	Remark          string             `json:"remark" gorm:"备注"`
	Menus           []*SysMenu         `gorm:"many2many:sys_role_menu;constraint:OnDelete:CASCADE"`
	Users           []*SysUser         `gorm:"many2many:sys_user_role;constraint:OnDelete:CASCADE"`
	AuthedResources []*SysRoleResource `gorm:"foreignKey:RoleId;references:Id;constraint:OnDelete:CASCADE;"`
	models.ControlBy
	models.ModelTime
}

func (SysRole) TableName() string {
	return "sys_role"
}

const (
	UserStatusEnable  = 1
	UserStatusDisable = 2

	UserSourceSystem     = 1
	UserSourceStarPortal = 2
)

type SysUser struct {
	models.BseId
	Username  string     `json:"username" gorm:"comment:用户登录名;not null;unique;size:50"`
	NickName  string     `json:"nickName" gorm:"comment:用户名;not null;size:50"`
	Password  string     `json:"password" gorm:"comment:密码;not null"`
	Email     *string    `json:"email" gorm:"comment:邮箱;size:255"`
	Tel       *string    `json:"tel" gorm:"comment:手机号;size:20"`
	Status    int64      `json:"status" gorm:"comment:账号状态，1启用，2禁用;default:1;size:1"`
	Roles     []*SysRole `gorm:"many2many:sys_user_role;constraint:OnDelete:CASCADE"`
	DeptId    int64      `json:"deptId" gorm:"comment:部门ID;index:,type:btree"`
	Dept      *SysDept   `json:"dept" gorm:"constraint:OnUpdate:SET NULL,OnDelete:SET NULL;references:Id;foreignKey:DeptId"`
	Avatar    string     `json:"avatar" gorm:"comment:头像地址;size:255;"`
	OpenId    *string    `json:"openId" gorm:"comment:第三方系统中的用户ID;size:255;index:uk_source_openid,unique,type:btree"`
	Source    int64      `json:"source" gorm:"comment:用户来源，1系统自建，2星云;default:1;size:1;index:uk_source_openid,unique,type:btree"`
	JobNumber *string    `json:"jobNumber" gorm:"comment:工号;default:null;index:uk_job_number,unique,type:btree"`
	models.ControlBy
	models.ModelTime
}

func (SysUser) TableName() string {
	return "sys_user"
}

const (
	DeptStatusEnable  = 1
	DeptStatusDisable = 2
)

type SysDept struct {
	models.BseId
	Name       string     `gorm:"comment:角色名称;size:255;not null"`
	Status     int64      `gorm:"角色状态;default:1;size:1"`
	OrderNo    int64      `gorm:"comment:部门排序"`
	Remark     string     `json:"remark" gorm:"备注"`
	ParentDept int64      `json:"parentDept" gorm:"comment:父级部门"`
	Users      []*SysUser `json:"Users" gorm:"constraint:OnUpdate:SET NULL,OnDelete:SET NULL;references:Id;foreignKey:DeptId"`
	models.ControlBy
	models.ModelTime
}

func (SysDept) TableName() string {
	return "sys_dept"
}

const (
	DictDataStatusEnable  = 1
	DictDataStatusDisable = 2
)
const (
	DictTypeStatusEnable  = 1
	DictTypeStatusDisable = 2
)

type SysDictType struct {
	models.BseId
	TypeName string         `gorm:"comment:字段类型名称;size:255;not null"`
	TypeCode string         `gorm:"comment:唯一标识;size:100;not null;unique"`
	Sort     int64          `gorm:"comment:排序;default:10"`
	Status   int64          `gorm:"comment:字典类型状态 1 启用，2禁用;default:1;size:1"`
	Remark   string         `json:"remark" gorm:"comment:备注"`
	DictData []*SysDictData `gorm:"foreignKey:DictTypeId;references:Id;constraint:OnDelete:CASCADE;"`
	models.ControlBy
	models.ModelTime
}

func (SysDictType) TableName() string {
	return "sys_dict_type"
}

type SysDictData struct {
	models.BseId
	DictTypeId int64        `gorm:"comment:字段类型Id;not null;uniqueIndex:uk_type_id_value"`
	DictType   *SysDictType `gorm:"foreignKey:DictTypeId;references:Id;constraint:OnDelete:CASCADE;"`
	DictLabel  string       `gorm:"comment:字段标签;size:255;not null"`
	DictValue  string       `gorm:"comment:字典值;size:255;uniqueIndex:uk_type_id_value"`
	Sort       int64        `gorm:"comment:排序;default:10"`
	Status     int64        `gorm:"comment:字典状态 1 启用，2禁用;default:1;size:1"`
	Remark     string       `gorm:"comment:备注"`
	Color      string       `gorm:"comment:展示颜色，用于展示时前端显示的颜色（如有需要）"`
	models.ControlBy
	models.ModelTime
}

func (SysDictData) TableName() string {
	return "sys_dict_data"
}
