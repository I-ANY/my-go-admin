package dto

type AlarmData struct {
	Level  string `json:"level" form:"level"`
	Topic  string `json:"topic" form:"topic"`
	Detail string `json:"detail" form:"detail"`
	Status int    `json:"status" form:"status"`
}
