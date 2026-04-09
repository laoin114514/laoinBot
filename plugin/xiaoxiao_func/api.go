package xiaoxiaofunc

import "github.com/go-resty/resty/v2"

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
