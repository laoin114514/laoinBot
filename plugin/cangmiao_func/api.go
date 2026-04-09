package cangmiao_func

import (
	"errors"
	"laoinBot/config"
	"net/http"

	"github.com/go-resty/resty/v2"
)

var ImageRepos = map[string]string{
	"cos图":  "cos",
	"瑟瑟":    "sese",
	"美女":    "meinv",
	"黑丝":    "heisitu",
	"白丝":    "baisitu",
	"美腿":    "meitui",
	"jk":    "jktu",
	"小姐姐":   "pcxjj",
	"原神cos": "yscos",
	"买家秀":   "maijiaxiu",
	"壁纸":    "bizhi",
	"4K壁纸":  "4kacg",
}

type ImageResponse struct {
	Img string `json:"img"`
}

type CangMiaoApi struct {
	baseURL string
	apikey  string
	client  *resty.Client
}

var CangMiaoGetter *CangMiaoApi

func init() {
	config.LoadConfig()
	CangMiaoGetter = NewCangMiaoApi("http://api.tinise.cn/api/", config.BotConfig.CangMiaoKey)
}
func NewCangMiaoApi(baseURL string, apikey string) *CangMiaoApi {
	return &CangMiaoApi{
		baseURL: baseURL,
		apikey:  apikey,
		client:  resty.New(),
	}
}
func (c *CangMiaoApi) GetImageURL(repo string) (*ImageResponse, error) {
	var result ImageResponse
	resp, err := c.client.R().
		SetResult(&result).
		SetQueryParams(map[string]string{
			"apikey": c.apikey,
		}).
		Get(c.baseURL + repo)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, errors.New("获取图片失败: " + resp.Status())
	}
	return &result, nil
}

func (i ImageResponse) ToByte() []byte {
	resp, err := CangMiaoGetter.client.R().
		Get(i.Img)
	if err != nil {
		return nil
	}
	if resp.IsError() {
		return nil
	}
	return resp.Body()
}

func (c *CangMiaoApi) GetPhoneAdress(phoneNum string) (string, error) {
	if len(phoneNum) < 11 {
		return "", errors.New("手机号长度不正确")
	}
	resp, err := c.client.R().
		SetQueryParams(map[string]string{
			"id": phoneNum,
		}).
		Get(c.baseURL + "phone")
	if err != nil {
		return "", errors.New("获取号码归属地失败: " + err.Error())
	}
	if resp.IsError() {
		return "", errors.New("获取号码归属地失败: " + resp.Status())
	}
	if resp.StatusCode() != http.StatusOK {
		return "", errors.New("获取号码归属地失败: " + resp.Status())
	}
	return string(resp.Body()), nil
}

type IPResponse struct {
	IP        string  `json:"ip"`
	Latitude  float64 `json:"latitude"`
	Rectangle string  `json:"rectangle"`
	Location  string  `json:"location"`
	Timestamp int64   `json:"timestamp"`
}

func (c *CangMiaoApi) GetIPAdress(ip string) (string, error) {
	var result IPResponse
	resp, err := c.client.R().
		SetResult(&result).
		SetQueryParams(map[string]string{
			"ip": ip,
		}).
		Get(c.baseURL + "ip")
	if err != nil {
		return "", errors.New("获取IP归属地失败: " + err.Error())
	}
	if resp.IsError() {
		return "", errors.New("获取IP归属地失败: " + resp.Status())
	}
	if resp.StatusCode() != http.StatusOK {
		return "", errors.New("获取IP归属地失败: " + resp.Status())
	}
	return result.Location, nil
}
