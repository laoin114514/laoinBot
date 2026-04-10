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

// GoldPriceResponse 整体返回结构体
type GoldPriceResponse struct {
	Code      int           `json:"code"`
	Msg       string        `json:"msg"`
	Data      GoldPriceData `json:"data"`
	RequestId string        `json:"request_id"`
}

// GoldPriceData data 部分
type GoldPriceData struct {
	BankGoldBarPrice   []BankGoldBar   `json:"bank_gold_bar_price"`
	GoldRecyclePrice   []GoldRecycle   `json:"gold_recycle_price"`
	PreciousMetalPrice []PreciousMetal `json:"precious_metal_price"`
}

// BankGoldBar 银行金条价格
type BankGoldBar struct {
	Bank  string `json:"bank"`
	Price string `json:"price"`
}

// GoldRecycle 黄金回收价格
type GoldRecycle struct {
	GoldType     string `json:"gold_type"`
	RecyclePrice string `json:"recycle_price"`
	UpdatedDate  string `json:"updated_date"`
}

// PreciousMetal 品牌贵金属价格
type PreciousMetal struct {
	Brand         string `json:"brand"`
	BullionPrice  string `json:"bullion_price"`
	GoldPrice     string `json:"gold_price"`
	PlatinumPrice string `json:"platinum_price"`
	UpdatedDate   string `json:"updated_date"`
}

func (c *XiaoxiaoApi) GetGoldPrice() (GoldPriceResponse, error) {
	var response GoldPriceResponse
	resp, err := c.client.R().
		SetResult(&response).
		Get(c.baseUrl + "goldprice")
	if err != nil {
		return GoldPriceResponse{}, errors.New("获取黄金价格失败: " + err.Error())
	}
	if resp.IsError() {
		return GoldPriceResponse{}, errors.New("获取黄金价格失败: " + resp.Status())
	}
	if response.Code != 200 {
		return GoldPriceResponse{}, errors.New("获取黄金价格失败: " + response.Msg)
	}
	return response, nil
}
