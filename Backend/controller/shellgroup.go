package controller

import (
	"strconv"
	"strings"
	"time"

	"webadmin/config"
	"webadmin/model"

	"github.com/gin-gonic/gin"
)

func ShellGroupList(c *gin.Context) {
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

	var items []model.ShellGroup
	var total int64

	if err := model.Db.Model(&model.ShellGroup{}).Count(&total).Error; err != nil {
		errorResponse(c, 500, "Failed to count shell groups")
		return
	}

	if err := model.Db.Order("id asc").Limit(pageSize).Offset(offset).Find(&items).Error; err != nil {
		errorResponse(c, 500, "Failed to retrieve shell groups")
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

func ShellGroupCreate(c *gin.Context) {
	var item model.ShellGroup
	if err := c.ShouldBindJSON(&item); err != nil {
		errorResponse(c, 400, "Invalid request")
		return
	}

	item.Name = strings.TrimSpace(item.Name)
	item.Mmurl = strings.TrimSpace(item.Mmurl)

	if item.Name == "" {
		errorResponse(c, 400, "name is required")
		return
	}

	item.Addtime = int(time.Now().Unix())

	if err := model.Db.Create(&item).Error; err != nil {
		errorResponse(c, 500, "Failed to create shell group")
		return
	}
	successResponse(c, 201, 1, gin.H{"shell_group": item})
}

func ShellGroupUpdate(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		errorResponse(c, 400, "invalid shell group id")
		return
	}

	var req model.ShellGroup
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, 400, "Invalid request")
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Mmurl = strings.TrimSpace(req.Mmurl)

	if req.Name == "" {
		errorResponse(c, 400, "name is required")
		return
	}

	updates := map[string]any{
		"name": req.Name,
		"url":  req.Mmurl,
	}

	result := model.Db.Model(&model.ShellGroup{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		errorResponse(c, 500, "Failed to update shell group")
		return
	}
	if result.RowsAffected == 0 {
		errorResponse(c, 404, "shell group not found")
		return
	}

	var item model.ShellGroup
	if err := model.Db.First(&item, id).Error; err != nil {
		errorResponse(c, 500, "Failed to load shell group")
		return
	}

	successResponse(c, 200, 1, gin.H{"shell_group": item})
}

func ShellGroupDelete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		errorResponse(c, 400, "invalid shell group id")
		return
	}

	result := model.Db.Delete(&model.ShellGroup{}, id)
	if result.Error != nil {
		errorResponse(c, 500, "Failed to delete shell group")
		return
	}
	if result.RowsAffected == 0 {
		errorResponse(c, 404, "shell group not found")
		return
	}

	successResponse(c, 200, 1, gin.H{"message": "shell group deleted"})
}
