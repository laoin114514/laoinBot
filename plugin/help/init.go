package help

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"html/template"
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

	const htmlTpl = `<!doctype html>
<html lang="zh-CN">
<head>
<meta charset="utf-8" />
<meta name="viewport" content="width=device-width, initial-scale=1" />
<style>
  * { box-sizing: border-box; }
  body {
    margin: 0;
    padding: 24px;
    background: linear-gradient(135deg, #f8fafc, #eef2ff);
    font-family: "PingFang SC", "Microsoft YaHei", "Noto Sans SC", sans-serif;
    color: #0f172a;
  }
  .card {
    width: 1100px;
    margin: 0 auto;
    background: #ffffff;
    border-radius: 16px;
    box-shadow: 0 12px 28px rgba(15, 23, 42, .12);
    overflow: hidden;
  }
  .header {
    padding: 22px 26px;
    background: linear-gradient(90deg, #4f46e5, #7c3aed);
    color: white;
  }
  .title {
    font-size: 30px;
    font-weight: 700;
    margin: 0;
  }
  .sub {
    margin-top: 8px;
    opacity: .92;
    font-size: 16px;
  }
  .table-wrap { padding: 16px 18px 22px; }
  table {
    width: 100%;
    border-collapse: collapse;
    table-layout: fixed;
    background: #fff;
    border-radius: 10px;
    overflow: hidden;
  }
  th {
    text-align: left;
    background: #eef2ff;
    color: #312e81;
    font-weight: 700;
    padding: 12px 10px;
    border-bottom: 1px solid #e5e7eb;
    font-size: 14px;
  }
  td {
    padding: 12px 10px;
    border-bottom: 1px solid #f1f5f9;
    vertical-align: top;
    word-break: break-word;
    line-height: 1.5;
    font-size: 14px;
  }
  tr:nth-child(even) td { background: #fafbff; }
  .idx { width: 64px; text-align: center; }
  .cmd { width: 240px; font-weight: 700; color: #1d4ed8; }
  .desc { width: 330px; }
  .ex { width: 430px; color: #334155; }
</style>
</head>
<body>
  <div class="card">
    <div class="header">
      <h1 class="title">机器人帮助菜单</h1>
      <div class="sub">共 {{len .}} 条指令</div>
    </div>
    <div class="table-wrap">
      <table>
        <thead>
          <tr>
            <th class="idx">#</th>
            <th class="cmd">指令</th>
            <th class="desc">说明</th>
            <th class="ex">示例</th>
          </tr>
        </thead>
        <tbody>
          {{range .}}
          <tr>
            <td class="idx">{{.Index}}</td>
            <td class="cmd">{{.Order}}</td>
            <td class="desc">{{.Explain}}</td>
            <td class="ex">{{.Example}}</td>
          </tr>
          {{end}}
        </tbody>
      </table>
    </div>
  </div>
</body>
</html>`

	tpl, err := template.New("help").Parse(htmlTpl)
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
			chromedp.WindowSize(1200, 900),
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
		chromedp.WaitVisible("table", chromedp.ByQuery),
		chromedp.FullScreenshot(&imageBytes, 95),
	)
	if err != nil {
		return nil, err
	}
	return imageBytes, nil
}

func init() {
	nova.OnFullMatch("帮助").Handle(func(ctx *nova.Ctx) {
		helperList := HelpInstance.GetHelperList()
		if len(helperList) == 0 {
			ctx.Send("暂无帮助内容")
			return
		}

		imageBytes, err := renderHelpAsImage(helperList)
		if err != nil {
			for _, helper := range helperList {
				ctx.Send(helper.GetOrder() + " - " + helper.GetExplain() + " - " + helper.GetExample())
			}
			return
		}

		ctx.Send(message.ImageBytes(imageBytes).String())
	})
}
