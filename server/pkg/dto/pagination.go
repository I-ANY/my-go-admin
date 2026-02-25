package dto

import (
	"gorm.io/gorm"
)

const (
	DefaultPageSize    = 10
	DefaultPageIndex   = 1
	DefaultMaxPageSize = 500
)

type PaginationReq struct {
	PageIndex   int64 `form:"pageIndex" json:"pageIndex"`
	PageSize    int64 `form:"pageSize" json:"pageSize"`
	maxPageSize int64 // 单页最大数据
}

func (m *PaginationReq) SetMaxPageSize(maxPageSize int64) {
	m.maxPageSize = maxPageSize
}

func (m *PaginationReq) GetPageIndex() int64 {
	if m.PageIndex <= 0 {
		m.PageIndex = DefaultPageIndex
	}
	return m.PageIndex
}

func (m *PaginationReq) GetPageSize() int64 {

	// 设置单页最大数据
	var maxSize int64 = DefaultMaxPageSize
	if m.maxPageSize > 0 {
		maxSize = m.maxPageSize
	}

	if m.PageSize <= 0 {
		m.PageSize = DefaultPageSize
	}
	// 超过最大值时使用最大值
	if m.PageSize > maxSize {
		m.PageSize = maxSize
	}
	return m.PageSize
}

func (m *PaginationReq) MakePagination() func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Offset(int((m.GetPageIndex() - 1) * m.GetPageSize())).Limit(int(m.GetPageSize()))
	}
}
