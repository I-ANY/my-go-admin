package dto

type GetDeptTreeReq struct {
}
type Dept struct {
	Id         int64   `json:"id"`
	Name       string  `json:"name"`
	Status     int64   `json:"status"`
	Remark     string  `json:"remark"`
	Children   []*Dept `json:"children"`
	ParentDept int64   `json:"parentDept"`
	OrderNo    int64   `json:"orderNo"`
}
type GetDeptTreeRes struct {
	Items []*Dept `json:"items"`
	Total int64   `json:"total"`
}
