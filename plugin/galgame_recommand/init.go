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
	help.HelpInstance.SetHelper("黄油推荐", "推荐黄油", "黄油推荐")
	nova.OnFullMatch("黄油推荐").Handle(func(ctx *nova.Ctx) {
		var galgame table.Galgame
		var count int64
		db.Db.Model(&table.Galgame{}).Count(&count)
		if count == 0 {
			ctx.SendChain(message.Text("没有黄油推荐"))
			return
		}
		id := rand.Intn(int(count))
		db.Db.Where("id = ?", id).Take(&galgame)

		imageBytes, err := renderGalgameAsImage(galgame)
		if err != nil {
			ctx.SendChain(message.Text(galgame.Code, "\n\n", galgame.Type, "\n\n", galgame.Tag, "\n\n", galgame.Name, "\n\n", galgame.Platform, "\n\n", galgame.Size, "\n\n", galgame.Version))
			return
		}

		ctx.Send(message.ImageBytes(imageBytes).String() + galgame.Url)
	})
}
