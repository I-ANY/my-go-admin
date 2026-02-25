package service

import (
	"biz-auto-api/internal/apps/system/service/dto"
	"biz-auto-api/internal/models"
	"biz-auto-api/pkg/service"
	"biz-auto-api/pkg/tools"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

type SysOperaLog struct {
	service.Service
}

func (s SysOperaLog) GetOperaLogList(req *dto.GetOperaLogListReq) (*dto.GetOperaLogListRes, error) {
	var result = &dto.GetOperaLogListRes{}
	db := s.DB
	log := s.Log
	var total int64
	var operaLogs []*models.SysOperaLog

	err := db.Model(&models.SysOperaLog{}).Scopes(func(tx *gorm.DB) *gorm.DB {
		if len(req.Uri) > 0 {
			tx.Where("uri like ?", tools.FuzzyQuery(req.Uri))
		}
		if len(req.RequestMethod) > 0 {
			tx.Where("request_method = ?", req.RequestMethod)
		}
		if len(req.ClientIp) > 0 {
			tx.Where("client_ip like ?", tools.FuzzyQuery(req.ClientIp))
		}
		if len(req.Handler) > 0 {
			tx.Where("handler like ?", tools.FuzzyQuery(req.Handler))
		}
		if len(req.HttpCode) > 0 {
			tx.Where("http_code like ?", tools.FuzzyQuery(req.HttpCode))
		}
		if len(req.BizCode) > 0 {
			tx.Where("biz_code like ?", tools.FuzzyQuery(req.BizCode))
		}
		if len(req.RequestTimeRangeStart) > 0 && len(req.RequestTimeRangeEnd) > 0 {
			tx.Where("request_time >= ? and request_time <= ?", req.RequestTimeRangeStart, req.RequestTimeRangeEnd)
		}
		if req.UserId > 0 {
			tx.Where("create_by = ?", req.UserId)
		}
		if len(req.RequestId) > 0 {
			tx.Where("request_id = ?", req.RequestId)
		}
		if req.HandleSource > 0 {
			tx.Where("handle_source = ? ", req.HandleSource)
		}
		return tx
	}).Count(&total).Order("request_time desc").Scopes(req.MakePagination()).Find(&operaLogs).Error
	if err != nil {
		err = errors.Wrap(err, "query user failed")
		log.Errorf("%+v", err)
		return nil, err
	}
	var items = make([]*dto.SysOperaLog, 0, len(operaLogs))
	userIds := tools.GetSlice(operaLogs, func(e *models.SysOperaLog) int64 {
		return e.CreateBy
	})
	userIds = tools.RemoveDuplication(userIds)
	users := make([]*models.SysUser, 0)
	// 查询本次涉及到的用户信息
	err = db.Model(models.SysUser{}).Select("id", "username", "nick_name").Where("id in ? ", userIds).Find(&users).Error
	if err != nil {
		err = errors.Wrap(err, "query users info failed")
		log.Errorf("%+v", err)
		return nil, err
	}
	usersMap := tools.Slice2Map(users, func(user *models.SysUser) int64 {
		return user.Id
	})
	for _, operaLog := range operaLogs {
		var item = dto.SysOperaLog{}
		item.UserId = operaLog.CreateBy
		item.RequestTime = operaLog.RequestTime.Format(time.DateTime)
		if user, ok := usersMap[operaLog.CreateBy]; ok {
			item.UserInfo = dto.UserBaseInfo{
				Id:       user.Id,
				Username: user.Username,
				NickName: user.NickName,
			}
		}
		err = copier.Copy(&item, operaLog)
		if err != nil {
			err = errors.Wrap(err, "convert data failed")
			log.Errorf("%+v", err)
			return nil, err
		}
		items = append(items, &item)
	}
	result.Items = items
	result.Total = total
	return result, nil
}
