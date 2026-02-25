package business

import (
	"biz-auto-api/internal/models"
	"biz-auto-api/pkg/tools"
	"context"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Queryer struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewQueryer(db *gorm.DB, rdb *redis.Client) *Queryer {
	return &Queryer{
		db:  db,
		rdb: rdb,
	}
}

// GetAuthorizedCategoryList 查询用户已授权的业务大类列表
func (q *Queryer) GetAuthorizedCategoryList(ctx context.Context, userId int64, req *GetCategoryListReq) (*GetCategoryListRes, error) {
	var (
		total                   int64
		categories              = make([]*models.BusinessCategory, 0)
		result                  = &GetCategoryListRes{}
		err                     error
		authorizedSubcategories = make([]string, 0) // 已授权的业务子类ID
	)

	//TODO 校验用户权限，并只查询出用户有权限的业务小类，再通过小类去查询大类，目前没有接入权限系统，所以暂时不校验用户权限，直接查询出所有业务小类
	subQuery := q.db.Model(&models.BusinessSubcategory{}).Scopes(func(tx *gorm.DB) *gorm.DB {
		// 如果包含 * 则说明所有权限都有，则不进行权限过滤
		if len(authorizedSubcategories) > 0 && !tools.InSlice("*", authorizedSubcategories) {
			return tx.Where("id IN ?", authorizedSubcategories)
		}
		if len(req.Names) > 0 {
			tx = tx.Where("name in (?)", req.Names)
		}
		tx = tx.Where("status = ?", models.BusinessEnable)
		return tx
	}).Distinct("category_id").Select("category_id")
	tmpdb := q.db.Model(&models.BusinessCategory{}).Joins("join (?) as s on business_category.id = s.category_id", subQuery).
		Scopes(func(tx *gorm.DB) *gorm.DB {
			if req.Code != nil {
				tx = tx.Where("business_category.code = ?", req.Code)
			}
			if req.Name != nil {
				tx = tx.Where("business_category.name like  ?", tools.FuzzyQuery(*req.Name))
			}
			if req.Status != nil {
				tx = tx.Where("business_category.status = ?", *req.Status)
			}

			return tx
		}).Distinct().Session(&gorm.Session{}).Preload("Subcategories")
	// 分页
	if tools.ToValue(req.Paginate) {
		err = tmpdb.Order("business_category.name ASC").Scopes(req.MakePagination()).Find(&categories).Error
		if err != nil {
			err = errors.Wrap(err, "list category failed")
			return nil, err
		}
		err = tmpdb.Count(&total).Error
		if err != nil {
			err = errors.Wrap(err, "list category count failed")
			return nil, err
		}
	} else {
		err = tmpdb.Order("business_category.name ASC").Find(&categories).Error
		if err != nil {
			err = errors.Wrap(err, "list category failed")
			return nil, err
		}
		total = int64(len(categories))
	}
	items := make([]*CategoryItem, len(categories))
	err = copier.Copy(&items, categories)
	if err != nil {
		err = errors.Wrap(err, "copy category list failed")
		return nil, err
	}
	result.Total = total
	result.Items = items
	return result, nil
}

// GetAuthorizedSubcategoryList 查询用户已授权的业务小类列表
func (q *Queryer) GetAuthorizedSubcategoryList(ctx context.Context, userId int64, req *GetSubcategoryListReq) (*GetSubcategoryListRes, error) {
	var (
		total                   int64
		subcategories           = make([]*models.BusinessSubcategory, 0)
		result                  = &GetSubcategoryListRes{}
		err                     error
		authorizedSubcategories = make([]string, 0)
	)

	commonQuery := func() *gorm.DB {
		d := q.db.Model(&models.BusinessSubcategory{}).Scopes(func(tx *gorm.DB) *gorm.DB {
			if req.Name != nil {
				tx = tx.Where("name like ?", tools.FuzzyQuery(*req.Name))
			}
			if len(req.Names) > 0 {
				tx = tx.Where("name in (?)", req.Names)
			}
			if req.Status != nil {
				tx = tx.Where("status = ?", *req.Status)
			}
			if len(req.CategoryIds) > 0 {
				tx = tx.Where("category_id in ?", req.CategoryIds)
			}
			// 如果包含 * 则说明所有权限都有，则不进行权限过滤
			if len(authorizedSubcategories) > 0 && !tools.InSlice("*", authorizedSubcategories) {
				tx = tx.Where("id in (?)", authorizedSubcategories)
			}
			if req.IDCType != nil {
				if *req.IDCType == models.IDCTypeSelfBuild {
					tx.Where("business_subcategory.name like ?", "%_MF")
				} else {
					tx.Where("business_subcategory.name not like ?", "%_MF")
				}
			}
			if tools.ToValue(req.AddCategory) {
				tx.Preload("Category")
			}
			return tx
		})
		return d
	}

	joinQuery := func() *gorm.DB {
		subQuery := q.db.Model(&models.BusinessCategory{}).Scopes(func(tx *gorm.DB) *gorm.DB {
			tx = tx.Where("status = ?", models.BusinessEnable)
			if req.CategoryCode != nil {
				tx = tx.Where("code = ?", *req.CategoryCode)
			}
			if len(req.CategoryIds) > 0 {
				tx = tx.Where("c.id in ?", req.CategoryIds)
			}
			return tx
		}).Distinct("id")

		d := q.db.Model(&models.BusinessSubcategory{}).Joins("join (?) as c on business_subcategory.category_id = c.id", subQuery).Scopes(func(tx *gorm.DB) *gorm.DB {
			if req.Name != nil {
				tx = tx.Where("business_subcategory.name like ?", tools.FuzzyQuery(*req.Name))
			}
			if len(req.Names) > 0 {
				tx = tx.Where("business_subcategory.name in (?)", req.Names)
			}
			if req.Status != nil {
				tx = tx.Where("business_subcategory.status = ?", *req.Status)
			}
			if len(req.CategoryIds) > 0 {
				tx = tx.Where("category_id in ?", req.CategoryIds)
			}
			if len(authorizedSubcategories) > 0 {
				tx = tx.Where("business_subcategory.name in (?)", authorizedSubcategories)
			}
			if req.IDCType != nil {
				if *req.IDCType == models.IDCTypeSelfBuild {
					tx.Where("business_subcategory.name like ?", "%_MF")
				} else {
					tx.Where("business_subcategory.name not like ?", "%_MF")
				}
			}
			if tools.ToValue(req.AddCategory) {
				tx.Preload("Category")
			}
			return tx
		})
		return d
	}

	//TODO 校验用户权限，并只查询出用户有权限的业务小类，目前没有接入权限系统，所以暂时不校验用户权限，直接查询出所有业务小类
	var tmpDb *gorm.DB
	if req.CategoryCode != nil {
		tmpDb = joinQuery().Session(&gorm.Session{})
	} else {
		tmpDb = commonQuery().Session(&gorm.Session{})
	}

	// 分页
	if tools.ToValue(req.Paginate) {
		err = tmpDb.Order("business_subcategory.name ASC").Scopes(req.MakePagination()).Find(&subcategories).Error
		if err != nil {
			return nil, errors.Wrap(err, "query subcategory failed")
		}
		err = tmpDb.Count(&total).Error
		if err != nil {
			return nil, errors.Wrap(err, "count subcategory failed")
		}
	} else {
		err = tmpDb.Order("business_subcategory.name ASC").Find(&subcategories).Error
		if err != nil {
			return nil, errors.Wrap(err, "query subcategory failed")
		}
		total = int64(len(subcategories))
	}

	items := make([]*SubcategoryItem, 0, len(subcategories))
	err = copier.Copy(&items, subcategories)
	if err != nil {
		return nil, errors.Wrap(err, "copy subcategory list failed")
	}
	result.Items = items
	result.Total = total

	return result, nil
}

// GetAllAuthorizedSubcategoryNames 获取用户授权的所有子类的名称
func (q *Queryer) GetAllAuthorizedSubcategoryNames(ctx context.Context, userId int64, req *GetSubcategoryListReq) ([]string, error) {
	if req == nil {
		req = &GetSubcategoryListReq{}
	}
	// 显示设置分页为false
	req.Paginate = tools.ToPointer(false)
	subcategory, err := q.GetAuthorizedSubcategoryList(ctx, userId, req)
	if err != nil {
		return nil, err
	}
	names := tools.GetSlice(subcategory.Items, func(c *SubcategoryItem) string {
		return tools.ToValue(c.Name)
	})
	return names, nil
}

// GetAllAuthorizedSubcategoryNameWithWildcard 获取用户授权的所有子类的名称，如果存在通配符，则返回所有子类的名称
func (q *Queryer) GetAllAuthorizedSubcategoryNameWithWildcard(ctx context.Context, userId int64) ([]string, error) {
	var (
		names                   = make([]string, 0)
		authorizedSubcategories = make([]string, 0)
	)
	// TODO 获取用户有权限的业务子类ID
	// authorizedSubcategories=
	// 查看用户拥有的权限中包含*，包含则返回*并设置key
	if tools.InSlice("*", authorizedSubcategories) {
		names = []string{"*"}
		return names, nil
	}
	return q.GetAllAuthorizedSubcategoryNames(ctx, userId, nil)
}
