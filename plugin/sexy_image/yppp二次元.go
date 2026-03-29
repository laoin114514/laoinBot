package sexy_image

import (
	nova "github.com/laoin114514/NovaBot"
	"github.com/laoin114514/NovaBot/message"
)

func init() {
	nova.OnFullMatch("来张图").Handle(func(ctx *nova.Ctx) {
		img, err := getOne2DImage()
		if err != nil {
			ctx.Send("图片获取失败: " + err.Error())
			return
		}
		ctx.Send("发送中...")
		ctx.Send(message.ImageBytes(img).String())
	})
}
