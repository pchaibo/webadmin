package controller

import (
	"strconv"

	"webadmin/model"

	"github.com/gin-gonic/gin"
)

func ShellMinList(c *gin.Context) {
	shellIdStr := c.Query("shell_id")
	if shellIdStr == "" {
		errorResponse(c, 400, "shell_id is required")
		return
	}

	shellId, err := strconv.Atoi(shellIdStr)
	if err != nil || shellId <= 0 {
		errorResponse(c, 400, "invalid shell_id")
		return
	}

	var items []model.ShellMin
	if err := model.Db.Where("shell_id = ?", shellId).Order("id asc").Find(&items).Error; err != nil {
		errorResponse(c, 500, "Failed to query shell min")
		return
	}

	successResponse(c, 200, 1, gin.H{
		"shell_id": shellId,
		"data":     items,
	})
}
