package service

import (
	"biz-auto-api/internal/apps/business/service/dto"
	"biz-auto-api/pkg/service"
)

type BusinessDemoService struct {
	service.Service
}

func (s *BusinessDemoService) Demo(req *dto.DemoReq) (res dto.DemoRes, error error) {
	var items []dto.DemoItem
	return dto.DemoRes{
		Id:    req.Id,
		Items: items,
	}, nil
}
