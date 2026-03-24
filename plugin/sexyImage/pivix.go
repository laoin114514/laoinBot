package sexyimage

import (
	"flag"
	"strconv"
	"strings"
	"time"

	nova "github.com/laoin114514/NovaBot"
	"github.com/laoin114514/NovaBot/extension/shell"
	"github.com/laoin114514/NovaBot/message"
)

func init() {
	nova.OnCommand("涩图", nova.OnlyGroup).Handle(func(ctx *nova.Ctx) {
		fset := flag.NewFlagSet("涩图", flag.ContinueOnError)
		var (
			num           int
			pidRaw        string
			uidRaw        string
			author        string
			proxy         string
			aiType        int
			r18Type       int
			dateAfter     int64
			dateBefore    int64
			sizeListRaw   string
			imageSizeType int
		)

		fset.IntVar(&num, "num", 1, "")
		fset.StringVar(&pidRaw, "pid", "", "")
		fset.StringVar(&uidRaw, "uid", "", "")
		fset.StringVar(&author, "author", "", "")
		fset.StringVar(&proxy, "proxy", "i.pixiv.re", "")
		fset.IntVar(&aiType, "ai", 0, "")
		fset.IntVar(&r18Type, "r18", 0, "")
		fset.Int64Var(&dateAfter, "dateAfter", 0, "")
		fset.Int64Var(&dateBefore, "dateBefore", 0, "")
		fset.StringVar(&sizeListRaw, "sizeList", "original", "")
		fset.IntVar(&imageSizeType, "imageSizeType", 0, "")

		argsText, _ := ctx.State["args"].(string)
		arguments := shell.Parse(argsText)
		if err := fset.Parse(arguments); err != nil {
			ctx.Send("参数解析失败，请使用示例：涩图 -num 1 -r18 1 -sizeList original")
			return
		}

		pidList, err := parseInt64List(pidRaw)
		if err != nil {
			ctx.Send("pid 参数错误: " + err.Error())
			return
		}
		uidList, err := parseInt64List(uidRaw)
		if err != nil {
			ctx.Send("uid 参数错误: " + err.Error())
			return
		}

		if num < 1 || num > 20 {
			ctx.Send("num 参数范围是 1~20")
			return
		}
		if aiType < 0 || aiType > 2 {
			ctx.Send("aiType 参数范围是 0~2")
			return
		}
		if r18Type == 1 && ctx.Event.MessageType == "group" && ctx.Event.UserID != 2908451607 {
			ctx.Send("只有laoin能用r18哦！" + message.Face(122).String())
			return
		}
		if r18Type < 0 || r18Type > 1 {
			ctx.Send("r18Type 参数范围是 0~1")
			return
		}
		if imageSizeType < 0 || imageSizeType > 3 {
			ctx.Send("imageSizeType 参数范围是 0~3")
			return
		}

		params := SearchParams{
			Num:           num,
			Pid:           pidList,
			Uid:           uidList,
			Author:        author,
			Proxy:         proxy,
			AiType:        aiType,
			R18Type:       r18Type,
			DateAfter:     dateAfter,
			DateBefore:    dateBefore,
			SizeList:      parseStringList(sizeListRaw),
			ImageSizeType: imageSizeType,
		}
		url, err := getOneSexyImage(params)
		if err != nil {
			ctx.Send("获取涩图失败：" + err.Error())
			return
		}
		ctx.Send("获取成功，正在发送...")
		code := ctx.Send(message.Image(url).String())
		if code.ID() == 0 {
			ctx.Send("发送失败：可能图太涩了awa")
			return
		}
		time.Sleep(5 * time.Second)
		ctx.DeleteMessage(code.ID())
	})
}

func parseInt64List(raw string) ([]int64, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, nil
	}
	parts := strings.Split(raw, ",")
	result := make([]int64, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		v, err := strconv.ParseInt(part, 10, 64)
		if err != nil {
			return nil, err
		}
		result = append(result, v)
	}
	return result, nil
}

func parseStringList(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		result = append(result, part)
	}
	if len(result) == 0 {
		return nil
	}
	return result
}
