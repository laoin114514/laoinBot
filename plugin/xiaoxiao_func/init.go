package xiaoxiaofunc

import (
	"bytes"
	"context"
	"encoding/base64"
	"html/template"
	"laoinBot/plugin/help"
	"os"
	"time"

	"github.com/chromedp/chromedp"
	nova "github.com/laoin114514/NovaBot"
	"github.com/laoin114514/NovaBot/message"
)

func renderLotteryAsImage(data LotteryData) ([]byte, error) {
	htmlTplBytes, err := os.ReadFile("public/lottery_template.html")
	if err != nil {
		return nil, err
	}

	tpl, err := template.New("lottery").Parse(string(htmlTplBytes))
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err = tpl.Execute(&buf, data); err != nil {
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
			chromedp.WindowSize(900, 1200),
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
		chromedp.WaitVisible(".lottery-card", chromedp.ByQuery),
		chromedp.Screenshot(".lottery-card", &imageBytes, chromedp.NodeVisible, chromedp.ByQuery),
	)
	if err != nil {
		return nil, err
	}

	return imageBytes, nil
}

func init() {
	help.HelpInstance.SetHelper("历史上的今天", "查询历史上的今天", "历史上的今天")
	nova.OnFullMatch("历史上的今天").Handle(func(ctx *nova.Ctx) {
		msgID := ctx.Event.MessageID
		image, err := XiaoxiaoGetter.GetHistoryToday()
		if err != nil {
			ctx.Send(message.Reply(msgID).String() + "获取历史上的今天失败: " + err.Error())
			return
		}
		ctx.Send(message.ImageBytes(image))
	})

	help.HelpInstance.SetHelper("来一签", "来一签", "来一签")
	nova.OnFullMatch("来一签").Handle(func(ctx *nova.Ctx) {
		msgID := ctx.Event.MessageID
		response, err := XiaoxiaoGetter.GetLottery()
		if err != nil {
			ctx.Send(message.Reply(msgID).String() + "获取签文失败: " + err.Error())
			return
		}

		imageBytes, err := renderLotteryAsImage(response.Data)
		if err != nil {
			ctx.Send(message.Reply(msgID).String() + "生成签文图片失败: " + err.Error())
			return
		}

		ctx.Send(message.ImageBytes(imageBytes).String())
	})
}
