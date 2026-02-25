package common

import (
	"biz-auto-api/internal/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type DictTypeCode = string

type DictDataLabel = string
type DictDataValue = string
type DictData = map[DictDataValue]DictDataLabel

// K业务
const (
	KDeviceCategory DictTypeCode = "K_DEVICE_CATEGORY" // 设备大类
	KDeviceType     DictTypeCode = "K_DEVICE_TYPE"     // 设备小类
	KIsp            DictTypeCode = "K_ISP"             // 运营商
	KScheduleType   DictTypeCode = "K_SCHEDULE_TYPE"
	KIsCoverDiffIsp DictTypeCode = "K_IS_COVER_DIFF_ISP"
	KProviderType   DictTypeCode = "K_PROVIDER_TYPE"
	KHddLossReason  DictTypeCode = "K_HDD_LOSE_REASON"
)

// 巡检
const (
	InspectExecStatus DictTypeCode = "INSPECT_EXEC_STATUS"
	InspectResult     DictTypeCode = "INSPECT_INSPECT_RESULT"
)

// LA 数据字典
const (
	LAIsDiffIsp      DictTypeCode = "LA_IS_DIFF_ISP"
	LABusinessStatus DictTypeCode = "LA_BUSINESS_STATUS"
)

// 公共
const (
	CommonYesNo = "COMMON_YES_NO"
	//CommonOnlineStatus 公共_在线状态
	CommonOnlineStatus = "COMMON_ONLINE_STATUS"
)

func QueryDictByKeys(db *gorm.DB, keys []string) (map[DictTypeCode]DictData, error) {
	var (
		res   = make(map[DictTypeCode]DictData)
		dicts = make([]*models.SysDictType, 0, len(keys))
	)
	if len(keys) == 0 {
		return res, nil
	}
	err := db.Model(models.SysDictType{}).Preload("DictData").Scopes(func(tx *gorm.DB) *gorm.DB {
		tx.Where("type_code in ?", keys)
		return tx
	}).Find(&dicts).Error
	if err != nil {
		return res, errors.WithStack(err)
	}
	for _, dict := range dicts {
		values := make(DictData)
		for _, dictData := range dict.DictData {
			values[dictData.DictValue] = dictData.DictLabel
		}
		res[dict.TypeCode] = values
	}
	return res, nil
}

func QueryDictLabel(db *gorm.DB, key DictTypeCode, value DictDataValue) (DictDataLabel, error) {
	// TODO
	return "", nil
}

func QueryDictValue(db *gorm.DB, key DictTypeCode, label DictDataLabel) (DictDataValue, error) {
	// TODO
	return "", nil
}
