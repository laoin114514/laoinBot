package sendlike

import (
	"fmt"
	"math/rand"

	nova "github.com/laoin114514/NovaBot"
	"github.com/laoin114514/NovaBot/message"
)

func init() {
	nova.OnFullMatch("赞我", nova.OnlyGroup).Handle(func(ctx *nova.Ctx) {
		userID := ctx.Event.UserID
		msgID := ctx.Event.MessageID
		count := 0
		for {
			_, err := ctx.SendLike(ctx.Event.UserID, 10)
			if err != nil {
				break
			}
			count++
		}

		if count == 0 {
			ctx.Send(fmt.Sprintf("%s%s点赞失败,点赞到达上限了%s", message.At(userID), message.Reply(msgID), message.Face(rand.Intn(100)).String()))
			return
		}

		ctx.Send(fmt.Sprintf("%s%s点赞成功，已点赞%d次%s", message.At(userID), message.Reply(msgID), count*10, message.Face(rand.Intn(100)).String()))
	})
}
