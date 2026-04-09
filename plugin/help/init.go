package help

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"html/template"
	"os"
	"sort"
	"time"

	"github.com/chromedp/chromedp"
	nova "github.com/laoin114514/NovaBot"
	"github.com/laoin114514/NovaBot/message"
)

type Help struct {
	order   string
	explain string
	example string
}

type Helpers struct {
	helpersMap map[string]Help
}

var HelpInstance = NewHelpers()

func NewHelpers() *Helpers {
	return &Helpers{
		helpersMap: make(map[string]Help),
	}
}

func (h *Helpers) SetHelper(order string, explain string, example string) error {
	if _, ok := h.helpersMap[order]; ok {
		return errors.New("helper already exists")
	}
	h.helpersMap[order] = Help{
		order:   order,
		explain: explain,
		example: example,
	}
	return nil
}

func (h *Helpers) GetHelper(order string) (Help, error) {
	if _, ok := h.helpersMap[order]; !ok {
		return Help{}, errors.New("helper not found")
	}
	return h.helpersMap[order], nil
}

func (h *Helpers) GetHelperList() []Help {
	helperList := make([]Help, 0, len(h.helpersMap))
	for _, helper := range h.helpersMap {
		helperList = append(helperList, helper)
	}
	return helperList
}

func (h *Help) SetExplain(explain string) error {
	if len(explain) > 100 {
		return errors.New("explain too long")
	}
	h.explain = explain
	return nil
}

func (h *Help) SetExample(example string) error {
	if len(example) > 100 {
		return errors.New("example too long")
	}
	h.example = example
	return nil
}

func (h *Help) GetExplain() string {
	return h.explain
}

func (h *Help) GetExample() string {
	return h.example
}

func (h *Help) GetOrder() string {
	return h.order
}

type helpRow struct {
	Index   int
	Order   string
	Explain string
	Example string
}

func renderHelpAsImage(helperList []Help) ([]byte, error) {
	sort.Slice(helperList, func(i, j int) bool {
		return helperList[i].order < helperList[j].order
	})

	rows := make([]helpRow, 0, len(helperList))
	for i, h := range helperList {
		rows = append(rows, helpRow{
			Index:   i + 1,
			Order:   h.GetOrder(),
			Explain: h.GetExplain(),
			Example: h.GetExample(),
		})
	}

	htmlTplBytes, err := os.ReadFile("plugin/help/template.html")
	if err != nil {
		return nil, err
	}

	tpl, err := template.New("help").Parse(string(htmlTplBytes))
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err = tpl.Execute(&buf, rows); err != nil {
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
			chromedp.WindowSize(1100, 900),
		)...,
	)
	defer cancelAlloc()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	timeoutCtx, cancelTimeout := context.WithTimeout(ctx, 15*time.Second)
	defer cancelTimeout()

	var imageBytes []byte
	err = chromedp.Run(timeoutCtx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(".help-root", chromedp.ByQuery),
		chromedp.Screenshot(".help-root", &imageBytes, chromedp.NodeVisible, chromedp.ByQuery),
	)
	if err != nil {
		return nil, err
	}
	return imageBytes, nil
}

func init() {
	HelpInstance.SetHelper("帮助", "查看帮助", "帮助")
	nova.OnFullMatch("帮助").Handle(func(ctx *nova.Ctx) {
		helperList := HelpInstance.GetHelperList()
		if len(helperList) == 0 {
			ctx.Send("暂无帮助内容")
			return
		}

		imageBytes, err := renderHelpAsImage(helperList)
		if err != nil {
			ctx.Send("生成帮助图片失败: " + err.Error())
			return
		}

		ctx.Send(message.ImageBytes(imageBytes).String())
	})
}
