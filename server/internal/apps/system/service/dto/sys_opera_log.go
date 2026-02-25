package dto

import "biz-auto-api/pkg/dto"

type GetOperaLogListReq struct {
	dto.PaginationReq
	Uri                   string `json:"uri" form:"uri"`
	RequestMethod         string `json:"requestMethod" form:"requestMethod"`
	ClientIp              string `json:"clientIp" form:"clientIp"`
	Handler               string `json:"handler" form:"handler"`
	HttpCode              string `json:"httpCode" form:"httpCode"`
	BizCode               string `json:"bizCode" form:"bizCode"`
	RequestTimeRangeStart string `json:"requestTimeRangeStart" form:"requestTimeRangeStart" validate:"required,datetime=2006-01-02 15:04:05"`
	RequestTimeRangeEnd   string `json:"requestTimeRangeEnd" form:"requestTimeRangeEnd" validate:"required,datetime=2006-01-02 15:04:05"`
	UserId                int64  `json:"userId" form:"userId"`
	RequestId             string `json:"requestId" form:"requestId"`
	HandleSource          int64  `json:"handleSource" form:"handleSource" validate:"required"`
}

type SysOperaLog struct {
	Id            int64        `json:"id"`
	Uri           string       `json:"uri"`
	Api           string       `json:"api"`
	RequestId     string       `json:"requestId"`
	RequestMethod string       `json:"requestMethod"`
	ClientIp      string       `json:"clientIp"`
	LatencyTime   string       `json:"latencyTime"`
	UserAgent     string       `json:"userAgent"`
	ReqBody       string       `json:"reqBody"`
	JsonRes       string       `json:"jsonRes"`
	HttpCode      string       `json:"httpCode"`
	BizCode       string       `json:"bizCode"`
	RequestTime   string       `json:"requestTime"`
	Handler       string       `json:"handler"`
	UserId        int64        `json:"userId"`
	UserInfo      UserBaseInfo `json:"userInfo"`
	HandleSource  int64        `json:"handleSource"`
}

type GetOperaLogListRes struct {
	Items []*SysOperaLog `json:"items"`
	Total int64          `json:"total"`
}
