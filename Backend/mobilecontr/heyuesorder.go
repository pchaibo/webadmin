package mobilecontr

import (
	"strconv"
	"strings"

	"webadmin/config"
	"webadmin/model"

	"github.com/gin-gonic/gin"
)

func HeyuesorderList(c *gin.Context) {
	userId, _ := c.Get("user_id")
	uid, ok := userId.(int)
	if !ok || uid <= 0 {
		uid = 0
		errorResponse(c, 500, "uid  error")
		return
	}

	page := 1
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	pageSize := 20
	if ps := config.Get("pageSize"); ps != "" {
		if p, err := strconv.Atoi(ps); err == nil && p > 0 {
			pageSize = p
		}
	}
	offset := (page - 1) * pageSize

	symbol := strings.TrimSpace(c.Query("symbol"))
	username := strings.TrimSpace(c.Query("username"))
	statusStr := strings.TrimSpace(c.Query("status"))
	ordertypeStr := strings.TrimSpace(c.Query("ordertype"))
	sideStr := strings.TrimSpace(c.Query("side"))

	var items []model.Heyueorder
	var total int64
	var totalUsdt float64
	var totalUsdtLong float64
	var totalUsdtShort float64

	query := model.Db.Model(&model.Heyueorder{})
	statQuery := model.Db.Model(&model.Heyueorder{})
	if uid > 0 {
		query = query.Where("user_id = ?", uid)
		statQuery = statQuery.Where("user_id = ?", uid)
	}
	if symbol != "" {
		query = query.Where("symbol LIKE ?", "%"+symbol+"%")
		statQuery = statQuery.Where("symbol LIKE ?", "%"+symbol+"%")
	}
	if username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
		statQuery = statQuery.Where("username LIKE ?", "%"+username+"%")
	}
	if statusStr != "" {
		if s, err := strconv.Atoi(statusStr); err == nil {
			query = query.Where("status = ?", s)
			statQuery = statQuery.Where("status = ?", s)
		}
	}
	if ordertypeStr != "" {
		if s, err := strconv.Atoi(ordertypeStr); err == nil {
			query = query.Where("ordertype = ?", s)
			statQuery = statQuery.Where("ordertype = ?", s)
		}
	}
	if sideStr != "" {
		if s, err := strconv.Atoi(sideStr); err == nil {
			query = query.Where("side = ?", s)
			statQuery = statQuery.Where("side = ?", s)
		}
	}

	if err := query.Count(&total).Error; err != nil {
		errorResponse(c, 500, "Failed to count heyuesorders")
		return
	}

	type sideTotal struct {
		Side  int32
		Total float64
	}
	var sideTotals []sideTotal
	if err := statQuery.Select("side, COALESCE(SUM(usdt), 0) as total").Where("ordertype = ?", 2).Group("side").Scan(&sideTotals).Error; err != nil {
		errorResponse(c, 500, "Failed to sum heyuesorder usdt by side")
		return
	}
	for _, st := range sideTotals {
		if st.Side == 1 {
			totalUsdtLong = st.Total
		} else if st.Side == 2 {
			totalUsdtShort = st.Total
		}
	}
	totalUsdt = totalUsdtLong + totalUsdtShort

	if err := query.Order("id desc").Limit(pageSize).Offset(offset).Find(&items).Error; err != nil {
		errorResponse(c, 500, "Failed to retrieve heyuesorders")
		return
	}

	status := 0
	if len(items) > 0 {
		status = 1
	}
	successResponse(c, 200, status, gin.H{
		"page":             page,
		"pagesize":         pageSize,
		"total":            total,
		"data":             items,
		"total_usdt":       totalUsdt,
		"total_usdt_long":  totalUsdtLong,
		"total_usdt_short": totalUsdtShort,
	})
}
