package controller

import (
	"strconv"
	"strings"

	"webadmin/config"
	"webadmin/model"

	"github.com/gin-gonic/gin"
)

func AuthGroupList(c *gin.Context) {
	page := 1
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	pageSize := 10
	if ps := config.Get("pageSize"); ps != "" {
		if p, err := strconv.Atoi(ps); err == nil && p > 0 {
			pageSize = p
		}
	}
	offset := (page - 1) * pageSize

	var items []model.AuthGroup
	var total int64

	if err := model.Db.Model(&model.AuthGroup{}).Count(&total).Error; err != nil {
		errorResponse(c, 500, "Failed to count auth groups")
		return
	}

	if err := model.Db.Order("id asc").Limit(pageSize).Offset(offset).Find(&items).Error; err != nil {
		errorResponse(c, 500, "Failed to retrieve auth groups")
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

func AuthGroupAll(c *gin.Context) {
	var items []model.AuthGroup
	if err := model.Db.Where("status = ?", 1).Order("id asc").Find(&items).Error; err != nil {
		errorResponse(c, 500, "Failed to retrieve auth groups")
		return
	}

	successResponse(c, 200, 1, gin.H{"data": items})
}

func AuthGroupCreate(c *gin.Context) {
	var item model.AuthGroup
	if err := c.ShouldBindJSON(&item); err != nil {
		errorResponse(c, 400, "Invalid request")
		return
	}

	item.Title = strings.TrimSpace(item.Title)
	item.Rules = strings.TrimSpace(item.Rules)

	if item.Title == "" {
		errorResponse(c, 400, "title is required")
		return
	}

	if item.Status == 0 {
		item.Status = 1
	}

	if err := model.Db.Create(&item).Error; err != nil {
		errorResponse(c, 500, "Failed to create auth group")
		return
	}
	successResponse(c, 201, 1, gin.H{"auth_group": item})
}

func AuthGroupUpdate(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		errorResponse(c, 400, "invalid auth group id")
		return
	}

	var req model.AuthGroup
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, 400, "Invalid request")
		return
	}

	req.Title = strings.TrimSpace(req.Title)
	req.Rules = strings.TrimSpace(req.Rules)

	if req.Title == "" {
		errorResponse(c, 400, "title is required")
		return
	}

	updates := map[string]any{
		"title":  req.Title,
		"status": req.Status,
		"rules":  req.Rules,
	}

	result := model.Db.Model(&model.AuthGroup{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		errorResponse(c, 500, "Failed to update auth group")
		return
	}
	if result.RowsAffected == 0 {
		errorResponse(c, 404, "auth group not found")
		return
	}

	var item model.AuthGroup
	if err := model.Db.First(&item, id).Error; err != nil {
		errorResponse(c, 500, "Failed to load auth group")
		return
	}

	successResponse(c, 200, 1, gin.H{"auth_group": item})
}

func AuthGroupDelete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		errorResponse(c, 400, "invalid auth group id")
		return
	}

	result := model.Db.Delete(&model.AuthGroup{}, id)
	if result.Error != nil {
		errorResponse(c, 500, "Failed to delete auth group")
		return
	}
	if result.RowsAffected == 0 {
		errorResponse(c, 404, "auth group not found")
		return
	}

	successResponse(c, 200, 1, gin.H{"message": "auth group deleted"})
}
