package controller

import (
	"strconv"
	"strings"
	"time"

	"webadmin/config"
	"webadmin/model"

	"github.com/gin-gonic/gin"
)

func TaskList(c *gin.Context) {
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
	conditionStr := strings.TrimSpace(c.Query("condition"))

	var items []model.HeyueTask
	var total int64

	query := model.Db.Model(&model.HeyueTask{})
	if symbol != "" {
		query = query.Where("symbol LIKE ?", "%"+symbol+"%")
	}
	if username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}
	if statusStr != "" {
		if s, err := strconv.Atoi(statusStr); err == nil {
			query = query.Where("status = ?", s)
		}
	}
	if conditionStr != "" {
		if s, err := strconv.Atoi(conditionStr); err == nil {
			query = query.Where("`condition` = ?", s)
		}
	}

	if err := query.Count(&total).Error; err != nil {
		errorResponse(c, 500, "Failed to count tasks")
		return
	}

	if err := query.Order("id asc").Limit(pageSize).Offset(offset).Find(&items).Error; err != nil {
		errorResponse(c, 500, "Failed to retrieve tasks")
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

func TaskCreate(c *gin.Context) {
	var item model.HeyueTask
	if err := c.ShouldBindJSON(&item); err != nil {
		errorResponse(c, 400, "Invalid request")
		return
	}

	item.Symbol = strings.TrimSpace(item.Symbol)
	item.Username = strings.TrimSpace(item.Username)

	if item.Symbol == "" {
		errorResponse(c, 400, "symbol is required")
		return
	}

	if item.Coinid == 0 {
		var coin model.Coin
		if err := model.Db.Where("symbol = ?", item.Symbol).First(&coin).Error; err == nil {
			item.Coinid = int(coin.Id)
		}
	}

	if item.Userid == 0 && item.Username != "" {
		var user model.User
		if err := model.Db.Where("username = ?", item.Username).First(&user).Error; err == nil {
			item.Userid = user.Id
		}
	}

	item.AddTime = time.Now().Unix()
	item.UpdateTime = time.Now().Unix()

	if err := model.Db.Create(&item).Error; err != nil {
		errorResponse(c, 500, "Failed to create task")
		return
	}
	successResponse(c, 201, 1, gin.H{"task": item})
}

func TaskUpdate(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		errorResponse(c, 400, "invalid task id")
		return
	}

	var req model.HeyueTask
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, 400, "Invalid request")
		return
	}

	updates := map[string]any{}
	updates["coinid"] = req.Coinid
	updates["symbol"] = strings.TrimSpace(req.Symbol)
	updates["userid"] = req.Userid
	updates["username"] = strings.TrimSpace(req.Username)
	updates["price"] = req.Price
	updates["condition"] = req.Condition
	updates["status"] = req.Status
	updates["update_time"] = time.Now().Unix()

	result := model.Db.Model(&model.HeyueTask{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		errorResponse(c, 500, "Failed to update task")
		return
	}
	if result.RowsAffected == 0 {
		errorResponse(c, 404, "task not found")
		return
	}

	var item model.HeyueTask
	if err := model.Db.First(&item, id).Error; err != nil {
		errorResponse(c, 500, "Failed to load task")
		return
	}

	successResponse(c, 200, 1, gin.H{"task": item})
}

func TaskDelete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		errorResponse(c, 400, "invalid task id")
		return
	}

	result := model.Db.Delete(&model.HeyueTask{}, id)
	if result.Error != nil {
		errorResponse(c, 500, "Failed to delete task")
		return
	}
	if result.RowsAffected == 0 {
		errorResponse(c, 404, "task not found")
		return
	}

	successResponse(c, 200, 1, gin.H{"message": "task deleted"})
}
