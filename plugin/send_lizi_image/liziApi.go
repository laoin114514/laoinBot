package send_lizi_image

import (
	"github.com/go-resty/resty/v2"
)

var HasRepo = map[string]string{
	"二次元":  "ycy",
	"萌版":   "moez",
	"AI":   "ai",
	"原神":   "ysz",
	"随机横图": "pc",
	"风景":   "fj",
	"随机竖图": "mp",
	"萌版竖图": "moemp",
	"原神竖图": "ysmp",
	"AI竖图": "aimp",
	"头像":   "tx",
	"白底":   "bd",
}
var LiziGetter = NewLiziApi("https://t.alcy.cc/")

type LiziApi struct {
	baseURL string
	client  *resty.Client
}

func NewLiziApi(baseURL string) *LiziApi {
	return &LiziApi{
		baseURL: baseURL,
		client:  resty.New(),
	}
}

func (l *LiziApi) GetOneImage(suffix string) ([]byte, error) {
	resp, err := l.client.R().
		Get(l.baseURL + suffix)
	if err != nil {
		return nil, err
	}
	return resp.Body(), nil
}
