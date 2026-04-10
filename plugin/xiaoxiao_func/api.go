package xiaoxiaofunc

import (
	"errors"

	"github.com/go-resty/resty/v2"
)

type XiaoxiaoApi struct {
	client  *resty.Client
	baseUrl string
}

var XiaoxiaoGetter = NewXiaoxiaoApi("https://v2.xxapi.cn/api/")

func NewXiaoxiaoApi(baseUrl string) *XiaoxiaoApi {
	return &XiaoxiaoApi{
		baseUrl: baseUrl,
		client:  resty.New(),
	}
}

func (c *XiaoxiaoApi) GetHistoryToday() ([]byte, error) {
	resp, err := c.client.R().
		SetQueryParams(map[string]string{
			"return": "302",
		}).
		Get(c.baseUrl + "historypic")
	if err != nil {
		return nil, err
	}
	return resp.Body(), nil
}

type LotteryResponse struct {
	Code      int         `json:"code"`
	Msg       string      `json:"msg"`
	Data      LotteryData `json:"data"`
	RequestID string      `json:"request_id"`
}

// 签文数据结构体：对应 data 字段
type LotteryData struct {
	Content string `json:"content"` // 解签内容
	ID      int    `json:"id"`      // 签号
	Pic     string `json:"pic"`     // 图片地址
	Poem    string `json:"poem"`    // 签诗
	Title   string `json:"title"`   // 签文标题
}

func (c *XiaoxiaoApi) GetLottery() (LotteryResponse, error) {
	var response LotteryResponse
	resp, err := c.client.R().
		SetResult(&response).
		Get(c.baseUrl + "wenchangdijunrandom")
	if err != nil {
		return LotteryResponse{}, errors.New("获取签文失败: " + err.Error())
	}
	if resp.IsError() {
		return LotteryResponse{}, errors.New("获取签文失败: " + resp.Status())
	}
	if response.Code != 200 {
		return LotteryResponse{}, errors.New("获取签文失败: " + response.Msg)
	}
	return response, nil
}
