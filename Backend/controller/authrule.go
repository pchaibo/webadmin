package controller

import (
	"strconv"
	"strings"

	"webadmin/config"
	"webadmin/model"

	"github.com/gin-gonic/gin"
)

type authRuleRequest struct {
	Pid       int    `json:"pid"`
	Name      string `json:"name"`
	Title     string `json:"title"`
	Icon      string `json:"icon"`
	Type      int    `json:"type"`
	Status    int    `json:"status"`
	Condition string `json:"condition"`
}

func AuthRuleList(c *gin.Context) {
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

	var rules []model.AuthRule
	var total int64

	if err := model.Db.Model(&model.AuthRule{}).Count(&total).Error; err != nil {
		errorResponse(c, 500, "Failed to count auth rules")
		return
	}

	if err := model.Db.Order("id asc").Limit(pageSize).Offset(offset).Find(&rules).Error; err != nil {
		errorResponse(c, 500, "Failed to retrieve auth rules")
		return
	}

	status := 0
	if len(rules) > 0 {
		status = 1
	}
	successResponse(c, 200, status, gin.H{
		"page":     page,
		"pagesize": pageSize,
		"total":    total,
		"data":     rules,
	})
}

func AuthRuleCreate(c *gin.Context) {
	var req authRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, 400, "Invalid request")
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Title = strings.TrimSpace(req.Title)
	req.Icon = strings.TrimSpace(req.Icon)
	req.Condition = strings.TrimSpace(req.Condition)

	if req.Name == "" {
		errorResponse(c, 400, "name is required")
		return
	}
	if req.Title == "" {
		errorResponse(c, 400, "title is required")
		return
	}

	rule := model.AuthRule{
		Pid:       req.Pid,
		Name:      req.Name,
		Title:     req.Title,
		Icon:      req.Icon,
		Type:      req.Type,
		Status:    req.Status,
		Condition: req.Condition,
	}

	if err := model.Db.Create(&rule).Error; err != nil {
		errorResponse(c, 500, "Failed to create auth rule")
		return
	}

	successResponse(c, 201, 1, gin.H{"auth_rule": rule})
}

func AuthRuleUpdate(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		errorResponse(c, 400, "invalid auth rule id")
		return
	}

	var req authRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, 400, "Invalid request")
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Title = strings.TrimSpace(req.Title)
	req.Icon = strings.TrimSpace(req.Icon)
	req.Condition = strings.TrimSpace(req.Condition)

	if req.Name == "" {
		errorResponse(c, 400, "name is required")
		return
	}
	if req.Title == "" {
		errorResponse(c, 400, "title is required")
		return
	}

	updates := map[string]any{
		"pid":       req.Pid,
		"name":      req.Name,
		"title":     req.Title,
		"icon":      req.Icon,
		"type":      req.Type,
		"status":    req.Status,
		"condition": req.Condition,
	}

	result := model.Db.Model(&model.AuthRule{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		errorResponse(c, 500, "Failed to update auth rule")
		return
	}
	if result.RowsAffected == 0 {
		errorResponse(c, 404, "auth rule not found")
		return
	}

	var rule model.AuthRule
	if err := model.Db.First(&rule, id).Error; err != nil {
		errorResponse(c, 500, "Failed to load auth rule")
		return
	}

	successResponse(c, 200, 1, gin.H{"auth_rule": rule})
}

func AuthRuleDelete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		errorResponse(c, 400, "invalid auth rule id")
		return
	}

	result := model.Db.Delete(&model.AuthRule{}, id)
	if result.Error != nil {
		errorResponse(c, 500, "Failed to delete auth rule")
		return
	}
	if result.RowsAffected == 0 {
		errorResponse(c, 404, "auth rule not found")
		return
	}

	successResponse(c, 200, 1, gin.H{"message": "auth rule deleted"})
}
