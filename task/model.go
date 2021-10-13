package task

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

// 获取用户信息返回值
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

// 任务中心返回值
type CenterResponse struct {
	UserCredits        int                 `json:"userCredits"`
	NoviceList         []interface{}       `json:"noviceList"`
	EverydayList       []*List             `json:"everydayList"`
	BuyList            []*List             `json:"buyList"`
	StartTime          interface{}         `json:"startTime"`
	EndTime            interface{}         `json:"endTime"`
	CurrentTime        int64               `json:"currentTime"`
	WaterStartTime     int64               `json:"waterStartTime"`
	WaterEndTime       int64               `json:"waterEndTime"`
	IsStart            int                 `json:"isStart"`
	UserReportInfoForm *UserReportInfoForm `json:"userReportInfoForm"`
}

type List struct {
	Id                  string            `json:"id"`
	CreatedAt           int64             `json:"createdAt"`
	UpdatedAt           int64             `json:"updatedAt"`
	Marking             string            `json:"marking"`
	Sort                int               `json:"sort"`
	Name                string            `json:"name"`
	Title               string            `json:"title"`
	ButtonName          string            `json:"buttonName"`
	IconUrl             string            `json:"iconUrl"`
	Type                int               `json:"type"`
	Times               int               `json:"times"`
	Credits             int               `json:"credits"`
	WaterAmount         int               `json:"waterAmount"`
	Link                string            `json:"link"`
	Status              int               `json:"status"`
	VersionDetailss     interface{}       `json:"versionDetailss"`
	ChannelDetailss     []*ChannelDetails `json:"channelDetailss"`
	MobileBrandDetailss interface{}       `json:"mobileBrandDetailss"`
	ReadCount           int               `json:"readCount"`
	CompleteStatus      int               `json:"completeStatus"`
}

type ChannelDetails struct {
	Id               string `json:"id"`
	CreatedAt        int64  `json:"createdAt"`
	UpdatedAt        int64  `json:"updatedAt"`
	ConfigSourceType string `json:"configSourceType"`
	ChannelCode      string `json:"channelCode"`
	DetailId         int    `json:"detailId"`
	ConfigKey        string `json:"configKey"`
	ShieldProducts   string `json:"shieldProducts"`
}

type UserReportInfoForm struct {
	ReportDays int     `json:"reportDays"`
	Status     int     `json:"status"`
	Gifts      []*Gift `json:"gifts"`
}

type Gift struct {
	Date    string      `json:"date"`
	Gift    interface{} `json:"gift"` // 空值string 有值 float64
	Credits int         `json:"credits"`
	Today   bool        `json:"today"`
	Id      int         `json:"id"`
	Type    interface{} `json:"type"` // 空值时是string // 有值时是float64
}

// 签到返回值
type SignResponse struct {
	Result  string `json:"result"`
	Message string `json:"message"`
}
