package task

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/fghwett/heytap/config"
	"github.com/fghwett/heytap/util"
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

func (t *Task) Do() {
	if err := t.getUserInfo(); err != nil {
		t.result = append(t.result, fmt.Sprintf("【登陆失败】：%s", err))
		return
	}

	if err := t.signTask(); err != nil {
		t.result = append(t.result, fmt.Sprintf("【每日签到】：失败 %s", err))
		return
	}

	if err := t.viewGoodsTask(); err != nil {
		t.result = append(t.result, fmt.Sprintf("【每日浏览商品】：失败 %s", err))
		return
	}

	if err := t.shareGoodsTask(); err != nil {
		t.result = append(t.result, fmt.Sprintf("【每日分享商品】：失败 %s", err))
		return
	}
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
			reqData = fmt.Sprintf("amount=%d&type=%d&gift=", todayGist.Credits, int64(todayGist.Type.(float64)))
			if _, ok := todayGist.Gift.(float64); ok {
				reqData += fmt.Sprintf("%d", int64(todayGist.Gift.(float64)))
			} else {
				reqData += todayGist.Gift.(string)
			}
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

func (t *Task) viewGoodsTask() error {
	center, err := t.taskCenter()
	if err != nil {
		return err
	}

	util.SmallSleep(500, 1500)

	todayList := &List{}
	for _, list := range center.EverydayList {
		if list.Name == "浏览商品" {
			todayList = list
		}
	}
	if todayList.CompleteStatus == 0 {
		goodsDetails, err := t.getGoods()
		if err != nil {
			return err
		}

		for _, goodsDetail := range goodsDetails {
			log.Printf("正在浏览商品ID：%d", goodsDetail.Skuid)

			if err := t.viewGoods(goodsDetail.Skuid); err != nil {
				return err
			}

			util.SmallSleep(5000, 8000)
		}

		if err := t.cashingCredits(todayList.Marking, todayList.Type, todayList.Credits); err != nil {
			return err
		}

		t.result = append(t.result, fmt.Sprintf("【每日浏览商品】：任务完成！积分领取+%d", todayList.Credits))
	} else if todayList.CompleteStatus == 1 {
		if err := t.cashingCredits(todayList.Marking, todayList.Type, todayList.Credits); err != nil {
			return err
		}

		t.result = append(t.result, fmt.Sprintf("【每日浏览商品】：任务完成！积分领取+%d", todayList.Credits))
	} else {
		t.result = append(t.result, "【每日浏览商品】：任务已完成！")
	}

	return nil
}

func (t *Task) getGoods() ([]*GoodsDetail, error) {
	reqUrl := "https://msec.opposhop.cn/goods/v1/SeckillRound/goods/115?pageSize=12&currentPage=1"
	req, err := http.NewRequest(http.MethodGet, reqUrl, nil)
	if err != nil {
		return nil, err
	}

	resp, err := t.client.Do(req)

	response := &GoodsResponse{}
	err = util.GetHTTPResponse(resp, reqUrl, err, response)
	if err != nil {
		return nil, err
	}

	if response.Meta.Code != 200 {
		return nil, fmt.Errorf("获取商品列表失败")
	}

	return response.Detail, nil
}

func (t *Task) viewGoods(skuid int) error {
	reqUrl := fmt.Sprintf("https://msec.opposhop.cn/goods/v1/info/sku?skuId=%d", skuid)
	req, err := http.NewRequest(http.MethodGet, reqUrl, nil)
	if err != nil {
		return err
	}

	req.Header.Set("clientPackage", "com.oppo.store")
	req.Header.Set("Host", "store.oppo.com")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("User-Agent", "okhttp/3.12.12.200sp1")
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("cookie", t.cookie)

	resp, err := t.client.Do(req)
	if err != nil {
		return err
	}

	response := &ViewGoodsResponse{}
	err = util.GetHTTPResponse(resp, reqUrl, err, response)
	if err != nil {
		return err
	}

	if response.Meta.Code != 200 {
		return fmt.Errorf("浏览商品失败：%d", skuid)
	}

	return nil
}

func (t *Task) cashingCredits(infoMarking string, infoType int, infoCredits int) error {
	reqData := fmt.Sprintf("marking=%s&type=%d&amount=%d", infoMarking, infoType, infoCredits)

	reqUrl := "https://store.oppo.com/cn/oapi/credits/web/credits/cashingCredits"
	req, err := http.NewRequest(http.MethodPost, reqUrl, bytes.NewReader([]byte(reqData)))
	if err != nil {
		return err
	}

	req.Header.Set("Host", "store.oppo.com")
	req.Header.Set("clientPackage", "com.oppo.store")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("User-Agent", t.userAgent)
	req.Header.Set("Accept-Language", "zh-cn")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("cookie", t.cookie)
	req.Header.Set("Origin", "https://store.oppo.com")
	req.Header.Set("X-Requested-With", "com.oppo.store")
	req.Header.Set("referer", "https://store.oppo.com/cn/app/taskCenter/index?us=gerenzhongxin&um=hudongleyuan&uc=renwuzhongxin")

	resp, err := t.client.Do(req)
	if err != nil {
		return err
	}

	response := &Response{}
	err = util.GetHTTPResponse(resp, reqUrl, err, response)
	if err != nil {
		return err
	}

	if response.Code != 200 {
		return fmt.Errorf("领取积分奖励出错！%s", response.ErrorMessage)
	}

	return nil
}

func (t *Task) shareGoodsTask() error {
	center, err := t.taskCenter()
	if err != nil {
		return err
	}

	util.SmallSleep(500, 1500)

	todayList := &List{}
	for _, list := range center.EverydayList {
		if list.Name == "分享商品到微信" {
			todayList = list
		}
	}

	if todayList.CompleteStatus == 0 {
		count := todayList.ReadCount
		endCount := todayList.Times
		for {
			if count >= endCount {
				break
			}

			if shareErr := t.shareGoods(); shareErr != nil {
				return shareErr
			}

			util.SmallSleep(2000, 5000)
			count++
		}

		if err := t.cashingCredits(todayList.Marking, todayList.Type, todayList.Credits); err != nil {
			return err
		}

		t.result = append(t.result, fmt.Sprintf("【每日分享商品】：任务完成！积分领取+%d", todayList.Credits))
	} else if todayList.CompleteStatus == 1 {
		if err := t.cashingCredits(todayList.Marking, todayList.Type, todayList.Credits); err != nil {
			return err
		}

		t.result = append(t.result, fmt.Sprintf("【每日分享商品】：任务完成！积分领取+%d", todayList.Credits))
	} else {
		t.result = append(t.result, "【每日分享商品】：任务已完成！")
	}

	return nil
}

func (t *Task) shareGoods() error {
	reqUrl := fmt.Sprintf("https://msec.opposhop.cn/users/vi/creditsTask/pushTask?marking=daily_sharegoods")
	req, err := http.NewRequest(http.MethodGet, reqUrl, nil)
	if err != nil {
		return err
	}

	req.Header.Set("clientPackage", "com.oppo.store")
	req.Header.Set("Host", "store.oppo.com")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("User-Agent", "okhttp/3.12.12.200sp1")
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("cookie", t.cookie)

	resp, err := t.client.Do(req)

	response := &ShareGoodsResponse{}
	err = util.GetHTTPResponse(resp, reqUrl, err, response)
	if err != nil {
		return err
	}

	if response.Meta.Code != 200 {
		return fmt.Errorf("分享商品失败 %v", err)
	}

	return nil
}

func (t *Task) GetResult() string {
	return strings.Join(t.result, " \n\n ")
}
