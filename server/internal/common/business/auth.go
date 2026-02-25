package business

import (
	"biz-auto-api/pkg/consts"
	"biz-auto-api/pkg/grpc/clients"
	pbauth "biz-auto-api/pkg/grpc/pb/auth"
	"context"
	"github.com/pkg/errors"
	"strconv"
)

func HasPermission(ctx context.Context, userId int64, categoryId int64) (bool, error) {
	r := &pbauth.CheckPermissionReq{
		UserID:  userId,
		ResType: string(consts.Category),
		ResID:   strconv.FormatInt(categoryId, 10),
		Action:  string(consts.CategoryGeneralPermission),
	}
	res, err := clients.GetAuthGrpcClient().CheckPermission(ctx, r)
	if err != nil {
		return false, errors.Wrap(err, "check permission failed")
	}
	if res != nil {
		return res.Allowed, nil
	}
	return false, nil
}
func AuthedCategoryIds(ctx context.Context, userId int64) ([]int64, error) {
	r := &pbauth.GetUserAuthorizedResIdsReq{
		UserID:  userId,
		ResType: string(consts.Category),
		Action:  string(consts.CategoryGeneralPermission),
	}
	res, err := clients.GetAuthGrpcClient().GetUserAuthorizedResIds(ctx, r)
	if err != nil {
		return nil, errors.Wrap(err, "get authed category failed")
	}
	return res.ResIds, nil
}
