package galgamerecommand

import (
	"bytes"
	"context"
	"encoding/base64"
	"html/template"
	"laoinBot/config/db"
	"laoinBot/config/db/table"
	"laoinBot/plugin/help"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	nova "github.com/laoin114514/NovaBot"
	"github.com/laoin114514/NovaBot/message"
)

type galgameTemplateData struct {
	Code         string
	Name         string
	Type         string
	Tags         []string
	Platform     string
	Size         string
	Version      string
	Title        string
	Url          string
	Introduction string
}

func renderGalgameAsImage(galgame table.Galgame) ([]byte, error) {
	tags := make([]string, 0)
	for t := range strings.SplitSeq(galgame.Tag, ",") {
		t = strings.TrimSpace(t)
		if t != "" {
			tags = append(tags, t)
		}
	}
	if len(tags) == 0 {
		tags = append(tags, "暂无标签")
	}
	if galgame.Version == "" {
		galgame.Version = "Ver0.0.0"
	}
	tplData := galgameTemplateData{
		Code:         galgame.Code,
		Name:         galgame.Name,
		Type:         galgame.Type,
		Tags:         tags,
		Platform:     galgame.Platform,
		Size:         galgame.Size,
		Version:      galgame.Version,
		Title:        galgame.Title,
		Url:          galgame.Url,
		Introduction: galgame.Introduction,
	}

	htmlTplBytes, err := os.ReadFile("public/galgame_template.html")
	if err != nil {
		return nil, err
	}

	tpl, err := template.New("galgame").Parse(string(htmlTplBytes))
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err = tpl.Execute(&buf, tplData); err != nil {
		return nil, err
	}

	htmlData := base64.StdEncoding.EncodeToString(buf.Bytes())
	url := "data:text/html;base64," + htmlData

	allocCtx, cancelAlloc := chromedp.NewExecAllocator(
		context.Background(),
		append(
			chromedp.DefaultExecAllocatorOptions[:],
			chromedp.Flag("headless", true),
			chromedp.Flag("no-sandbox", true),
			chromedp.Flag("disable-gpu", true),
			chromedp.WindowSize(900, 1600),
		)...,
	)
	defer cancelAlloc()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	timeoutCtx, cancelTimeout := context.WithTimeout(ctx, 20*time.Second)
	defer cancelTimeout()

	var imageBytes []byte
	err = chromedp.Run(timeoutCtx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(".poster", chromedp.ByQuery),
		chromedp.Screenshot(".poster", &imageBytes, chromedp.NodeVisible, chromedp.ByQuery),
	)
	if err != nil {
		return nil, err
	}

	return imageBytes, nil
}

func init() {
	help.HelpInstance.SetHelper("黄油推荐", "按参数筛选推荐黄油", "黄油推荐 -tag 日系 -platform PC")
	nova.OnPrefix("黄油推荐").Handle(func(ctx *nova.Ctx) {
		argsText, _ := ctx.State["args"].(string)
		argsText = strings.TrimSpace(argsText)

		helpText := "可用参数：\n" +
			"-name 游戏名（模糊匹配）\n" +
			"-tag 标签（可多词）\n" +
			"-platform 平台\n" +
			"-code 编号\n" +
			"-type 类型\n\n" +
			"示例：\n" +
			"黄油推荐 -tag 日系 动态 -platform PC\n" +
			"黄油推荐 -name 透明人间\n" +
			"黄油推荐 help"

		if argsText == "help" {
			ctx.SendChain(message.Text(helpText))
			return
		}

		query := db.Db.Model(&table.Galgame{})
		if argsText != "" {
			tokens := strings.Fields(argsText)
			currentFlag := ""
			flagValues := map[string][]string{}

			for _, tok := range tokens {
				if strings.HasPrefix(tok, "-") {
					switch tok {
					case "-name", "-tag", "-platform", "-code", "-type":
						currentFlag = tok
						if _, ok := flagValues[currentFlag]; !ok {
							flagValues[currentFlag] = []string{}
						}
					default:
						ctx.SendChain(message.Text("不支持的参数: ", tok, "\n\n", helpText))
						return
					}
					continue
				}

				if currentFlag == "" {
					ctx.SendChain(message.Text("参数格式错误，请使用 -name/-tag/-platform/-code/-type\n\n", helpText))
					return
				}
				flagValues[currentFlag] = append(flagValues[currentFlag], tok)
			}

			for flag, values := range flagValues {
				if len(values) == 0 {
					ctx.SendChain(message.Text("参数缺少值: ", flag, "\n\n", helpText))
					return
				}
				value := strings.Join(values, " ")
				switch flag {
				case "-name":
					query = query.Where("name LIKE ?", "%"+value+"%")
				case "-platform":
					query = query.Where("platform LIKE ?", "%"+value+"%")
				case "-code":
					query = query.Where("code LIKE ?", "%"+value+"%")
				case "-type":
					query = query.Where("type LIKE ?", "%"+value+"%")
				case "-tag":
					for _, tag := range values {
						query = query.Where("tag LIKE ?", "%"+tag+"%")
					}
				}
			}
		}

		var count int64
		query.Count(&count)
		if count == 0 {
			ctx.SendChain(message.Text("没有匹配条件的黄油推荐"))
			return
		}

		var galgame table.Galgame
		offset := rand.Intn(int(count))
		err := query.Offset(offset).Limit(1).Take(&galgame).Error
		if err != nil {
			ctx.SendChain(message.Text("获取黄油推荐失败: ", err.Error()))
			return
		}

		imageBytes, err := renderGalgameAsImage(galgame)
		if err != nil {
			ctx.SendChain(message.Text("没有黄油推荐"))
			return
		}

		ctx.Send(message.ImageBytes(imageBytes).String())
	})
}
