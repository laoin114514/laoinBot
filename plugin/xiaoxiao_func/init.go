package xiaoxiaofunc

import (
	"laoinBot/plugin/help"

	nova "github.com/laoin114514/NovaBot"
	"github.com/laoin114514/NovaBot/message"
)

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
}
