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

// 谷歌收录
func GoogleRevenue() {
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
