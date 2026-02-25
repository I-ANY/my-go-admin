package api

import (
	"biz-auto-api/pkg/consts"
	"net/http"

	"github.com/gin-gonic/gin"
)

var Default = &Response{}

type Response struct {
	// 数据集
	RequestId string      `protobuf:"bytes,1,opt,name=requestId,proto3" json:"requestId,omitempty"`
	Code      int32       `protobuf:"varint,2,opt,name=code,proto3" json:"code,omitempty"`
	Msg       string      `protobuf:"bytes,3,opt,name=msg,proto3" json:"msg,omitempty"`
	Data      interface{} `json:"data"`
}

type PageDate struct {
	Total     int64       `json:"total"`
	PageIndex int64       `json:"pageIndex"`
	PageSize  int64       `json:"pageSize"`
	Items     interface{} `json:"items"`
}

func (e *Response) SetData(data interface{}) {
	e.Data = data
}

func (e Response) Clone() *Response {
	return &e
}

func (e *Response) SetRequestId(id string) {
	e.RequestId = id
}

func (e *Response) SetMsg(s string) {
	e.Msg = s
}

func (e *Response) SetCode(code int32) {
	e.Code = code
}

// Error 失败数据处理
func Error(c *gin.Context, code int, msg string) {
	res := Default.Clone()
	res.SetMsg(msg)
	res.SetRequestId(GenerateMsgIDFromContext(c))
	res.SetCode(int32(code))
	c.AbortWithStatusJSON(code, res)
}

func OKWithBizCode(c *gin.Context, code consts.BizCode, msg string) {
	res := Default.Clone()
	res.SetMsg(msg)
	res.SetRequestId(GenerateMsgIDFromContext(c))
	res.SetCode(int32(code))
	c.AbortWithStatusJSON(http.StatusOK, res)
}

// OK 通常成功数据处理
func OK(c *gin.Context, data interface{}, msg string) {
	res := Default.Clone()
	res.SetRequestId(GenerateMsgIDFromContext(c))
	res.SetData(data)
	res.SetMsg(msg)
	res.SetCode(int32(consts.BizCode200))
	c.AbortWithStatusJSON(http.StatusOK, res)
}
func OKWithCodeAndData(c *gin.Context, code consts.BizCode, data interface{}, msg string) {
	res := Default.Clone()
	res.SetRequestId(GenerateMsgIDFromContext(c))
	res.SetData(data)
	res.SetMsg(msg)
	res.SetCode(int32(code))
	c.AbortWithStatusJSON(http.StatusOK, res)
}

// PageOK 分页数据处理
func PageOK(c *gin.Context, items interface{}, count, pageIndex, pageSize int64, msg string) {
	var res PageDate
	res.Items = items
	res.Total = count
	res.PageIndex = pageIndex
	res.PageSize = pageSize
	OK(c, res, msg)
}
