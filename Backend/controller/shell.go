package controller

import (
	"strconv"
	"strings"
	"time"

	"webadmin/config"
	"webadmin/model"

	"github.com/gin-gonic/gin"
)

func ShellList(c *gin.Context) {
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

	var items []model.Shell
	var total int64

	if err := model.Db.Model(&model.Shell{}).Count(&total).Error; err != nil {
		errorResponse(c, 500, "Failed to count shells")
		return
	}

	if err := model.Db.Preload("Group").Order("id asc").Limit(pageSize).Offset(offset).Find(&items).Error; err != nil {
		errorResponse(c, 500, "Failed to retrieve shells")
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

func ShellCreate(c *gin.Context) {
	var item model.Shell
	if err := c.ShouldBindJSON(&item); err != nil {
		errorResponse(c, 400, "Invalid request")
		return
	}

	item.Host = strings.TrimSpace(item.Host)
	item.Scheme = strings.TrimSpace(item.Scheme)
	item.Maxurl = strings.TrimSpace(item.Maxurl)
	item.Minurl = strings.TrimSpace(item.Minurl)
	item.Remark = strings.TrimSpace(item.Remark)

	if item.Host == "" {
		errorResponse(c, 400, "host is required")
		return
	}

	now := int(time.Now().Unix())
	item.Addtime = now
	item.Uptime = now

	if err := model.Db.Create(&item).Error; err != nil {
		errorResponse(c, 500, "Failed to create shell")
		return
	}
	successResponse(c, 201, 1, gin.H{"shell": item})
}

func ShellUpdate(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		errorResponse(c, 400, "invalid shell id")
		return
	}

	var req model.Shell
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, 400, "Invalid request")
		return
	}

	req.Host = strings.TrimSpace(req.Host)
	req.Scheme = strings.TrimSpace(req.Scheme)
	req.Maxurl = strings.TrimSpace(req.Maxurl)
	req.Minurl = strings.TrimSpace(req.Minurl)
	req.Remark = strings.TrimSpace(req.Remark)

	if req.Host == "" {
		errorResponse(c, 400, "host is required")
		return
	}

	updates := map[string]any{
		"host":     req.Host,
		"scheme":   req.Scheme,
		"group_id": req.GroupId,
		"status":   req.Status,
		"num":      req.Num,
		"maxurl":   req.Maxurl,
		"minurl":   req.Minurl,
		"dir":      req.Dir,
		"lock":     req.Lock,
		"remark":   req.Remark,
		"uptime":   int(time.Now().Unix()),
	}

	result := model.Db.Model(&model.Shell{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		errorResponse(c, 500, "Failed to update shell")
		return
	}
	if result.RowsAffected == 0 {
		errorResponse(c, 404, "shell not found")
		return
	}

	var item model.Shell
	if err := model.Db.First(&item, id).Error; err != nil {
		errorResponse(c, 500, "Failed to load shell")
		return
	}

	successResponse(c, 200, 1, gin.H{"shell": item})
}

func ShellDelete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		errorResponse(c, 400, "invalid shell id")
		return
	}

	result := model.Db.Delete(&model.Shell{}, id)
	if result.Error != nil {
		errorResponse(c, 500, "Failed to delete shell")
		return
	}
	if result.RowsAffected == 0 {
		errorResponse(c, 404, "shell not found")
		return
	}

	successResponse(c, 200, 1, gin.H{"message": "shell deleted"})
}
