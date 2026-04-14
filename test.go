package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/mozillazg/go-pinyin"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	c := resty.New()
	var result map[string]any
	var ch string
	fmt.Scan(&ch)
	pin := pinyin.Slug(ch, pinyin.Args{})
	if pin == "" {
		pin = ch
	}
	params := map[string]string{
		"phrase":    ch,
		"transform": pin,
	}
	sec, err := GenerateSec(params)
	_, err = c.R().
		SetResult(&result).
		SetHeader("sec", sec).
		SetQueryParams(params).
		Get("https://www.vcg.com/api/common/searchAllImage")
	if err != nil {
		log.Fatalf("请求失败: %v", err)
	}
	os.Remove("public/result_cache.json")
	f, err := os.OpenFile("public/result_cache.json", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("打开文件失败: %v", err)
	}
	defer f.Close()
	jsonString, _ := json.MarshalIndent(result, "", "  ")
	f.Write(jsonString)
	fmt.Println(len(result["list"].([]any)))
}

// GenerateSec 生成 Sec 请求头（已修复72字节截断）
func GenerateSec(params any) (string, error) {
	// 1. JSON 序列化
	jsonBytes, _ := json.Marshal(params)
	jsonStr := string(jsonBytes)

	// 2. URL 编码
	urlEncoded := url.QueryEscape(jsonStr)

	// 3. Base64
	base64Str := base64.StdEncoding.EncodeToString([]byte(urlEncoded))

	// 4. 拼接密钥
	fullStr := "SECRET_VCG_" + base64Str

	// ---------------- FIX HERE ----------------
	// 关键：截断到 72 字节（和网站JS逻辑一致）
	if len(fullStr) > 72 {
		fullStr = fullStr[:72]
	}
	// -----------------------------------------

	// 5. bcrypt hash
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(fullStr), 8)
	if err != nil {
		return "", err
	}

	return string(hashBytes), nil
}
