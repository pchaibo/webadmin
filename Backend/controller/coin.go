package controller

import (
	"strconv"
	"strings"
	"time"

	"webadmin/config"
	"webadmin/model"

	"github.com/gin-gonic/gin"
)

func CoinList(c *gin.Context) {
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
	statusStr := strings.TrimSpace(c.Query("status"))

	var items []model.Coin
	var total int64

	query := model.Db.Model(&model.Coin{})
	if symbol != "" {
		query = query.Where("symbol LIKE ?", "%"+symbol+"%")
	}
	if statusStr != "" {
		if s, err := strconv.Atoi(statusStr); err == nil {
			query = query.Where("status = ?", s)
		}
	}

	if err := query.Count(&total).Error; err != nil {
		errorResponse(c, 500, "Failed to count coins")
		return
	}

	if err := query.Order("id asc").Limit(pageSize).Offset(offset).Find(&items).Error; err != nil {
		errorResponse(c, 500, "Failed to retrieve coins")
		return
	}

	status := 0
	if len(items) > 0 {
		status = 1
	}
	successResponse(c, 200, status, gin.H{
		"page":     page,
		"pagesize": pageSize,
		"total":    total,
		"data":     items,
	})
}

func CoinCreate(c *gin.Context) {
	var item model.Coin
	if err := c.ShouldBindJSON(&item); err != nil {
		errorResponse(c, 400, "Invalid request")
		return
	}

	item.Name = strings.TrimSpace(item.Name)
	item.Symbol = strings.TrimSpace(item.Symbol)

	if item.Name == "" {
		errorResponse(c, 400, "name is required")
		return
	}
	if item.Symbol == "" {
		errorResponse(c, 400, "symbol is required")
		return
	}

	item.AddTime = time.Now().Unix()
	item.UpdateTime = time.Now().Unix()

	if err := model.Db.Create(&item).Error; err != nil {
		errorResponse(c, 500, "Failed to create coin")
		return
	}
	successResponse(c, 201, 1, gin.H{"coin": item})
}

func CoinUpdate(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		errorResponse(c, 400, "invalid coin id")
		return
	}

	var req model.Coin
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, 400, "Invalid request")
		return
	}

	updates := map[string]any{}

	if req.Name != "" {
		req.Name = strings.TrimSpace(req.Name)
		updates["name"] = req.Name
	}
	if req.Symbol != "" {
		req.Symbol = strings.TrimSpace(req.Symbol)
		updates["symbol"] = req.Symbol
	}
	updates["close"] = req.Close
	updates["priceprecision"] = req.Priceprecision
	updates["open"] = req.Open
	updates["low"] = req.Low
	updates["high"] = req.High
	updates["status"] = req.Status
	updates["update_time"] = time.Now().Unix()

	result := model.Db.Model(&model.Coin{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		errorResponse(c, 500, "Failed to update coin")
		return
	}
	if result.RowsAffected == 0 {
		errorResponse(c, 404, "coin not found")
		return
	}

	var item model.Coin
	if err := model.Db.First(&item, id).Error; err != nil {
		errorResponse(c, 500, "Failed to load coin")
		return
	}

	successResponse(c, 200, 1, gin.H{"coin": item})
}

func CoinDelete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		errorResponse(c, 400, "invalid coin id")
		return
	}

	result := model.Db.Delete(&model.Coin{}, id)
	if result.Error != nil {
		errorResponse(c, 500, "Failed to delete coin")
		return
	}
	if result.RowsAffected == 0 {
		errorResponse(c, 404, "coin not found")
		return
	}

	successResponse(c, 200, 1, gin.H{"message": "coin deleted"})
}
