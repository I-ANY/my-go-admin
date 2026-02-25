package dto

type User struct {
	Id       *int64  `json:"id"`
	Username *string `json:"username"`
	NickName *string `json:"nickName"`
}
