package sexyimage

import (
	"errors"

	"github.com/go-resty/resty/v2"
)

type APIResponse struct {
	ErrCode string `json:"errCode"`
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    []struct {
		Pid         int64     `json:"pid"`
		Uid         int64     `json:"uid"`
		Title       string    `json:"title"`
		Author      string    `json:"author"`
		Width       int       `json:"width"`
		Height      int       `json:"height"`
		Ext         string    `json:"ext"`
		AiType      int       `json:"aiType"`
		TagsList    []Tag     `json:"tagsList"`
		UrlsList    []URLItem `json:"urlsList"`
		PcreateDate int64     `json:"pcreateDate"`
		PuploadDate int64     `json:"puploadDate"`
	} `json:"data"`
}

type Data struct {
	Pid         int64     `json:"pid"`
	Uid         int64     `json:"uid"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	Width       int       `json:"width"`
	Height      int       `json:"height"`
	Ext         string    `json:"ext"`
	AiType      int       `json:"aiType"`
	TagsList    []Tag     `json:"tagsList"`
	UrlsList    []URLItem `json:"urlsList"`
	PcreateDate int64     `json:"pcreateDate"`
	PuploadDate int64     `json:"puploadDate"`
}

type Tag struct {
	TagName string `json:"tagName"`
	TagEn   string `json:"tagEn,omitempty"`
}

type URLItem struct {
	URLSize string `json:"urlSize"`
	URL     string `json:"url"`
}

type SearchParams struct {
	Num           int      `json:"num,omitempty"`           // 1~20
	Pid           []int64  `json:"pid,omitempty"`           // 指定 pid
	Uid           []int64  `json:"uid,omitempty"`           // 指定 uid
	Author        string   `json:"author,omitempty"`        // 作者名模糊搜索
	Proxy         string   `json:"proxy,omitempty"`         // 默认 i.pixiv.re
	AiType        int      `json:"aiType,omitempty"`        // 0-未知[旧画作] 1-否 2-是
	R18Type       int      `json:"r18Type,omitempty"`       // 0-不是 1-是
	DateAfter     int64    `json:"dateAfter,omitempty"`     // 毫秒时间戳
	DateBefore    int64    `json:"dateBefore,omitempty"`    // 毫秒时间戳
	SizeList      []string `json:"sizeList,omitempty"`      // 如 ["original"]
	ImageSizeType int      `json:"imageSizeType,omitempty"` // 1-横图 2-竖图 3-方图
}

func getOneSexyImage(params SearchParams) (url string, err error) {
	if params.Num <= 0 {
		params.Num = 1
	}
	c := resty.New()
	var result APIResponse
	resp, err := c.R().
		SetResult(APIResponse{}).
		SetResult(&result).
		SetBody(&params).
		Post("https://api.mossia.top/duckMo")
	if err != nil {
		return url, err
	}
	if resp.IsError() {
		return url, errors.New(resp.Status())
	}
	if len(result.Data) == 0 {
		return url, errors.New("图片获取失败")
	}
	if len(result.Data[0].UrlsList) == 0 {
		return url, errors.New("图片URL获取失败")
	}
	return result.Data[0].UrlsList[0].URL, nil
}
