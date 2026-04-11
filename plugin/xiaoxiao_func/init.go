package xiaoxiaofunc

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"html/template"
	"laoinBot/plugin/help"
	"os"
	"strings"
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

type GoldPriceTemplateData struct {
	Mode     string
	SubTitle string
	GoldPriceData
}

func parseGoldPriceMode(argsText string) (string, string, error) {
	mode := strings.TrimSpace(argsText)
	switch mode {
	case "银行", "银行金条", "bank":
		return "bank", "银行投资金条", nil
	case "回收", "黄金回收", "recycle":
		return "recycle", "黄金回收", nil
	case "品牌", "品牌金店", "brand":
		return "brand", "品牌金店", nil
	default:
		return "", "", errors.New("参数错误，可用参数：\n银行\n回收\n品牌")
	}
}

func renderGoldPriceAsImage(data GoldPriceData, mode string, subTitle string) ([]byte, error) {
	htmlTplBytes, err := os.ReadFile("public/gold_price_template.html")
	if err != nil {
		return nil, err
	}

	tpl, err := template.New("gold_price").Parse(string(htmlTplBytes))
	if err != nil {
		return nil, err
	}

	tplData := GoldPriceTemplateData{
		Mode:          mode,
		SubTitle:      subTitle,
		GoldPriceData: data,
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
			chromedp.WindowSize(1320, 1500),
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
		chromedp.WaitVisible(".container", chromedp.ByQuery),
		chromedp.Screenshot(".container", &imageBytes, chromedp.NodeVisible, chromedp.ByQuery),
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

	help.HelpInstance.SetHelper("来一签", "来一签，多了就不准了awa", "来一签")
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

	help.HelpInstance.SetHelper("金价查询", "查询黄金价格（银行/回收/品牌）", "金价查询 银行")
	nova.OnPrefix("金价查询").Handle(func(ctx *nova.Ctx) {
		msgID := ctx.Event.MessageID
		argsText, _ := ctx.State["args"].(string)

		mode, subTitle, err := parseGoldPriceMode(argsText)
		if err != nil {
			ctx.Send(message.Reply(msgID).String() + err.Error())
			return
		}

		response, err := XiaoxiaoGetter.GetGoldPrice()
		if err != nil {
			ctx.Send(message.Reply(msgID).String() + "获取黄金价格失败: " + err.Error())
			return
		}

		imageBytes, err := renderGoldPriceAsImage(response.Data, mode, subTitle)
		if err != nil {
			ctx.Send(message.Reply(msgID).String() + "生成黄金价格图片失败: " + err.Error())
			return
		}

		ctx.Send(message.ImageBytes(imageBytes).String())
	})

	help.HelpInstance.SetHelper("写真图二", "来一张写真图", "写真图二 <参数>")
	nova.OnPrefix("写真图二").Handle(func(ctx *nova.Ctx) {
		argsText, _ := ctx.State["args"].(string)
		keyword := strings.TrimSpace(argsText)
		image, err := XiaoxiaoGetter.GetImage(keyword)
		if err != nil {
			ctx.Send(err.Error())
			return
		}
		ctx.Send(message.ImageBytes(image).String())
	})
}
