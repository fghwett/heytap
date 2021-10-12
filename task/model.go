package task

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

type UserInfoResponse struct {
	Id            string `json:"id"`
	AccountName   string `json:"accountName"`
	UserName      string `json:"userName"`
	AvatarDefault string `json:"avatarDefault"`
	Country       string `json:"country"`
	Birthday      string `json:"birthday"`
	Sex           string `json:"sex"`
	RealName      string `json:"realName"`
	Oid           string `json:"oid"`
	UserIdStr     string `json:"userIdStr"`
}
