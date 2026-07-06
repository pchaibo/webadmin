package controller

import (
	"strconv"
	"strings"
	"time"

	"webadmin/config"
	"webadmin/model"

	"github.com/gin-gonic/gin"
)

func HeyueList(c *gin.Context) {
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

	var items []model.Heyue
	var total int64

	query := model.Db.Model(&model.Heyue{})
	if symbol != "" {
		query = query.Where("symbol LIKE ?", "%"+symbol+"%")
	}
	if username != "" {
		query = query.Where("user_name LIKE ?", "%"+username+"%")
	}
	if statusStr != "" {
		if s, err := strconv.Atoi(statusStr); err == nil {
			query = query.Where("status = ?", s)
		}
	}

	if err := query.Count(&total).Error; err != nil {
		errorResponse(c, 500, "Failed to count heyues")
		return
	}

	if err := query.Order("id asc").Limit(pageSize).Offset(offset).Find(&items).Error; err != nil {
		errorResponse(c, 500, "Failed to retrieve heyues")
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

func HeyueCreate(c *gin.Context) {
	var item model.Heyue
	if err := c.ShouldBindJSON(&item); err != nil {
		errorResponse(c, 400, "Invalid request")
		return
	}

	item.UserName = strings.TrimSpace(item.UserName)
	item.Symbol = strings.TrimSpace(item.Symbol)
	if item.UserId == 0 {
		var user model.User
		if err := model.Db.Where("username = ?", item.UserName).First(&user).Error; err != nil {
			errorResponse(c, 400, "user not found")
			return
		}
		item.UserId = uint(user.Id)
	}

	if item.Symbol == "" {
		errorResponse(c, 400, "symbol is required")
		return
	}

	item.AddTime = time.Now().Unix()
	item.UpdateTime = time.Now().Unix()

	if err := model.Db.Create(&item).Error; err != nil {
		errorResponse(c, 500, "Failed to create heyue")
		return
	}
	successResponse(c, 201, 1, gin.H{"heyue": item})
}

func HeyueUpdate(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		errorResponse(c, 400, "invalid heyue id")
		return
	}

	var req model.Heyue
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, 400, "Invalid request")
		return
	}

	updates := map[string]any{}

	// if req.UserId != 0 {
	// 	updates["user_id"] = req.UserId
	// }
	// if req.UserName != "" {
	// 	req.UserName = strings.TrimSpace(req.UserName)
	// 	updates["username"] = req.UserName
	// }
	// if req.Symbol != "" {
	// 	req.Symbol = strings.TrimSpace(req.Symbol)
	// 	updates["symbol"] = req.Symbol
	// }
	updates["side"] = req.Side
	updates["num"] = req.Num
	updates["status"] = req.Status
	updates["sellprice"] = req.Sellprice
	updates["oneprice"] = req.Oneprice
	updates["repeatprice"] = req.Repeatprice
	updates["rangetype"] = req.Rangetype
	updates["rangeprice"] = req.Rangeprice
	updates["rangepercent"] = req.Rangepercent
	updates["rangeclosingpct"] = req.Rangeclosingpct
	updates["rangeclosing"] = req.Rangeclosing
	updates["closingprice"] = req.Closingprice
	updates["risk"] = req.Risk
	updates["risk_time"] = req.RiskTime
	updates["update_time"] = time.Now().Unix()

	result := model.Db.Model(&model.Heyue{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		errorResponse(c, 500, "Failed to update heyue")
		return
	}
	if result.RowsAffected == 0 {
		errorResponse(c, 404, "heyue not found")
		return
	}

	var item model.Heyue
	if err := model.Db.First(&item, id).Error; err != nil {
		errorResponse(c, 500, "Failed to load heyue")
		return
	}

	successResponse(c, 200, 1, gin.H{"heyue": item})
}

func HeyueDelete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		errorResponse(c, 400, "invalid heyue id")
		return
	}

	result := model.Db.Delete(&model.Heyue{}, id)
	if result.Error != nil {
		errorResponse(c, 500, "Failed to delete heyue")
		return
	}
	if result.RowsAffected == 0 {
		errorResponse(c, 404, "heyue not found")
		return
	}

	successResponse(c, 200, 1, gin.H{"message": "heyue deleted"})
}
