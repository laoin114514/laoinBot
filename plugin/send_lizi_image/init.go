package send_lizi_image

import (
	"errors"
	"strings"

	nova "github.com/laoin114514/NovaBot"
	"github.com/laoin114514/NovaBot/message"
)

func init() {
	nova.OnCommand("随机图").Handle(func(ctx *nova.Ctx) {
		argsText, _ := ctx.State["args"].(string)
		param, err := handleParams(ctx, argsText)
		if err != nil {
			ctx.Send(err.Error())
			return
		}
		img, err := LiziGetter.GetOneImage(HasRepo[param])
		if err != nil {
			ctx.Send("图片获取失败: " + err.Error())
			return
		}
		ctx.Send(message.ImageBytes(img).String())
	})
}

func handleParams(ctx *nova.Ctx, argsText string) (string, error) {
	args := strings.Split(argsText, " ")
	if len(args) == 0 {
		return "", errors.New("参数不能为空")
	}
	param := args[0]
	if _, ok := HasRepo[param]; !ok {
		repoList := make([]string, 0, len(HasRepo))
		for repo := range HasRepo {
			repoList = append(repoList, repo)
		}
		return "", errors.New("参数错误: " + param + "\n可用参数: \n" + strings.Join(repoList, "\n"))
	}
	return param, nil
}
