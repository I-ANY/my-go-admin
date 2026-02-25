package dto

type DictType struct {
	Id       int64       `json:"id"`
	TypeName string      `json:"typeName"`
	TypeCode string      `json:"typeCode"`
	Sort     int64       `json:"sort"`
	Status   int64       `json:"status"`
	Remark   string      `json:"remark"`
	DictData []*DictData `json:"dictData,omitempty"`
}
type DictData struct {
	Id         int64  `json:"id"`
	DictTypeId int64  `json:"dictTypeId"`
	DictLabel  string `json:"dictLabel"`
	DictValue  string `json:"dictValue"`
	Sort       int64  `json:"sort"`
	Status     int64  `json:"status"`
	Remark     string `json:"remark"`
	Color      string `json:"color"`
}

type GetAllDictDataReq struct {
}
