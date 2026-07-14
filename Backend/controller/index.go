package controller

import (
	"log"
	"webadmin/model"

	"github.com/gin-gonic/gin"
)

type DashboardStats struct {
	AdminCount          int64   `json:"admin_count"`
	ShellCount          int64   `json:"shell_count"`
	HeyueorderUsdt      float64 `json:"heyueorder_usdt"`
	HeyueorderShortUsdt float64 `json:"heyueorder_short_usdt"`
	CoinCount           int64   `json:"coin_count"`
	HeyueCount          int64   `json:"heyue_count"`
	UserCount           int64   `json:"user_count"`
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
	//按 side 分组统计做多/做空收益
	type sideTotal struct {
		Side  int32
		Total float64
	}
	var sideTotals []sideTotal
	if err := model.Db.Model(&model.Heyueorder{}).Select("side, COALESCE(SUM(usdt), 0) as total").Where("ordertype = ?", 2).Group("side").Scan(&sideTotals).Error; err != nil {
		errorResponse(c, 500, "Failed to sum heyueorder usdt by side")
		return
	}
	for _, st := range sideTotals {
		if st.Side == 1 {
			stats.HeyueorderUsdt = st.Total
		} else if st.Side == 2 {
			stats.HeyueorderShortUsdt = st.Total
		}
	}

	log.Printf("Dashboard stats: %+v\n", stats)

	successResponse(c, 200, 1, gin.H{
		"admin_count":           stats.AdminCount,
		"shell_count":           stats.ShellCount,
		"user_count":            stats.UserCount,
		"coin_count":            stats.CoinCount,
		"heyue_count":           stats.HeyueCount,
		"heyueorder_usdt":       stats.HeyueorderUsdt,
		"heyueorder_short_usdt": stats.HeyueorderShortUsdt,
	})
}
