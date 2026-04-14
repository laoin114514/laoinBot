package viewchina

import (
	"fmt"
	"laoinBot/plugin/help"
	"math/rand"
	"strconv"
	"strings"

	nova "github.com/laoin114514/NovaBot"
	"github.com/laoin114514/NovaBot/message"
)

var viewChinaClient = NewViewChinaApi()

func init() {
	help.HelpInstance.SetHelper("图库搜索", "搜索图库中的图片", "图库搜索 <关键词> [页码] [排序(可选，fresh或hot)]")
	nova.OnPrefix("图库搜索").SetBlock(true).Handle(func(ctx *nova.Ctx) {
		argsText, _ := ctx.State["args"].(string)
		args := strings.Fields(strings.TrimSpace(argsText))
		if len(args) == 0 {
			ctx.Send("参数不能为空，格式：图库搜索 <关键词> [页码] [排序(可选，fresh或hot)]")
			return
		}

		phrase := args[0]
		page := 1
		sort := ""
		var err error

		if len(args) >= 2 {
			page, err = strconv.Atoi(args[1])
			if err != nil || page <= 0 {
				ctx.Send("页码必须是大于 0 的整数")
				return
			}
		}
		if len(args) >= 3 {
			sort = strings.ToLower(args[2])
		}
		if sort != "fresh" && sort != "hot" {
			sort = ""
		}

		resp, err := viewChinaClient.SearchImages(phrase, page, sort)
		if err != nil {
			ctx.Send("搜索失败: " + err.Error())
			return
		}
		if len(resp.List) == 0 {
			ctx.Send("没有搜索到图片")
			return
		}

		item := resp.List[rand.Intn(len(resp.List))]
		imageURL := item.URL800
		if imageURL == "" {
			imageURL = item.EqualwURL
		}
		if imageURL == "" {
			imageURL = item.EqualhURL
		}
		if imageURL == "" {
			ctx.Send("图片链接为空，请重试")
			return
		}

		imgBytes, err := viewChinaClient.DownloadImage(imageURL)
		if err != nil {
			ctx.Send("下载图片失败: " + err.Error())
			return
		}

		caption := item.Title
		if caption == "" {
			caption = item.Caption
		}
		if caption == "" {
			caption = "图库随机图"
		}
		ctx.Send(fmt.Sprintf("%s\n关键词: %s | 页码: %d | 排序: %s", caption, phrase, page, sort) + message.ImageBytes(imgBytes).String())
	})
}
