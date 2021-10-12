package task

import (
	"net/http"

	"github.com/fghwett/heytap/config"
	"github.com/fghwett/heytap/util"
)

type Task struct {
	cookie    string
	userAgent string
	client    *http.Client

	username string
}

func New(config *config.Config) *Task {
	return &Task{
		cookie:    config.Cookie,
		userAgent: config.UserAgent,
		client:    &http.Client{},
	}
}

func (t *Task) Do() error {

	return nil
}

func (t *Task) getUserInfo() error {
	reqUrl := "https://www.heytap.com/cn/oapi/users/web/member/info"
	req, err := http.NewRequest(http.MethodGet, reqUrl, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Host", "www.heytap.com")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("User-Agent", t.userAgent)
	req.Header.Set("Accept-Language", "zh-cn")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("cookie", t.cookie)

	resp, err := t.client.Do(req)

	response := &Response{Data: &UserInfoResponse{}}
	err = util.GetHTTPResponse(resp, reqUrl, err, response)
	if err != nil {
		return err
	}
	t.username = response.Data.(*UserInfoResponse).RealName

	return nil
}
