package cangmiao_func

import (
	"errors"
	"laoinBot/plugin/help"
	"os"
	"strings"

	"github.com/go-resty/resty/v2"
	nova "github.com/laoin114514/NovaBot"
	"github.com/laoin114514/NovaBot/message"
)

func init() {
	help.HelpInstance.SetHelper("写真图", "发送写真图", "写真图 <参数>")
	nova.OnPrefix("写真图").Handle(func(ctx *nova.Ctx) {
		argsText, _ := ctx.State["args"].(string)
		if argsText == "瑟瑟" {
			ctx.Send("瑟瑟只有laoin能用")
			return
		}
		param, err := handleParams(ctx, argsText)
		if err != nil {
			ctx.Send(err.Error())
			return
		}
		img, err := CangMiaoGetter.GetImageURL(param)
		if err != nil {
			ctx.Send("图片获取失败" + err.Error())
			return
		}
		saveImage(img.Img)
		msgID := ctx.Send(message.ImageBytes(img.ToByte()).String())
		if msgID.ID() == 0 {
			ctx.Send("图片发送失败")
			return
		}
	})

	help.HelpInstance.SetHelper("查询号码", "查询号码归属地", "查询号码 <参数>")
	nova.OnPrefix("查询号码").Handle(func(ctx *nova.Ctx) {
		msgID := ctx.Event.MessageID
		argsText, _ := ctx.State["args"].(string)
		address, err := CangMiaoGetter.GetPhoneAdress(argsText)
		if err != nil {
			ctx.Send(message.Reply(msgID).String() + "获取号码归属地失败: " + err.Error())
			return
		}
		ctx.Send(message.Reply(msgID).String() + address)
	})

	help.HelpInstance.SetHelper("查询IP", "查询IP归属地", "查询IP <参数>")
	nova.OnPrefix("查询IP").Handle(func(ctx *nova.Ctx) {
		msgID := ctx.Event.MessageID
		argsText, _ := ctx.State["args"].(string)
		address, err := CangMiaoGetter.GetIPAdress(argsText)
		if err != nil {
			ctx.Send(message.Reply(msgID).String() + "获取IP归属地失败: " + err.Error())
			return
		}
		ctx.Send(message.Reply(msgID).String() + address)
	})
}

func handleParams(ctx *nova.Ctx, argsText string) (string, error) {
	args := strings.Split(argsText, " ")
	if len(args) == 0 {
		return "", errors.New("参数不能为空")
	}
	param := args[0]
	if _, ok := ImageRepos[param]; !ok {
		repoList := make([]string, 0, len(ImageRepos))
		for repo := range ImageRepos {
			repoList = append(repoList, repo)
		}
		return "", errors.New("参数错误: " + param + "\n可用参数: \n" + strings.Join(repoList, "\n"))
	}
	return ImageRepos[param], nil
}

func saveImage(img string) error {
	c := resty.New()
	resp, err := c.R().
		Get(img)
	if err != nil {
		return err
	}
	if resp.IsError() {
		return errors.New("图片获取失败: " + resp.Status())
	}
	file, err := os.Create("image.jpg")
	if err != nil {
		return err
	}
	defer file.Close()
	file.Write(resp.Body())
	return nil
}
