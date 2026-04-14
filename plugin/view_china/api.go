package viewchina

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"
	transformer "github.com/mozillazg/go-pinyin"
	"golang.org/x/crypto/bcrypt"
)

type ViewChinaApi struct {
	baseURL string
	client  *resty.Client
}

func NewViewChinaApi() *ViewChinaApi {
	return &ViewChinaApi{
		baseURL: "https://www.vcg.com",
		client:  resty.New(),
	}
}

// GenerateSec 生成 Sec 请求头（已修复72字节截断）
func (v *ViewChinaApi) GenerateSec(params any) (string, error) {
	// 1. JSON 序列化
	jsonBytes, _ := json.Marshal(params)
	jsonStr := string(jsonBytes)

	// 2. URL 编码
	urlEncoded := url.QueryEscape(jsonStr)

	// 3. Base64
	base64Str := base64.StdEncoding.EncodeToString([]byte(urlEncoded))

	// 4. 拼接密钥
	fullStr := "SECRET_VCG_" + base64Str

	// ---------------- FIX HERE ----------------
	// 关键：截断到 72 字节（和网站JS逻辑一致）
	if len(fullStr) > 72 {
		fullStr = fullStr[:72]
	}
	// -----------------------------------------

	// 5. bcrypt hash
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(fullStr), 8)
	if err != nil {
		return "", err
	}

	return string(hashBytes), nil
}

type SearchImagesParams struct {
	Phrase    string `json:"phrase"`
	Transform string `json:"transform"`
	Page      int    `json:"page"`
	Sort      string `json:"sort"` //fresh or hot
}

func (v *ViewChinaApi) SearchImages(phrase string, page int, sort string) (Response, error) {
	pinyin := transformer.Slug(phrase, transformer.Args{})
	if pinyin == "" {
		pinyin = phrase
	}
	params := map[string]string{
		"phrase":    phrase,
		"transform": pinyin,
		"page":      strconv.Itoa(page),
		"sort":      sort,
	}
	sec, err := v.GenerateSec(params)
	if err != nil {
		return Response{}, err
	}
	fmt.Println(sec)
	var result Response
	resp, err := v.client.R().
		SetResult(&result).
		SetHeader("sec", sec).
		SetQueryParams(params).
		Get(v.baseURL + "/api/common/searchAllImage")
	if err != nil {
		return Response{}, err
	}
	if resp.IsError() {
		return Response{}, fmt.Errorf("search images failed: http %d", resp.StatusCode())
	}
	return result, nil
}

func (v *ViewChinaApi) DownloadImage(rawURL string) ([]byte, error) {
	imageURL := strings.TrimSpace(rawURL)
	if imageURL == "" {
		return nil, fmt.Errorf("empty image url")
	}
	if strings.HasPrefix(imageURL, "//") {
		imageURL = "https:" + imageURL
	}

	resp, err := v.client.R().
		SetHeader("referer", v.baseURL).
		SetHeader("user-agent", "Mozilla/5.0").
		Get(imageURL)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("download image failed: http %d", resp.StatusCode())
	}
	return resp.Body(), nil
}

// 根结构体
type Response struct {
	List             []Item `json:"list"`
	ResUrl           ResUrl `json:"resUrl"`
	RecommendCount   int    `json:"recommend_count"`
	TotalCount       int    `json:"total_count"`
	ErrorCode        any    `json:"error_code"`
	AmbigInfo        any    `json:"ambigInfo"`
	ExtendKeywords   any    `json:"extendKeywords"`
	Sign             any    `json:"sign"`
	ActivateAISearch int    `json:"activateAISearch"`
	Message          any    `json:"message"`
	IP               string `json:"ip"`
	UserID           string `json:"user_id"`
	RealTotalCount   int    `json:"realTotalCount"`
}

// 列表里的每一项
type Item struct {
	ID                      int64      `json:"id"`
	ResID                   any        `json:"resId"`
	ProviderResID           any        `json:"providerResId"`
	BrandID                 int        `json:"brandId"`
	ProviderID              any        `json:"providerId"`
	ProviderAgentType       any        `json:"providerAgentType"`
	CollectionID            any        `json:"collectionId"`
	CfpGicID                any        `json:"cfpGicId"`
	AssetFamily             int        `json:"assetFamily"`
	AssetType               any        `json:"assetType"`
	LicenseType             int        `json:"licenseType"`
	GraphicalStyle          int        `json:"graphicalStyle"`
	AssetFormat             string     `json:"assetFormat"`
	Title                   string     `json:"title"`
	OneCategory             any        `json:"oneCategory"`
	Category                string     `json:"category"`
	QualityRank             int        `json:"qualityRank"`
	Keywords                string     `json:"keywords"`
	ImageState              any        `json:"imageState"`
	OnlineState             int        `json:"onlineState"`
	OnlineTime              string     `json:"onlineTime"`
	OfflineTime             any        `json:"offlineTime"`
	OfflineReason           any        `json:"offlineReason"`
	OfflineMark             any        `json:"offlineMark"`
	KeywordsRejectReason    any        `json:"keywordsRejectReason"`
	ImageRejectReason       any        `json:"imageRejectReason"`
	CreditLine              any        `json:"creditLine"`
	PicWidth                any        `json:"picWidth"`
	PicHeight               any        `json:"picHeight"`
	PicSize                 any        `json:"picSize"`
	ColorType               any        `json:"colorType"`
	DateCameraShot          any        `json:"dateCameraShot"`
	Country                 any        `json:"country"`
	Province                any        `json:"province"`
	City                    any        `json:"city"`
	Location                any        `json:"location"`
	People                  any        `json:"people"`
	UploadTime              string     `json:"uploadTime"`
	OssYuantu               any        `json:"ossYuantu"`
	Oss176                  any        `json:"oss176"`
	Oss800                  any        `json:"oss800"`
	Oss800Watermark         any        `json:"oss800Watermark"`
	Oss400                  any        `json:"oss400"`
	OssEpsJpg               any        `json:"ossEpsJpg"`
	OssId7                  any        `json:"ossId7"`
	EditTime                string     `json:"editTime"`
	EditUserID              any        `json:"editUserId"`
	IsPostil                any        `json:"isPostil"`
	PostilTime              any        `json:"postilTime"`
	IsCopyright             any        `json:"isCopyright"`
	Copyright               any        `json:"copyright"`
	Maincolors              any        `json:"maincolors"`
	Orientation             any        `json:"orientation"`
	OnlineType              any        `json:"onlineType"`
	MinPrice                any        `json:"minPrice"`
	Price                   any        `json:"price"`
	Memo                    any        `json:"memo"`
	CreatedTime             string     `json:"createdTime"`
	CreatedBy               any        `json:"createdBy"`
	UpdatedTime             string     `json:"updatedTime"`
	UpdatedBy               any        `json:"updatedBy"`
	Caption                 string     `json:"caption"`
	KeywordsAudit           any        `json:"keywordsAudit"`
	ReadIptc                any        `json:"readIptc"`
	ReadFileName            any        `json:"readFileName"`
	OriginFileName          any        `json:"originFileName"`
	KeywordsSource          any        `json:"keywordsSource"`
	ResGroupIds             any        `json:"resGroupIds"`
	TopicIds                any        `json:"topicIds"`
	Dpi                     any        `json:"dpi"`
	ColorPattern            any        `json:"colorPattern"`
	ProcessingSoftware      any        `json:"processingSoftware"`
	SoftwareVersion         any        `json:"softwareVersion"`
	ExtType                 any        `json:"extType"`
	ExtOneCategory          any        `json:"extOneCategory"`
	ExtCategory             any        `json:"extCategory"`
	OriginID                any        `json:"originId"`
	ParentID                any        `json:"parentId"`
	ChildCount              any        `json:"childCount"`
	ResIDStr                string     `json:"res_id"`
	EqualwURL               string     `json:"equalw_url"`
	EqualhURL               string     `json:"equalh_url"`
	URL800                  string     `json:"url800"`
	PriceType               string     `json:"price_type"`
	ImgDate                 string     `json:"img_date"`
	Width                   int        `json:"width"`
	Height                  int        `json:"height"`
	Dlsize                  any        `json:"dlsize"`
	BrandName               string     `json:"brandName"`
	Restrict                int        `json:"restrict"`
	OneCategoryCn           any        `json:"oneCategoryCn"`
	ProductIds              string     `json:"productIds"`
	Source                  string     `json:"source"`
	GroupID                 any        `json:"groupId"`
	IstockCollection        any        `json:"istockCollection"`
	IsWhollyOwned           any        `json:"isWhollyOwned"`
	ImgAuditStatus          any        `json:"imgAuditStatus"`
	KeywordsAuditStatus     any        `json:"keywordsAuditStatus"`
	PeopleAuditStatus       any        `json:"peopleAuditStatus"`
	SecurityAuditStatus     any        `json:"securityAuditStatus"`
	ImgAuditOutsource       any        `json:"imgAuditOutsource"`
	KeywordsAuditOutsource  any        `json:"keywordsAuditOutsource"`
	PeopleAuditOutsource    any        `json:"peopleAuditOutsource"`
	SecurityAuditOutsource  any        `json:"securityAuditOutsource"`
	ImgOutsourceStatus      any        `json:"imgOutsourceStatus"`
	KeywordsOutsourceStatus any        `json:"keywordsOutsourceStatus"`
	PeopleOutsourceStatus   any        `json:"peopleOutsourceStatus"`
	SecurityOutsourceStatus any        `json:"securityOutsourceStatus"`
	DefineLabel             any        `json:"defineLabel"`
	PolitcsLabel            any        `json:"politcsLabel"`
	TerrorismLabel          any        `json:"terrorismLabel"`
	AbuseLabel              any        `json:"abuseLabel"`
	PornLabel               any        `json:"pornLabel"`
	ContrabandLabel         any        `json:"contrabandLabel"`
	CopyrightType           any        `json:"copyrightType"`
	CopyrightExclusive      string     `json:"copyrightExclusive"`
	UsageRestrictions       any        `json:"usageRestrictions"`
	Recommend               any        `json:"recommend"`
	Extend4                 any        `json:"extend4"`
	Flags                   any        `json:"flags"`
	AiGenerate              AiGenerate `json:"aiGenerate"`
	ProductMode             string     `json:"productMode"`
	DownLoad                bool       `json:"downLoad"`
}

// AI 生成信息
type AiGenerate struct {
	Extend       bool `json:"extend"`
	Face         bool `json:"face"`
	Background   bool `json:"background"`
	Illustration bool `json:"illustration"`
}

// 搜索结果 URL 信息
type ResUrl struct {
	Phrase    string `json:"phrase"`
	Transform string `json:"transform"`
	Page      string `json:"page"`
	Sort      string `json:"sort"`
	UUID      string `json:"uuid"`
}
