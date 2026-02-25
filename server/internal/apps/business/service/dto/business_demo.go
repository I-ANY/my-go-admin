package dto

type DemoItem struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type DemoReq struct {
	Id   int64  `json:"id" form:"id"`
	Name string `json:"name" form:"name" binding:"required"`
}

type DemoRes struct {
	Id    int64      `json:"id"`
	Items []DemoItem `json:"items"`
}
