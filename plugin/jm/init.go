package jm

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"html/template"
	"laoinBot/config"
	"laoinBot/plugin/help"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	nova "github.com/laoin114514/NovaBot"
	"github.com/laoin114514/NovaBot/message"
	"github.com/laoin114514/jmapi"
)

type jmPosterTemplateData struct {
	ID           string
	Name         string
	Author       []string
	PageCount    int
	UpdateDate   string
	Tags         []string
	Description  string
	Views        string
	Likes        string
	CommentCount int
	CoverSrc     string
}

func renderJMPosterAsImage(album *jmapi.AlbumDetail, coverBytes []byte) ([]byte, error) {
	if album == nil {
		return nil, fmt.Errorf("album is nil")
	}

	tags := album.Tags
	if len(tags) == 0 {
		tags = []string{"暂无标签"}
	}
	authors := album.Author
	if len(authors) == 0 {
		authors = []string{"未知"}
	}
	updateDate := "未知"
	if album.Raw != nil {
		if v, ok := album.Raw["update_date"].(string); ok && strings.TrimSpace(v) != "" {
			updateDate = v
		} else if v, ok := album.Raw["updateDate"].(string); ok && strings.TrimSpace(v) != "" {
			updateDate = v
		}
	}

	// 优先写到本地文件，让 HTML 用相对路径引用（Windows 下更稳）
	// 失败时再回退到 data URI
	tempDir, err := os.MkdirTemp("", "jm_poster_*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tempDir)

	coverSrc := "data:image/gif;base64,R0lGODlhAQABAIAAAAAAAP///ywAAAAAAQABAAACAUwAOw=="
	coverFilename := "cover.jpg"
	if len(coverBytes) != 0 {
		ct := http.DetectContentType(coverBytes)
		switch ct {
		case "image/png":
			coverFilename = "cover.png"
		case "image/gif":
			coverFilename = "cover.gif"
		case "image/webp":
			coverFilename = "cover.webp"
		default:
			coverFilename = "cover.jpg"
		}

		coverPath := filepath.Join(tempDir, coverFilename)
		if err := os.WriteFile(coverPath, coverBytes, 0o644); err == nil {
			// 与 HTML 同目录，直接用相对路径
			coverSrc = coverFilename
		} else {
			// 写文件失败，回退 data URI
			if strings.HasPrefix(ct, "image/") {
				coverSrc = "data:" + ct + ";base64," + base64.StdEncoding.EncodeToString(coverBytes)
			} else {
				coverSrc = "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(coverBytes)
			}
		}
	}

	tplData := jmPosterTemplateData{
		ID:           album.ID,
		Name:         album.Name,
		Author:       authors,
		PageCount:    album.PageCount,
		UpdateDate:   updateDate,
		Tags:         tags,
		Description:  album.Description,
		Views:        album.Views,
		Likes:        album.Likes,
		CommentCount: album.CommentCount,
		CoverSrc:     coverSrc,
	}

	htmlTplBytes, err := os.ReadFile("public/jm_post_template.html")
	if err != nil {
		return nil, err
	}

	tpl, err := template.New("jm_poster").Parse(string(htmlTplBytes))
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err = tpl.Execute(&buf, tplData); err != nil {
		return nil, err
	}

	// 写出 HTML 文件，用 file:/// 打开，从而能引用同目录 cover.*
	htmlPath := filepath.Join(tempDir, "index.html")
	if err := os.WriteFile(htmlPath, buf.Bytes(), 0o644); err != nil {
		return nil, err
	}
	url := "file:///" + filepath.ToSlash(htmlPath)

	allocCtx, cancelAlloc := chromedp.NewExecAllocator(
		context.Background(),
		append(
			chromedp.DefaultExecAllocatorOptions[:],
			chromedp.Flag("headless", true),
			chromedp.Flag("no-sandbox", true),
			chromedp.Flag("disable-gpu", true),
			chromedp.Flag("allow-file-access-from-files", true),
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
		chromedp.WaitVisible(".card", chromedp.ByQuery),
		// 等封面图真正加载完成（否则可能截图到空白/占位）
		chromedp.ActionFunc(func(ctx context.Context) error {
			deadline := time.Now().Add(10 * time.Second)
			for time.Now().Before(deadline) {
				var ok bool
				err := chromedp.EvaluateAsDevTools(`(() => {
					const img = document.querySelector('.cover');
					if (!img) return false;
					return img.complete && img.naturalWidth > 0;
				})()`, &ok).Do(ctx)
				if err == nil && ok {
					return nil
				}
				time.Sleep(120 * time.Millisecond)
			}
			// 超时也继续截图（至少不会卡死）
			return nil
		}),
		chromedp.Screenshot(".card", &imageBytes, chromedp.NodeVisible, chromedp.ByQuery),
	)
	if err != nil {
		return nil, err
	}
	return imageBytes, nil
}

func init() {
	help.HelpInstance.SetHelper("本子推荐", "随机本子", "本子推荐 <关键词> <排序> <时间范围> <分类> <子分类>")
	nova.OnPrefix("本子推荐").Handle(func(ctx *nova.Ctx) {
		argsText, _ := ctx.State["args"].(string)
		argsText = strings.TrimSpace(argsText)

		// 默认参数（与 jmapi.SearchSite 入参顺序一致，page 由程序随机决定）
		searchQuery := "纯爱"
		orderBy := jmapi.OrderByLatest
		timeRange := jmapi.TimeAll
		category := ""
		subCategory := ""

		// 用户参数：按顺序覆盖
		// 本子 <关键词> <排序> <时间范围> <分类> <子分类>
		if argsText != "" {
			parts := strings.Fields(argsText)
			if len(parts) >= 1 {
				searchQuery = parts[0]
			}
			if len(parts) >= 2 {
				orderBy = parts[1]
			}
			if len(parts) >= 3 {
				timeRange = parts[2]
			}
			if len(parts) >= 4 {
				category = parts[3]
			}
			if len(parts) >= 5 {
				subCategory = parts[4]
			}
		}

		ctx.Send("正在获取本子...")

		// 第一次请求：拿总数，用于计算总页数（每页按 80 条）
		firstPage, err := config.JMClient.SearchSite(searchQuery, 1, orderBy, timeRange, category, subCategory)
		if err != nil {
			ctx.Send(err.Error())
			return
		}
		if firstPage.Total <= 0 || len(firstPage.Items) == 0 {
			ctx.Send("没有搜到结果")
			return
		}

		totalPages := (firstPage.Total + 80 - 1) / 80
		if totalPages <= 0 {
			totalPages = 1
		}
		randomPage := rand.Intn(totalPages) + 1

		pageList := firstPage
		if randomPage != 1 {
			pageList, err = config.JMClient.SearchSite(searchQuery, randomPage, orderBy, timeRange, category, subCategory)
			if err != nil {
				ctx.Send(err.Error())
				return
			}
		}
		if len(pageList.Items) == 0 {
			ctx.Send("随机页没有结果，请稍后重试")
			return
		}

		id := pageList.Items[rand.Intn(len(pageList.Items))].ID
		album, err := config.JMClient.GetAlbumDetail(id)
		if err != nil {
			ctx.Send(err.Error())
			return
		}
		cover, err := config.JMClient.DownloadAlbumCover(id)
		if err != nil {
			ctx.Send(err.Error())
			return
		}

		poster, err := renderJMPosterAsImage(album, cover)
		if err != nil {
			// 渲染失败时回退到原始文本 + 封面图，避免影响功能可用性
			ctx.Send(err.Error())
			// ctx.Send(fmt.Sprintf("本子名称: %s\n本子描述: %s\n本子评论数: %d", album.Name, album.Description, album.CommentCount) + message.ImageBytes(cover).String())
			return
		}
		ctx.Send(message.ImageBytes(poster).String())
	})
}
