package controller

import (
	"strconv"
	"strings"
	"time"

	"webadmin/config"
	"webadmin/model"

	"github.com/gin-gonic/gin"
)

func HeyuesorderList(c *gin.Context) {
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

var items []model.Heyueorder
	var total int64
	var totalUsdt float64
	var totalUsdtLong float64
	var totalUsdtShort float64

	query := model.Db.Model(&model.Heyueorder{})
	statQuery := model.Db.Model(&model.Heyueorder{})
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
		"page":       page,
		"pagesize":   pageSize,
		"total":      total,
		"data":       items,
		"total_usdt": totalUsdt,
		"total_usdt_long":  totalUsdtLong,
		"total_usdt_short": totalUsdtShort,
	})
}

func HeyuesorderCreate(c *gin.Context) {
	var item model.Heyueorder
	if err := c.ShouldBindJSON(&item); err != nil {
		errorResponse(c, 400, "Invalid request")
		return
	}

	item.Username = strings.TrimSpace(item.Username)
	item.Symbol = strings.TrimSpace(item.Symbol)

	if item.Symbol == "" {
		errorResponse(c, 400, "symbol is required")
		return
	}

	item.AddTime = time.Now().Unix()
	item.UpdateTime = time.Now().Unix()

	if err := model.Db.Create(&item).Error; err != nil {
		errorResponse(c, 500, "Failed to create heyuesorder")
		return
	}
	successResponse(c, 201, 1, gin.H{"heyuesorder": item})
}

func HeyuesorderUpdate(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		errorResponse(c, 400, "invalid heyuesorder id")
		return
	}

	var req model.Heyueorder
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, 400, "Invalid request")
		return
	}

	updates := map[string]any{}
	updates["ordertype"] = req.Ordertype
	updates["side"] = req.Side
	updates["price"] = req.Price
	updates["total"] = req.Total
	updates["quantity"] = req.Quantity
	updates["num"] = req.Num
	updates["status"] = req.Status
	updates["usdt"] = req.Usdt
	updates["log"] = req.Log
	updates["update_time"] = time.Now().Unix()

	result := model.Db.Model(&model.Heyueorder{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		errorResponse(c, 500, "Failed to update heyuesorder")
		return
	}
	if result.RowsAffected == 0 {
		errorResponse(c, 404, "heyuesorder not found")
		return
	}

	var item model.Heyueorder
	if err := model.Db.First(&item, id).Error; err != nil {
		errorResponse(c, 500, "Failed to load heyuesorder")
		return
	}

	successResponse(c, 200, 1, gin.H{"heyuesorder": item})
}

func HeyuesorderDelete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		errorResponse(c, 400, "invalid heyuesorder id")
		return
	}

	result := model.Db.Delete(&model.Heyueorder{}, id)
	if result.Error != nil {
		errorResponse(c, 500, "Failed to delete heyuesorder")
		return
	}
	if result.RowsAffected == 0 {
		errorResponse(c, 404, "heyuesorder not found")
		return
	}

	successResponse(c, 200, 1, gin.H{"message": "heyuesorder deleted"})
}
