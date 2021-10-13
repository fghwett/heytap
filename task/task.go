package task

import (
	"bytes"
	"fmt"
	"github.com/fghwett/heytap/config"
	"github.com/fghwett/heytap/util"
	"log"
	"net/http"
	"time"
)

type Task struct {
	cookie    string
	userAgent string
	client    *http.Client
	result    []string
}

func New(config *config.Config) *Task {
	return &Task{
		cookie:    config.Cookie,
		userAgent: config.UserAgent,
		client:    &http.Client{},
		result:    []string{"==== 欢太商城任务 ===="},
	}
}

func (t *Task) Do() error {
	if err := t.getUserInfo(); err != nil {
		return err
	}

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

	t.result = append(t.result, fmt.Sprintf("【登陆成功】：%s", response.Data.(*UserInfoResponse).RealName))

	return nil
}

func (t *Task) taskCenter() (*CenterResponse, error) {
	reqUrl := "https://store.oppo.com/cn/oapi/credits/web/credits/show"
	req, err := http.NewRequest(http.MethodGet, reqUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Host", "store.oppo.com")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("User-Agent", t.userAgent)
	req.Header.Set("Accept-Language", "zh-cn")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("cookie", t.cookie)
	req.Header.Set("referer", "https://store.oppo.com/cn/app/taskCenter/index")

	resp, err := t.client.Do(req)
	response := &Response{Data: &CenterResponse{}}
	err = util.GetHTTPResponse(resp, reqUrl, err, response)
	if err != nil {
		return nil, err
	}

	return response.Data.(*CenterResponse), nil
}

func (t *Task) signTask() error {
	center, err := t.taskCenter()
	if err != nil {
		return err
	}

	util.SmallSleep(500, 1500)

	if center.UserReportInfoForm.Status != 0 {
		t.result = append(t.result, "【每日签到】: 已经签到过了！")
		return nil
	}

	dated := time.Now().Format("2006-01-02")
	todayGist := &Gift{}
	for _, gift := range center.UserReportInfoForm.Gifts {
		if gift.Date == dated {
			todayGist = gift
		}
	}

	var reqData string
	if todayGist.Today {
		log.Printf("credits: %v, type: %v, gift: %v", todayGist.Credits, todayGist.Type, todayGist.Gift)

		if _, ok := todayGist.Type.(string); ok {
			reqData = fmt.Sprintf("amount=%d", todayGist.Credits)
		} else {
			reqData = fmt.Sprintf("amount=%d&type=%d&gift=%d", todayGist.Credits, int64(todayGist.Type.(float64)), int64(todayGist.Gift.(float64)))
		}
	} else {
		reqData = fmt.Sprintf("amount=%d", todayGist.Credits)
	}

	reqUrl := "https://store.oppo.com/cn/oapi/credits/web/report/immediately"
	req, err := http.NewRequest(http.MethodPost, reqUrl, bytes.NewReader([]byte(reqData)))
	if err != nil {
		return err
	}

	req.Header.Set("Host", "store.oppo.com")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("User-Agent", t.userAgent)
	req.Header.Set("Accept-Language", "zh-cn")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("cookie", t.cookie)
	req.Header.Set("referer", "https://store.oppo.com/cn/app/taskCenter/index")

	resp, err := t.client.Do(req)
	response := &Response{Data: &SignResponse{}}
	err = util.GetHTTPResponse(resp, reqUrl, err, response)
	if err != nil {
		return err
	}

	t.result = append(t.result, fmt.Sprintf("【每日签到】：成功 %s", response.Data.(*SignResponse).Message))

	return nil
}
