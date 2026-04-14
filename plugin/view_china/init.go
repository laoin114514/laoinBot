package viewchina

import (
	"laoinBot/plugin/help"
	"math/rand"
	"strconv"
	"strings"

	nova "github.com/laoin114514/NovaBot"
	"github.com/laoin114514/NovaBot/message"
)

var viewChinaClient = NewViewChinaApi()

func init() {
	help.HelpInstance.SetHelper("图库搜索", "搜索图库中的图片", "图库搜索 <关键词> [-page 页码] [-total 数量] [-sort fresh|hot]")
	nova.OnPrefix("图库搜索", nova.SuperUserPermission).SetBlock(true).Handle(func(ctx *nova.Ctx) {
		argsText, _ := ctx.State["args"].(string)
		args := strings.Fields(strings.TrimSpace(argsText))
		if len(args) == 0 {
			ctx.Send("参数不能为空，格式：图库搜索 <关键词> [-page 页码] [-total 数量] [-sort fresh|hot]")
			return
		}

		page := 1
		sort := ""
		total := 10
		phraseEnd := 0
		for phraseEnd < len(args) && !strings.HasPrefix(args[phraseEnd], "-") {
			phraseEnd++
		}
		if phraseEnd == 0 {
			ctx.Send("关键词不能为空，格式：图库搜索 <关键词> [-page 页码] [-total 数量] [-sort fresh|hot]")
			return
		}
		phrase := strings.Join(args[:phraseEnd], " ")

		for i := phraseEnd; i < len(args); i += 2 {
			key := strings.ToLower(strings.TrimPrefix(args[i], "-"))
			if key == "" {
				ctx.Send("参数格式错误，示例：-page 2 -total 5")
				return
			}
			if i+1 >= len(args) {
				ctx.Send("参数缺少值，请使用 -变量 值")
				return
			}
			val := args[i+1]
			switch key {
			case "page", "p":
				n, err := strconv.Atoi(val)
				if err != nil || n <= 0 {
					ctx.Send("页码必须是大于 0 的整数")
					return
				}
				page = n
			case "total", "t", "size":
				n, err := strconv.Atoi(val)
				if err != nil || n <= 0 {
					ctx.Send("数量必须是大于 0 的整数")
					return
				}
				total = n
			case "sort", "s":
				v := strings.ToLower(val)
				if v != "fresh" && v != "hot" {
					ctx.Send("排序只支持 fresh 或 hot")
					return
				}
				sort = v
			default:
				ctx.Send("不支持的参数: -" + key)
				return
			}
		}
		if total > 50 {
			total = 50
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

		if total > len(resp.List) {
			total = len(resp.List)
		}
		if total <= 0 {
			ctx.Send("可用图片数量不足")
			return
		}

		nodes := make(message.Message, 0, total)
		for _, idx := range rand.Perm(len(resp.List)) {
			if len(nodes) >= total {
				break
			}

			item := resp.List[idx]
			imageURL := item.URL800
			if imageURL == "" {
				imageURL = item.EqualwURL
			}
			if imageURL == "" {
				imageURL = item.EqualhURL
			}
			if imageURL == "" {
				continue
			}

			imgBytes, err := viewChinaClient.DownloadImage(imageURL)
			if err != nil {
				continue
			}

			caption := item.Title
			if caption == "" {
				caption = item.Caption
			}
			if caption == "" {
				caption = "图库随机图"
			}

			content := message.Message{
				message.ImageBytes(imgBytes),
			}
			nodes = append(nodes, message.CustomNode("图库搜索", ctx.Event.SelfID, content))
		}

		if len(nodes) == 0 {
			ctx.Send("图片下载失败，请稍后重试")
			return
		}
		ctx.Send(nodes)
	})
}
