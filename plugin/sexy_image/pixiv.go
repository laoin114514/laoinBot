package sexy_image

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"strings"

	nova "github.com/laoin114514/NovaBot"
	"github.com/laoin114514/NovaBot/extension/shell"
	"github.com/laoin114514/NovaBot/message"
)

type PixivApiParams struct {
	Num           int
	Pid           []int64
	Uid           []int64
	RawPid        string
	RawUid        string
	Author        string
	Proxy         string
	AiType        int
	R18Type       int
	DateAfter     int64
	DateBefore    int64
	SizeList      []string
	RawSizeList   string
	ImageSizeType int
}

func init() {
	nova.OnCommand("涩图").Handle(func(ctx *nova.Ctx) {
		if nova.OnlyGroup(ctx) {
			ctx.Send("请在私聊中使用awa")
			return
		}
		params := PixivApiParams{}

		if err := handleParams(ctx, params); err != nil {
			ctx.Send("参数处理失败：" + err.Error())
			return
		}

		if params.R18Type == 1 && !nova.SuperUserPermission(ctx) {
			ctx.Send("只有laoin能用r18哦！" + message.Face(122).String())
			return
		}

		msgID := ctx.Event.MessageID
		searchParams := SearchParams{
			Num:           params.Num,
			Pid:           params.Pid,
			Uid:           params.Uid,
			Author:        params.Author,
			Proxy:         params.Proxy,
			AiType:        params.AiType,
			R18Type:       params.R18Type,
			DateAfter:     params.DateAfter,
			DateBefore:    params.DateBefore,
			SizeList:      parseStringList(params.RawSizeList),
			ImageSizeType: params.ImageSizeType,
		}
		url, err := getOneSexyImage(searchParams)
		if err != nil {
			ctx.Send("获取涩图失败：" + err.Error() + message.Reply(msgID).String())
			return
		}
		ctx.Send("获取成功，正在发送...")
		code := ctx.Send(message.Image(url).String())
		if code.ID() == 0 {
			ctx.Send("发送失败：可能图太涩了awa" + message.Reply(msgID).String())
			return
		}
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

func handleParams(ctx *nova.Ctx, params PixivApiParams) error {
	fset := flag.NewFlagSet("涩图", flag.ContinueOnError)

	fset.IntVar(&params.Num, "num", 1, "")
	fset.StringVar(&params.RawPid, "pid", "", "")
	fset.StringVar(&params.RawUid, "uid", "", "")
	fset.StringVar(&params.Author, "author", "", "")
	fset.StringVar(&params.Proxy, "proxy", "i.pixiv.re", "")
	fset.IntVar(&params.AiType, "ai", 0, "")
	fset.IntVar(&params.R18Type, "r18", 0, "")
	fset.Int64Var(&params.DateAfter, "dateAfter", 0, "")
	fset.Int64Var(&params.DateBefore, "dateBefore", 0, "")
	fset.StringVar(&params.RawSizeList, "sizeList", "original", "")
	fset.IntVar(&params.ImageSizeType, "imageSizeType", 0, "")

	argsText, _ := ctx.State["args"].(string)
	fmt.Println(argsText)
	arguments := shell.Parse(argsText)
	fmt.Println(arguments)
	if err := fset.Parse(arguments); err != nil {
		return errors.New("参数解析失败: " + err.Error())
	}

	pidList, err := parseInt64List(params.RawPid)
	if err != nil {
		return errors.New("pid 参数错误: " + err.Error())
	}
	params.Pid = pidList
	uidList, err := parseInt64List(params.RawUid)
	if err != nil {
		return errors.New("uid 参数错误: " + err.Error())
	}
	params.Uid = uidList
	if params.Num < 1 || params.Num > 20 {
		return errors.New("num 参数范围是 1~20")
	}
	if params.AiType < 0 || params.AiType > 2 {
		return errors.New("aiType 参数范围是 0~2")
	}
	if params.R18Type < 0 || params.R18Type > 1 {
		return errors.New("r18Type 参数范围是 0~1")
	}
	if params.ImageSizeType < 0 || params.ImageSizeType > 3 {
		return errors.New("imageSizeType 参数范围是 0~3")
	}
	return nil
}
