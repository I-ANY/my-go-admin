package assessment

import (
	"strings"
	"sync"

	"biz-auto-api/pkg/tools"
	"encoding/json"

	"gorm.io/gorm"
)

// AssessmentBizResult 考核业务结果
type AssessmentBizResult struct {
	Bizs   []string          `json:"bizs"`   // 考核业务列表
	Groups map[string]string `json:"groups"` // 考核业务 -> 业务组映射
}

// BusinessMapping 业务映射服务（基于DB动态加载）
type BusinessMapping struct {
	db *gorm.DB

	// 缓存
	once          sync.Once
	loadErr       error
	subToCategory map[string]string   // 主线业务 -> 考核业务(大类名)
	categoryToGrp map[string]string   // 考核业务 -> 业务组名
	groupToCats   map[string][]string // 业务组 -> []考核业务
}

// NewBusinessMapping 创建业务映射服务实例
func NewBusinessMapping(db *gorm.DB) *BusinessMapping {
	return &BusinessMapping{db: db}
}

// loadMapping 从数据库加载映射关系（仅执行一次）
func (b *BusinessMapping) loadMapping() error {
	b.once.Do(func() {
		b.subToCategory = make(map[string]string)
		b.categoryToGrp = make(map[string]string)
		b.groupToCats = make(map[string][]string)

		// 查询主线业务 -> 考核业务(业务大类) + 业务组
		// 1) 查询业务大类(id->name)
		type catRow struct {
			Id   int64  `gorm:"column:id"`
			Name string `gorm:"column:name"`
		}
		var categories []catRow
		err := b.db.Table("business_category").
			Select("id, name").
			Where("deleted_at IS NULL").
			Find(&categories).Error
		if err != nil {
			b.loadErr = err
			return
		}
		catIdToName := make(map[int64]string)
		for _, c := range categories {
			if c.Name != "" {
				catIdToName[c.Id] = c.Name
			}
		}

		// 2) 查询业务组 -> category_ids（JSON数组）
		type grpRow struct {
			Name        string          `gorm:"column:name"`
			CategoryIds json.RawMessage `gorm:"column:category_ids;type:json"`
		}
		var groups []grpRow
		err = b.db.Table("business_group").
			Select("name, category_ids").
			Where("deleted_at IS NULL").
			Find(&groups).Error
		if err != nil {
			b.loadErr = err
			return
		}
		for _, g := range groups {
			if g.Name == "" {
				continue
			}
			var categoryIds []int64
			err = json.Unmarshal(g.CategoryIds, &categoryIds)
			if err != nil {
				b.loadErr = err
				return
			}
			for _, catId := range categoryIds {
				catName := catIdToName[catId]
				if catName == "" {
					continue
				}
				b.categoryToGrp[catName] = g.Name
				b.groupToCats[g.Name] = append(b.groupToCats[g.Name], catName)
			}
		}

		// 3) 查询主线业务 -> 考核业务(大类名)
		type subRow struct {
			SubName      string `gorm:"column:sub_name"`
			CategoryName string `gorm:"column:category_name"`
		}
		var subRows []subRow
		err = b.db.Table("business_subcategory AS bs").
			Select("bs.name AS sub_name, bc.name AS category_name").
			Joins("LEFT JOIN business_category bc ON bs.category_id = bc.id").
			Where("bs.deleted_at IS NULL").
			Where("bc.deleted_at IS NULL").
			Find(&subRows).Error
		if err != nil {
			b.loadErr = err
			return
		}

		for _, r := range subRows {
			if r.SubName == "" || r.CategoryName == "" {
				continue
			}
			b.subToCategory[r.SubName] = r.CategoryName
			if grp := b.categoryToGrp[r.CategoryName]; grp != "" {
				b.groupToCats[grp] = append(b.groupToCats[grp], r.CategoryName)
			}
		}

		// 去重 groupToCats
		for g, cats := range b.groupToCats {
			b.groupToCats[g] = tools.RemoveDuplication(cats)
		}
	})

	return b.loadErr
}

// GetAssessmentBusiness 获取主线业务对应的考核业务（业务大类名）和业务组
// 返回：考核业务列表和业务组映射（key=考核业务，value=业务组）
func (b *BusinessMapping) GetAssessmentBusiness(mainBiz string, deployedBizs []string) *AssessmentBizResult {
	result := &AssessmentBizResult{
		Bizs:   make([]string, 0),
		Groups: make(map[string]string),
	}

	if err := b.loadMapping(); err != nil {
		// 发生错误时回退：直接返回主线业务
		if mainBiz != "" {
			result.Bizs = []string{mainBiz}
			return result
		}
		result.Bizs = deployedBizs
		return result
	}

	// 单业务场景
	if mainBiz != "" {
		if cat, ok := b.subToCategory[strings.TrimSpace(mainBiz)]; ok && cat != "" {
			result.Bizs = []string{cat}
			if grp := b.categoryToGrp[cat]; grp != "" {
				result.Groups[cat] = grp
			}
			return result
		}
		result.Bizs = []string{mainBiz}
		return result
	}

	// 多业务场景 (专线节点多个业务或MCDN部署多个业务)
	if len(deployedBizs) > 0 {
		assessmentBizs := make([]string, 0, len(deployedBizs))
		for _, biz := range deployedBizs {
			biz = strings.TrimSpace(biz)
			if biz == "" {
				continue
			}
			if cat, ok := b.subToCategory[biz]; ok && cat != "" {
				assessmentBizs = append(assessmentBizs, cat)
				if grp := b.categoryToGrp[cat]; grp != "" {
					result.Groups[cat] = grp
				}
			} else {
				assessmentBizs = append(assessmentBizs, biz)
			}
		}
		result.Bizs = tools.RemoveDuplication(assessmentBizs)
		return result
	}

	return result
}

// GetBusinessGroup 获取考核业务对应的业务组
func (b *BusinessMapping) GetBusinessGroup(assessmentBiz string) string {
	if err := b.loadMapping(); err != nil {
		return ""
	}
	return b.categoryToGrp[strings.TrimSpace(assessmentBiz)]
}

// GetAllBusinessGroups 获取所有业务组及其包含的考核业务
func (b *BusinessMapping) GetAllBusinessGroups(groupName string) map[string][]string {
	if err := b.loadMapping(); err != nil {
		return map[string][]string{}
	}
	if groupName != "" {
		return map[string][]string{groupName: b.groupToCats[groupName]}
	}
	return b.groupToCats
}

// GetBusinessGroupsOrder 获取业务组排序顺序
func (b *BusinessMapping) GetBusinessGroupsOrder() []string {
	if err := b.loadMapping(); err != nil {
		return []string{}
	}

	// 查询数据库中的业务组排序
	type BusinessGroup struct {
		Name string `gorm:"column:name"`
	}

	var groups []BusinessGroup
	err := b.db.Table("business_group").
		Select("name").
		Where("deleted_at IS NULL").
		Order("name ASC").
		Find(&groups).Error

	if err != nil {
		return []string{}
	}
	businessGroups := make([]string, 0, len(groups))
	for _, group := range groups {
		businessGroups = append(businessGroups, group.Name)
	}

	return businessGroups
}
