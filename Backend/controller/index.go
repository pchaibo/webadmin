package controller

import (
	"log"
	"webadmin/model"

	"github.com/gin-gonic/gin"
)

type DashboardStats struct {
	AdminCount     int64   `json:"admin_count"`
	ShellCount     int64   `json:"shell_count"`
	HeyueorderUsdt float64 `json:"heyueorder_usdt"`
	CoinCount      int64   `json:"coin_count"`
	HeyueCount     int64   `json:"heyue_count"`
	UserCount      int64   `json:"user_count"`
}

func Index(c *gin.Context) {
	var stats DashboardStats

	if err := model.Db.Model(&model.Admin{}).Count(&stats.AdminCount).Error; err != nil {
		errorResponse(c, 500, "Failed to count admins")
		return
	}

	if err := model.Db.Model(&model.Shell{}).Count(&stats.ShellCount).Error; err != nil {
		errorResponse(c, 500, "Failed to count shells")
		return
	}

	if err := model.Db.Model(&model.User{}).Count(&stats.UserCount).Error; err != nil {
		errorResponse(c, 500, "Failed to count users")
		return
	}

	if err := model.Db.Model(&model.Coin{}).Where("status = ?", 1).Count(&stats.CoinCount).Error; err != nil {
		errorResponse(c, 500, "Failed to count coins")
		return
	}

	if err := model.Db.Model(&model.Heyue{}).Count(&stats.HeyueCount).Error; err != nil {
		errorResponse(c, 500, "Failed to count heyue")
		return
	}

	if err := model.Db.Model(&model.Heyueorder{}).Where("ordertype = ?", 2).Select("COALESCE(SUM(usdt), 0)").Scan(&stats.HeyueorderUsdt).Error; err != nil {
		errorResponse(c, 500, "Failed to sum heyueorder usdt")
		return
	}
	log.Printf("Dashboard stats: %+v\n", stats)

	successResponse(c, 200, 1, gin.H{
		"admin_count":     stats.AdminCount,
		"shell_count":     stats.ShellCount,
		"user_count":      stats.UserCount,
		"coin_count":      stats.CoinCount,
		"heyue_count":     stats.HeyueCount,
		"heyueorder_usdt": stats.HeyueorderUsdt,
	})
}
